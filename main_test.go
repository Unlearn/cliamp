package main

import "testing"

func TestSelectDefaultProvider(t *testing.T) {
	tests := []struct {
		name                string
		configured          string
		offline             bool
		providerOverrideSet bool
		want                string
		wantErr             bool
	}{
		{name: "default no args uses radio", want: "radio"},
		{name: "empty config keeps radio as active provider", want: "radio"},
		{name: "configured provider wins", configured: "navidrome", want: "navidrome"},
		{name: "provider override keeps radio", configured: "radio", providerOverrideSet: true, want: "radio"},
		{name: "offline ignores configured provider without override", configured: "radio", offline: true, want: "local"},
		{name: "offline accepts local override", configured: "local", offline: true, providerOverrideSet: true, want: "local"},
		{name: "offline rejects remote override", configured: "radio", offline: true, providerOverrideSet: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectDefaultProvider(tt.configured, tt.offline, tt.providerOverrideSet)
			if (err != nil) != tt.wantErr {
				t.Fatalf("selectDefaultProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("selectDefaultProvider() = %q, want %q", got, tt.want)
			}
		})
	}
}
