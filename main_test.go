package main

import (
	"testing"

	"cliamp/playlist"
)

func TestOfflineRemoteInputDetection(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"local absolute", "/Volumes/music/song.flac", false},
		{"local relative", "Music/song.flac", false},
		{"windows unc", `\\server\share\song.flac`, false},
		{"http", "https://example.com/song.mp3", true},
		{"yt search", "ytsearch1:lofi", true},
		{"spotify", "spotify:track:abc", true},
		{"ssh", "ssh://host/path/song.flac", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRemoteTrackPath(tt.path); got != tt.want {
				t.Fatalf("isRemoteTrackPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFirstRemoteTrackPath(t *testing.T) {
	tracks := []playlist.Track{
		{Path: "/music/local.flac"},
		{Path: "https://example.com/stream", Stream: true},
	}
	if got := firstRemoteTrackPath(tracks); got != "https://example.com/stream" {
		t.Fatalf("firstRemoteTrackPath = %q", got)
	}
}
