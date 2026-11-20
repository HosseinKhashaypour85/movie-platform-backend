package transcoder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Rendition struct {
	Name         string
	Height       int
	VideoBitrate string
	MaxRate      string
	BufferSize   string
}

var renditionLadder = []Rendition{
	{Name: "360p", Height: 360, VideoBitrate: "800k", MaxRate: "856k", BufferSize: "1200k"},
	{Name: "480p", Height: 480, VideoBitrate: "1400k", MaxRate: "1498k", BufferSize: "2100k"},
	{Name: "720p", Height: 720, VideoBitrate: "2800k", MaxRate: "2996k", BufferSize: "4200k"},
	{Name: "1080p", Height: 1080, VideoBitrate: "5000k", MaxRate: "5350k", BufferSize: "7500k"},
}

// SelectRenditions never upscales a source: only output heights at or below
// the source's display height are included.
func SelectRenditions(sourceHeight int) []Rendition {
	selected := make([]Rendition, 0, len(renditionLadder))
	for _, r := range renditionLadder {
		if r.Height <= sourceHeight {
			selected = append(selected, r)
		}
	}
	return selected
}

type probeResult struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	} `json:"streams"`
}

func (w *Worker) probe(ctx context.Context, input string) (int, int, bool, error) {
	cmd := exec.CommandContext(ctx, w.FFprobePath, "-v", "error", "-show_streams", "-of", "json", input)
	raw, err := cmd.Output()
	if err != nil {
		return 0, 0, false, fmt.Errorf("ffprobe failed: %w", err)
	}
	var result probeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return 0, 0, false, fmt.Errorf("parse ffprobe output: %w", err)
	}
	width, height, hasAudio := 0, 0, false
	for _, stream := range result.Streams {
		if stream.CodecType == "video" && stream.Width > 0 && stream.Height > 0 {
			width, height = stream.Width, stream.Height
			break
		}
	}
	for _, stream := range result.Streams {
		if stream.CodecType == "audio" {
			hasAudio = true
			break
		}
	}
	if width == 0 || height == 0 {
		return 0, 0, false, errors.New("source has no valid video stream")
	}
	return width, height, hasAudio, nil
}

// BuildABRArguments creates a multi-variant HLS VOD output. The caller prepends
// the input arguments and ensures the rendition directories exist.
func BuildABRArguments(outDir string, sourceWidth, sourceHeight int, renditions []Rendition, hasAudio bool) ([]string, error) {
	if sourceWidth <= 0 || sourceHeight <= 0 || len(renditions) == 0 {
		return nil, errors.New("invalid source dimensions or empty rendition ladder")
	}
	args := make([]string, 0, 64)
	filters := make([]string, 0, len(renditions)+1)
	splitLabels := make([]string, len(renditions))
	for i := range renditions {
		splitLabels[i] = fmt.Sprintf("[v%d]", i)
	}
	filters = append(filters, fmt.Sprintf("[0:v:0]split=%d%s", len(renditions), strings.Join(splitLabels, "")))
	for i, r := range renditions {
		if r.Height <= 0 || r.Height > sourceHeight {
			return nil, fmt.Errorf("rendition %s would upscale the source", r.Name)
		}
		filters = append(filters, fmt.Sprintf("[v%d]scale=-2:%d:flags=lanczos[vout%d]", i, r.Height, i))
	}
	args = append(args, "-filter_complex", strings.Join(filters, ";"))
	var streamMap strings.Builder
	for i, r := range renditions {
		args = append(args, "-map", fmt.Sprintf("[vout%d]", i))
		if hasAudio {
			args = append(args, "-map", "0:a:0")
		}
		if i > 0 {
			streamMap.WriteByte(' ')
		}
		fmt.Fprintf(&streamMap, "v:%d", i)
		if hasAudio {
			fmt.Fprintf(&streamMap, ",a:%d", i)
		}
		fmt.Fprintf(&streamMap, ",name:%s", r.Name)
	}
	args = append(args, "-c:v", "libx264", "-preset", "veryfast", "-pix_fmt", "yuv420p")
	for i, r := range renditions {
		args = append(args, fmt.Sprintf("-b:v:%d", i), r.VideoBitrate, fmt.Sprintf("-maxrate:v:%d", i), r.MaxRate, fmt.Sprintf("-bufsize:v:%d", i), r.BufferSize)
	}
	if hasAudio {
		args = append(args, "-c:a", "aac", "-b:a", "128k", "-ac", "2")
	}
	args = append(args, "-f", "hls", "-hls_time", "6", "-hls_playlist_type", "vod", "-hls_flags", "independent_segments",
		"-master_pl_name", "master.m3u8", "-var_stream_map", streamMap.String(),
		"-hls_segment_filename", filepath.ToSlash(filepath.Join(outDir, "%v", "segment_%05d.ts")),
		filepath.ToSlash(filepath.Join(outDir, "%v", "index.m3u8")))
	return args, nil
}

func createRenditionDirs(root string, renditions []Rendition) error {
	for _, r := range renditions {
		if err := os.MkdirAll(filepath.Join(root, r.Name), 0700); err != nil {
			return err
		}
	}
	return nil
}
