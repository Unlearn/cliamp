package model

import "cliamp/playlist"

func (m *Model) rejectOfflineRemoteTrack(track playlist.Track) bool {
	if !m.offline || !playlist.IsRemotePlaybackTrack(track) {
		return false
	}
	m.status.Showf(statusTTLDefault, "Offline mode: remote track %q is disabled", track.Path)
	return true
}

func (m *Model) rejectOfflineRemoteTracks(tracks []playlist.Track) bool {
	if !m.offline {
		return false
	}
	if path := playlist.FirstRemotePlaybackPath(tracks); path != "" {
		m.status.Showf(statusTTLDefault, "Offline mode: remote track %q is disabled", path)
		return true
	}
	return false
}
