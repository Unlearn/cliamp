package ui

import (
	"testing"
	"time"
)

func TestParseJumpTarget(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    time.Duration
		wantErr bool
	}{
		{name: "seconds only", in: "10", want: 10 * time.Second},
		{name: "minutes seconds", in: "58:05", want: 58*time.Minute + 5*time.Second},
		{name: "minutes one-digit seconds", in: "58:6", want: 58*time.Minute + 6*time.Second},
		{name: "minutes trailing colon", in: "58:", want: 58 * time.Minute},
		{name: "missing minutes accepted", in: ":49", want: 49 * time.Second},
		{name: "spaces trimmed", in: "  12:3  ", want: 12*time.Minute + 3*time.Second},
		{name: "empty", in: "", wantErr: true},
		{name: "not number", in: "abc", wantErr: true},
		{name: "bad minutes", in: "x:05", wantErr: true},
		{name: "bad seconds", in: "10:x", wantErr: true},
		{name: "seconds too large", in: "10:60", wantErr: true},
		{name: "high minute high second", in: "99:99", wantErr: true},
		{name: "too many colons", in: "1:2:3", wantErr: true},
		{name: "too many second digits", in: "10:123", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJumpTarget(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestFormatJumpClock(t *testing.T) {
	tests := []struct {
		in   time.Duration
		want string
	}{
		{in: 0, want: "00:00"},
		{in: 10 * time.Second, want: "00:10"},
		{in: 58*time.Minute + 5*time.Second, want: "58:05"},
		{in: 75*time.Minute + 48*time.Second, want: "75:48"},
	}

	for _, tt := range tests {
		if got := formatJumpClock(tt.in); got != tt.want {
			t.Fatalf("formatJumpClock(%v) = %q want %q", tt.in, got, tt.want)
		}
	}
}
