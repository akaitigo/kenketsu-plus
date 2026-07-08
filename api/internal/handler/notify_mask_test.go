package handler

import (
	"strings"
	"testing"
)

// TestMaskEndpoint verifies that the per-device subscription token in the URL
// path is redacted from log output while the push service host is retained (#21).
func TestMaskEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "fcm endpoint redacts token path",
			endpoint: "https://fcm.googleapis.com/fcm/send/abc123-secret-token",
			want:     "https://fcm.googleapis.com/[redacted]",
		},
		{
			name:     "mozilla endpoint redacts token path",
			endpoint: "https://updates.push.services.mozilla.com/wpush/v2/xyz-secret",
			want:     "https://updates.push.services.mozilla.com/[redacted]",
		},
		{
			name:     "endpoint with query redacts query",
			endpoint: "https://push.example.com/send?token=secret123",
			want:     "https://push.example.com/[redacted]",
		},
		{
			name:     "empty endpoint",
			endpoint: "",
			want:     "[redacted]",
		},
		{
			name:     "no scheme yields no host",
			endpoint: "fcm.googleapis.com/fcm/send/token",
			want:     "[redacted]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := maskEndpoint(tt.endpoint)
			if got != tt.want {
				t.Errorf("maskEndpoint(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
			if strings.Contains(got, "secret") || strings.Contains(got, "token") {
				t.Errorf("maskEndpoint(%q) leaked token material: %q", tt.endpoint, got)
			}
		})
	}
}
