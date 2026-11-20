package transcoder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"movie-platform/modules/media"
)

type Worker struct {
	DB *gorm.DB
	Storage *media.S3Presigner
	FFmpegPath string
	FFprobePath string
	TempDir string
	HTTP *http.Client
	MaxSourceBytes int64
}

func New(db *gorm.DB, storage *media.S3Presigner) (*Worker, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil { return nil, errors.New("ffmpeg executable was not found in PATH") }
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil { return nil, errors.New("ffprobe executable was not found in PATH") }
	temp := os.Getenv("TRANSCODE_TEMP_DIR")
	if temp == "" { temp = os.TempDir() }
	return &Worker{DB:db, Storage:storage, FFmpegPath:ffmpeg, FFprobePath:ffprobe, TempDir:temp, HTTP:&http.Client{Timeout:30*time.Minute}, MaxSourceBytes:2<<30}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	job, err := w.claim(ctx)
	if err != nil || job == nil { return err }
	if err := w.process(ctx, job); err != nil {
		msg := err.Error()
		if len(msg)>1000 { msg=msg[:1000] }
		status := "failed"
		if job.Attempts < 3 { status = "queued" }
		_ = w.DB.WithContext(ctx).Model(job).Updates(map[string]any{"status":status,"error_message":msg,"locked_at":nil}).Error
		return fmt.Errorf("transcode job %d: %w", job.ID, err)
	}
	return nil
}

func (w *Worker) claim(ctx context.Context) (*media.TranscodeJob, error) {
	var job media.TranscodeJob
	now:=time.Now().UTC()
	err:=w.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err:=tx.Clauses(clause.Locking{Strength:"UPDATE",Options:"SKIP LOCKED"}).
			Where("(status = ? AND attempts < ?) OR (status = ? AND locked_at < ?)","queued",3,"processing",now.Add(-2*time.Hour)).Order("id ASC").First(&job).Error; err!=nil { return err }
		return tx.Model(&job).Updates(map[string]any{"status":"processing","locked_at":now,"attempts":gorm.Expr("attempts + 1")}).Error
	})
	if errors.Is(err,gorm.ErrRecordNotFound) { return nil,nil }
	if err!=nil { return nil,err }
	job.Attempts++
	return &job,nil
}

func (w *Worker) process(ctx context.Context, job *media.TranscodeJob) error {
	var source media.Asset
	if err:=w.DB.WithContext(ctx).First(&source,job.SourceAssetID).Error; err!=nil { return err }
	if source.Kind!="video" || source.MimeType!="video/mp4" || source.IsActive { return errors.New("source asset is not an inactive MP4 video") }
	downloadURL,err:=w.Storage.PresignGet(source.StorageKey,time.Now(),30*time.Minute)
	if err!=nil { return err }
	work,err:=os.MkdirTemp(w.TempDir,"movie-hls-*")
	if err!=nil { return err }
	defer os.RemoveAll(work)
	input:=filepath.Join(work,"source.mp4")
	if err:=w.download(ctx,downloadURL,input); err!=nil { return err }
	outDir:=filepath.Join(work,"hls")
	if err:=os.MkdirAll(outDir,0700); err!=nil { return err }
	masterPlaylist:=filepath.Join(outDir,"master.m3u8")
	ffCtx,cancel:=context.WithTimeout(ctx,6*time.Hour)
	defer cancel()
	width,height,err:=w.probe(ffCtx,input)
	if err!=nil { return err }
	renditions:=SelectRenditions(height)
	if len(renditions)==0 { return errors.New("source video is too small for supported HLS renditions") }
	args,err:=BuildABRArguments(outDir,width,height,renditions)
	if err!=nil { return err }
	args=append([]string{"-hide_banner","-nostdin","-y","-i",input},args...)
	cmd:=exec.CommandContext(ffCtx,w.FFmpegPath,args...)
	logPath:=filepath.Join(work,"ffmpeg.log")
	logFile,_:=os.Create(logPath)
	if logFile!=nil { cmd.Stdout=logFile; cmd.Stderr=logFile; defer logFile.Close() }
	if err:=cmd.Run(); err!=nil {
		if ffCtx.Err()!=nil { return ffCtx.Err() }
		if logFile!=nil { _=logFile.Sync(); raw,_:=os.ReadFile(logPath); if len(raw)>2500 { raw=raw[len(raw)-2500:] }; return fmt.Errorf("ffmpeg failed: %s",strings.TrimSpace(string(raw))) }
		return fmt.Errorf("ffmpeg failed: %w",err)
	}
	files:=make([]string,0,32)
	err=filepath.WalkDir(outDir,func(file string,entry fs.DirEntry,walkErr error) error {
		if walkErr!=nil { return walkErr }
		if !entry.IsDir() { files=append(files,file) }
		return nil
	})
	if err!=nil { return err }
	sort.Strings(files)
	if len(files)<3 { return errors.New("ffmpeg did not produce an ABR master playlist, variant playlists and segments") }
	masterBytes,err:=os.ReadFile(masterPlaylist)
	if err!=nil || !strings.Contains(string(masterBytes),"#EXT-X-STREAM-INF") { return errors.New("generated HLS master playlist is invalid") }
	for _,rendition:=range renditions {
		variant:=filepath.Join(outDir,rendition.Name,"index.m3u8")
		data,readErr:=os.ReadFile(variant)
		if readErr!=nil || !strings.Contains(string(data),"#EXTINF:") { return fmt.Errorf("generated %s variant playlist is invalid",rendition.Name) }
	}
	var random [12]byte
	if _,err:=rand.Read(random[:]); err!=nil { return err }
	prefix:=fmt.Sprintf("hls/%d/%s",source.ID,hex.EncodeToString(random[:]))
	var totalOutputBytes uint64
	for _,file:=range files {
		info,statErr:=os.Stat(file)
		if statErr!=nil { return statErr }
		totalOutputBytes += uint64(info.Size())
		rel,relErr:=filepath.Rel(outDir,file)
		if relErr!=nil { return relErr }
		name:=filepath.ToSlash(rel)
		contentType:="video/mp2t"
		if strings.HasSuffix(name,".m3u8") { contentType="application/vnd.apple.mpegurl" }
		if err:=w.upload(ctx,prefix+"/"+name,file,contentType); err!=nil { return err }
	}
	outputKey:=prefix+"/master.m3u8"
	output:=media.Asset{MovieID:source.MovieID,EpisodeID:source.EpisodeID,Kind:"video",StorageKey:outputKey,MimeType:"application/vnd.apple.mpegurl",SizeBytes:totalOutputBytes,IsActive:true}
	return w.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err:=tx.Create(&output).Error; err!=nil { return err }
		return tx.Model(job).Updates(map[string]any{"status":"completed","error_message":"","locked_at":nil}).Error
	})
}

func (w *Worker) download(ctx context.Context, signedURL, destination string) error {
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,signedURL,nil); if err!=nil { return err }
	resp,err:=w.HTTP.Do(req); if err!=nil { return err }; defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK { return fmt.Errorf("source download returned status %d",resp.StatusCode) }
	if resp.ContentLength<=0 || resp.ContentLength>w.MaxSourceBytes { return errors.New("source object has invalid or excessive size") }
	out,err:=os.OpenFile(destination,os.O_CREATE|os.O_EXCL|os.O_WRONLY,0600); if err!=nil { return err }; defer out.Close()
	n,err:=io.Copy(out,io.LimitReader(resp.Body,w.MaxSourceBytes+1)); if err!=nil { return err }
	if n==0 || n>w.MaxSourceBytes { return errors.New("downloaded source is empty or too large") }
	return out.Sync()
}

func (w *Worker) upload(ctx context.Context,key,file,contentType string) error {
	signedURL,err:=w.Storage.PresignPut(key,contentType,time.Now(),10*time.Minute); if err!=nil { return err }
	f,err:=os.Open(file); if err!=nil { return err }; defer f.Close()
	info,err:=f.Stat(); if err!=nil { return err }
	req,err:=http.NewRequestWithContext(ctx,http.MethodPut,signedURL,f); if err!=nil { return err }
	req.ContentLength=info.Size()
	req.Header.Set("Content-Type",contentType)
	resp,err:=w.HTTP.Do(req); if err!=nil { return err }; defer resp.Body.Close()
	_,_=io.Copy(io.Discard,io.LimitReader(resp.Body,4096))
	if resp.StatusCode<200 || resp.StatusCode>=300 { return fmt.Errorf("HLS object upload returned status %d",resp.StatusCode) }
	return nil
}
