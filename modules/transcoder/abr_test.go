package transcoder

import (
	"strings"
	"testing"
)

func TestSelectRenditionsDoesNotUpscale(t *testing.T) {
	cases:=[]struct{height int; want []string}{
		{240,nil},
		{360,[]string{"360p"}},
		{720,[]string{"360p","480p","720p"}},
		{1080,[]string{"360p","480p","720p","1080p"}},
	}
	for _,tc:=range cases {
		got:=SelectRenditions(tc.height)
		if len(got)!=len(tc.want) { t.Fatalf("height %d: got %d renditions, want %d",tc.height,len(got),len(tc.want)) }
		for i,name:=range tc.want { if got[i].Name!=name { t.Fatalf("height %d: rendition %d = %s, want %s",tc.height,i,got[i].Name,name) } }
	}
}

func TestBuildABRArguments(t *testing.T) {
	renditions:=SelectRenditions(720)
	args,err:=BuildABRArguments("/tmp/hls",1280,720,renditions,true)
	if err!=nil { t.Fatal(err) }
	joined:=strings.Join(args," ")
	for _,want:=range []string{"master.m3u8","360p","480p","720p","libx264","aac","%v/segment_%05d.ts"} {
		if !strings.Contains(joined,want) { t.Errorf("arguments do not contain %q: %s",want,joined) }
	}
	if strings.Contains(joined,"1080p") { t.Fatal("720p source must not produce a 1080p rendition") }
}

func TestBuildABRArgumentsWithoutAudio(t *testing.T) {
	args,err:=BuildABRArguments("/tmp/hls",640,360,SelectRenditions(360),false)
	if err!=nil { t.Fatal(err) }
	joined:=strings.Join(args," ")
	if strings.Contains(joined,"-c:a") || strings.Contains(joined,",a:") { t.Fatalf("audio options included for silent source: %s",joined) }
}

func TestBuildABRArgumentsRejectsUpscale(t *testing.T) {
	_,err:=BuildABRArguments("/tmp/hls",640,360,[]Rendition{{Name:"720p",Height:720}},false)
	if err==nil { t.Fatal("expected upscale to be rejected") }
}
