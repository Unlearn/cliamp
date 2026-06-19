package model

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"cliamp/playlist"
)

type initFakeProvider struct{}

func (initFakeProvider) Name() string { return "Radio" }
func (initFakeProvider) Playlists() ([]playlist.PlaylistInfo, error) {
	return []playlist.PlaylistInfo{{ID: "1", Name: "One"}}, nil
}
func (initFakeProvider) Tracks(string) ([]playlist.Track, error) { return nil, nil }

func TestInitFetchesProviderOnlyWhenStartingInProviderFocus(t *testing.T) {
	player := &playbackFakeEngine{}
	pl := playlist.New()
	pl.Add(playlist.Track{Title: "Local", Path: "local.mp3"})
	prov := initFakeProvider{}
	m := New(player, pl, []ProviderEntry{{Key: "radio", Name: "Radio", Provider: prov}}, "radio", nil, nil, nil, nil)

	cmd := m.Init()
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("Init() = %T, want tea.BatchMsg", cmd())
	}
	if len(batch) != 2 {
		t.Fatalf("Init() command count = %d, want 2 without provider fetch", len(batch))
	}

	m.StartInProvider()
	cmd = m.Init()
	batch, ok = cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("Init() after StartInProvider = %T, want tea.BatchMsg", cmd())
	}
	if len(batch) != 3 {
		t.Fatalf("Init() command count = %d, want 3 with provider fetch", len(batch))
	}
}

func TestFocusProviderPaneFetchesOnce(t *testing.T) {
	player := &playbackFakeEngine{}
	pl := playlist.New()
	prov := initFakeProvider{}
	m := New(player, pl, []ProviderEntry{{Key: "radio", Name: "Radio", Provider: prov}}, "radio", nil, nil, nil, nil)

	cmd := m.focusProviderPane()
	if cmd == nil {
		t.Fatal("focusProviderPane() = nil, want provider fetch command")
	}
	if m.focus != focusProvider {
		t.Fatalf("focus = %v, want focusProvider", m.focus)
	}
	if !m.provLoading {
		t.Fatal("provLoading = false, want true")
	}

	if cmd := m.focusProviderPane(); cmd != nil {
		t.Fatal("focusProviderPane() while loading returned command, want nil")
	}

	m.provLoading = false
	m.providerLists = []playlist.PlaylistInfo{{ID: "cached", Name: "Cached"}}
	if cmd := m.focusProviderPane(); cmd != nil {
		t.Fatal("focusProviderPane() with cached lists returned command, want nil")
	}
}
