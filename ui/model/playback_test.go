package model

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"cliamp/playlist"
	"cliamp/ui"
)

type playbackFakeEngine struct {
	playing           bool
	playPaths         []string
	preloadPaths      []string
	clearPreloadCalls int
}

func (f *playbackFakeEngine) Play(path string, _ time.Duration) error {
	f.playing = true
	f.playPaths = append(f.playPaths, path)
	return nil
}
func (f *playbackFakeEngine) PlayYTDL(path string, _ time.Duration) error {
	f.playing = true
	f.playPaths = append(f.playPaths, path)
	return nil
}
func (f *playbackFakeEngine) Preload(path string, _ time.Duration) error {
	f.preloadPaths = append(f.preloadPaths, path)
	return nil
}
func (f *playbackFakeEngine) PreloadYTDL(path string, _ time.Duration) error {
	f.preloadPaths = append(f.preloadPaths, path)
	return nil
}
func (f *playbackFakeEngine) ClearPreload()                { f.clearPreloadCalls++ }
func (f *playbackFakeEngine) Stop()                        { f.playing = false }
func (f *playbackFakeEngine) Close()                       {}
func (f *playbackFakeEngine) TogglePause()                 {}
func (f *playbackFakeEngine) Seek(time.Duration) error     { return nil }
func (f *playbackFakeEngine) SeekYTDL(time.Duration) error { return nil }
func (f *playbackFakeEngine) CancelSeekYTDL()              {}
func (f *playbackFakeEngine) IsPlaying() bool              { return f.playing }
func (f *playbackFakeEngine) IsPaused() bool               { return false }
func (f *playbackFakeEngine) Drained() bool                { return false }
func (f *playbackFakeEngine) HasPreload() bool             { return false }
func (f *playbackFakeEngine) Seekable() bool               { return false }
func (f *playbackFakeEngine) IsStreamSeek() bool           { return false }
func (f *playbackFakeEngine) IsYTDLSeek() bool             { return false }
func (f *playbackFakeEngine) GaplessAdvanced() bool        { return false }
func (f *playbackFakeEngine) Position() time.Duration      { return 0 }
func (f *playbackFakeEngine) Duration() time.Duration      { return 0 }
func (f *playbackFakeEngine) PositionAndDuration() (time.Duration, time.Duration) {
	return 0, 0
}
func (f *playbackFakeEngine) SetVolumeMin(float64)                   {}
func (f *playbackFakeEngine) VolumeMin() float64                     { return -50 }
func (f *playbackFakeEngine) SetVolume(float64)                      {}
func (f *playbackFakeEngine) Volume() float64                        { return 0 }
func (f *playbackFakeEngine) SetSpeed(float64)                       {}
func (f *playbackFakeEngine) Speed() float64                         { return 1 }
func (f *playbackFakeEngine) ToggleMono()                            {}
func (f *playbackFakeEngine) Mono() bool                             { return false }
func (f *playbackFakeEngine) SetEQBand(int, float64)                 {}
func (f *playbackFakeEngine) EQBands() [10]float64                   { return [10]float64{} }
func (f *playbackFakeEngine) StreamErr() error                       { return nil }
func (f *playbackFakeEngine) StreamTitle() string                    { return "" }
func (f *playbackFakeEngine) StreamBytes() (downloaded, total int64) { return 0, 0 }
func (f *playbackFakeEngine) SamplesInto([]float64) int              { return 0 }
func (f *playbackFakeEngine) SampleRate() int                        { return 44100 }

func TestNavTrackListQueueStartsQueuedTrackWhenStopped(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	p.Replace([]playlist.Track{
		{Title: "Existing", Path: "https://example.com/existing", Stream: true},
		{Title: "Other", Path: "https://example.com/other", Stream: true},
	})
	p.SetIndex(0)

	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
		navBrowser: navBrowserState{
			tracks: []playlist.Track{
				{Title: "Queued", Path: "https://example.com/queued", Stream: true},
			},
		},
	}

	cmd := m.handleNavTrackListKey(tea.KeyPressMsg{Text: "q"})
	if cmd == nil {
		t.Fatal("handleNavTrackListKey(q) = nil, want command")
	}
	if current, idx := m.playlist.Current(); current.Title != "Queued" || idx != 2 {
		t.Fatalf("current = (%q,%d), want (\"Queued\",2)", current.Title, idx)
	}
	if m.plCursor != 2 {
		t.Fatalf("plCursor = %d, want 2", m.plCursor)
	}
	if p.QueueLen() != 0 {
		t.Fatalf("QueueLen() = %d, want 0 after starting queued track", p.QueueLen())
	}
}

func TestTogglePlayPauseRestartsQueuedCurrentTrack(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	p.Replace([]playlist.Track{
		{Title: "Base", Path: "base.mp3", DurationSecs: 180},
		{Title: "Queued", Path: "queued.mp3", DurationSecs: 180},
	})
	p.SetIndex(0)
	p.Queue(1)
	if track, ok := p.Next(); !ok || track.Title != "Queued" {
		t.Fatalf("Next() = (%q,%t), want (\"Queued\",true)", track.Title, ok)
	}
	if !p.CurrentIsQueued() {
		t.Fatal("CurrentIsQueued() = false, want true")
	}

	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
	}

	if cmd := m.togglePlayPause(); cmd != nil {
		_ = cmd()
	}

	if len(player.playPaths) != 1 || player.playPaths[0] != "queued.mp3" {
		t.Fatalf("playPaths = %v, want [queued.mp3]", player.playPaths)
	}
	if current, idx := m.playlist.Current(); current.Title != "Queued" || idx != 1 {
		t.Fatalf("current = (%q,%d), want (\"Queued\",1)", current.Title, idx)
	}
}

func TestPlayCurrentTrackUnplayableUsesSelectionOrder(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	p.Replace([]playlist.Track{
		{Title: "Queued", Path: "https://example.com/queued", Stream: true},
		{Title: "Missing", Unplayable: true},
		{Title: "Replacement", Path: "https://example.com/replacement", Stream: true},
	})
	p.SetIndex(1)
	p.Queue(0)

	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
	}

	cmd := m.playCurrentTrack()
	if cmd == nil {
		t.Fatal("playCurrentTrack() = nil, want command")
	}
	if idx := m.playlist.Index(); idx != 2 {
		t.Fatalf("playlist.Index() = %d, want 2", idx)
	}
	if m.plCursor != 2 {
		t.Fatalf("plCursor = %d, want 2", m.plCursor)
	}
	if m.status.text != "Track unavailable, skipping..." {
		t.Fatalf("status.text = %q, want %q", m.status.text, "Track unavailable, skipping...")
	}
	if p.QueueLen() != 1 {
		t.Fatalf("QueueLen() = %d, want 1", p.QueueLen())
	}
}

func TestPlayCurrentTrackUnplayableStopsWhenNoReplacementExists(t *testing.T) {
	player := &playbackFakeEngine{playing: true}
	p := playlist.New()
	p.Replace([]playlist.Track{
		{Title: "Playing", Path: "playing.mp3", DurationSecs: 2},
		{Title: "Missing", Unplayable: true},
	})
	p.SetIndex(1)

	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
	}

	if cmd := m.playCurrentTrack(); cmd != nil {
		t.Fatalf("playCurrentTrack() = %v, want nil", cmd)
	}
	if len(player.playPaths) != 0 {
		t.Fatalf("playPaths = %v, want none", player.playPaths)
	}
	if player.IsPlaying() {
		t.Fatal("player.IsPlaying() = true, want false")
	}
	if _, idx := m.playlist.Current(); idx != 1 {
		t.Fatalf("current index = %d, want 1", idx)
	}
	if m.status.text != "No available tracks" {
		t.Fatalf("status.text = %q, want %q", m.status.text, "No available tracks")
	}
}

func modelAfterProviderPlaylistLoadWhilePlaying(t *testing.T) (Model, *playbackFakeEngine) {
	t.Helper()

	player := &playbackFakeEngine{playing: true}
	p := playlist.New()
	p.Replace([]playlist.Track{
		{Title: "Old", Path: "old.mp3", DurationSecs: 180},
	})
	p.SetIndex(0)

	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
	}

	updated, _ := m.Update(tracksLoadedMsg{
		{Title: "New 1", Path: "new1.mp3", DurationSecs: 180},
		{Title: "New 2", Path: "new2.mp3", DurationSecs: 180},
	})
	m = updated.(Model)

	return m, player
}

func TestProviderPlaylistLoadWhilePlayingKeepsNowPlayingTrack(t *testing.T) {
	m, player := modelAfterProviderPlaylistLoadWhilePlaying(t)

	track, idx := m.currentPlaybackTrack()
	if idx < 0 || track.Title != "Old" {
		t.Fatalf("currentPlaybackTrack() = (%q,%d), want old playing track", track.Title, idx)
	}
	if !m.playbackDetached {
		t.Fatal("playbackDetached = false, want true")
	}
	if player.clearPreloadCalls != 1 {
		t.Fatalf("ClearPreload calls = %d, want 1", player.clearPreloadCalls)
	}
	tracks := m.playlist.Tracks()
	if len(tracks) != 2 || tracks[0].Title != "New 1" || tracks[1].Title != "New 2" {
		t.Fatalf("playlist tracks = %#v, want new provider playlist only", tracks)
	}
	if info := m.renderTrackInfo(); !strings.Contains(info, "Old") {
		t.Fatalf("renderTrackInfo() = %q, want old playing track", info)
	}
}

func TestNextAfterProviderPlaylistLoadStartsFirstNewTrack(t *testing.T) {
	m, player := modelAfterProviderPlaylistLoadWhilePlaying(t)

	cmd := m.nextTrack()
	if cmd != nil {
		_ = cmd()
	}

	if len(player.playPaths) == 0 || player.playPaths[0] != "new1.mp3" {
		t.Fatalf("playPaths = %v, want first new track", player.playPaths)
	}
	if m.playbackDetached {
		t.Fatal("playbackDetached = true, want false")
	}
	track, _ := m.currentPlaybackTrack()
	if track.Title != "New 1" {
		t.Fatalf("currentPlaybackTrack() = %q, want New 1", track.Title)
	}
}

func TestPreloadAfterProviderPlaylistLoadUsesFirstNewTrack(t *testing.T) {
	m, player := modelAfterProviderPlaylistLoadWhilePlaying(t)

	cmd := m.preloadNext()
	if cmd == nil {
		t.Fatal("preloadNext() = nil, want preload command")
	}
	_ = cmd()

	if len(player.preloadPaths) != 1 || player.preloadPaths[0] != "new1.mp3" {
		t.Fatalf("preloadPaths = %v, want first new track", player.preloadPaths)
	}
}

func TestOfflinePreloadSkipsRemoteNextTracks(t *testing.T) {
	tests := []playlist.Track{
		{Title: "HTTP Stream", Path: "https://example.com/stream.mp3", Stream: true},
		{Title: "YT Search", Path: "ytsearch1:lofi"},
		{Title: "Spotify", Path: "spotify:track:abc"},
	}

	for _, next := range tests {
		t.Run(next.Title, func(t *testing.T) {
			player := &playbackFakeEngine{playing: true}
			p := playlist.New()
			p.Replace([]playlist.Track{
				{Title: "Now", Path: "now.mp3", DurationSecs: 180},
				next,
			})
			p.SetIndex(0)
			m := Model{
				player:   player,
				playlist: p,
				vis:      ui.NewVisualizer(float64(player.SampleRate())),
				offline:  true,
			}

			if cmd := m.preloadNext(); cmd != nil {
				_ = cmd()
				t.Fatalf("preloadNext() returned command in offline mode for %q", next.Path)
			}
			if len(player.preloadPaths) != 0 {
				t.Fatalf("preloadPaths = %v, want none", player.preloadPaths)
			}
		})
	}
}

func TestOfflineProviderPlaylistLoadRejectsRemoteTrack(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	p.Replace([]playlist.Track{{Title: "Existing", Path: "existing.mp3"}})
	m := Model{
		player:      player,
		playlist:    p,
		vis:         ui.NewVisualizer(float64(player.SampleRate())),
		offline:     true,
		provLoading: true,
	}

	updated, _ := m.Update(tracksLoadedMsg{
		{Title: "Remote", Path: "https://example.com/stream.mp3", Stream: true},
	})
	m = updated.(Model)

	tracks := m.playlist.Tracks()
	if len(tracks) != 1 || tracks[0].Path != "existing.mp3" {
		t.Fatalf("playlist tracks = %#v, want existing local playlist unchanged", tracks)
	}
	if m.provLoading {
		t.Fatal("provLoading = true, want false")
	}
	if !strings.Contains(m.status.text, "Offline mode: remote track") {
		t.Fatalf("status.text = %q, want offline rejection", m.status.text)
	}
	if player.clearPreloadCalls != 0 {
		t.Fatalf("ClearPreload calls = %d, want 0", player.clearPreloadCalls)
	}
}

func TestOfflineFileBrowserRejectsRemoteTracks(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
		offline:  true,
	}

	updated, _ := m.Update(fbTracksResolvedMsg{
		tracks: []playlist.Track{{Title: "Remote", Path: "https://example.com/stream.mp3", Stream: true}},
	})
	m = updated.(Model)

	if m.playlist.Len() != 0 {
		t.Fatalf("playlist.Len() = %d, want 0", m.playlist.Len())
	}
	if !strings.Contains(m.status.text, "Offline mode: remote track") {
		t.Fatalf("status.text = %q, want offline rejection", m.status.text)
	}
}

func TestOfflinePlaylistManagerPlayRejectsRemoteTracks(t *testing.T) {
	player := &playbackFakeEngine{}
	p := playlist.New()
	p.Replace([]playlist.Track{{Title: "Existing", Path: "existing.mp3"}})
	m := Model{
		player:   player,
		playlist: p,
		vis:      ui.NewVisualizer(float64(player.SampleRate())),
		offline:  true,
		plManager: plManagerState{
			selPlaylist: "mixed",
			tracks: []playlist.Track{
				{Title: "Remote", Path: "https://example.com/stream.mp3", Stream: true},
			},
		},
	}

	if cmd := m.plMgrLoadAndPlay(0); cmd != nil {
		t.Fatalf("plMgrLoadAndPlay() = %v, want nil", cmd)
	}
	tracks := m.playlist.Tracks()
	if len(tracks) != 1 || tracks[0].Path != "existing.mp3" {
		t.Fatalf("playlist tracks = %#v, want existing local playlist unchanged", tracks)
	}
	if len(player.playPaths) != 0 {
		t.Fatalf("playPaths = %v, want none", player.playPaths)
	}
	if !strings.Contains(m.status.text, "Offline mode: remote track") {
		t.Fatalf("status.text = %q, want offline rejection", m.status.text)
	}
}

func TestOfflineDirectTrackActionsRejectRemoteTracks(t *testing.T) {
	remote := playlist.Track{Title: "Remote", Path: "https://example.com/stream.mp3", Stream: true}
	tests := []struct {
		name string
		run  func(*Model) tea.Cmd
	}{
		{"play immediate", func(m *Model) tea.Cmd { return m.playTrackImmediate(remote) }},
		{"append", func(m *Model) tea.Cmd { return m.appendTrack(remote) }},
		{"queue next", func(m *Model) tea.Cmd { return m.queueTrackNext(remote) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &playbackFakeEngine{}
			p := playlist.New()
			p.Replace([]playlist.Track{{Title: "Existing", Path: "existing.mp3"}})
			m := Model{
				player:   player,
				playlist: p,
				vis:      ui.NewVisualizer(float64(player.SampleRate())),
				offline:  true,
			}

			if cmd := tt.run(&m); cmd != nil {
				t.Fatalf("%s returned command in offline mode", tt.name)
			}
			tracks := m.playlist.Tracks()
			if len(tracks) != 1 || tracks[0].Path != "existing.mp3" {
				t.Fatalf("playlist tracks = %#v, want existing local playlist unchanged", tracks)
			}
			if len(player.playPaths) != 0 {
				t.Fatalf("playPaths = %v, want none", player.playPaths)
			}
			if !strings.Contains(m.status.text, "Offline mode: remote track") {
				t.Fatalf("status.text = %q, want offline rejection", m.status.text)
			}
		})
	}
}
