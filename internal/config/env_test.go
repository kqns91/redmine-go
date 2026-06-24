package config

import "testing"

const (
	canonicalURL = "https://canonical.example.com"
	aliasURL     = "https://alias.example.com"
)

func TestURLFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		urlAlias string
		want     string
	}{
		{
			name:     "canonical only",
			url:      canonicalURL,
			urlAlias: "",
			want:     canonicalURL,
		},
		{
			name:     "alias only",
			url:      "",
			urlAlias: aliasURL,
			want:     aliasURL,
		},
		{
			name:     "canonical takes precedence over alias",
			url:      canonicalURL,
			urlAlias: aliasURL,
			want:     canonicalURL,
		},
		{
			name:     "neither set",
			url:      "",
			urlAlias: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvURL, tt.url)
			t.Setenv(EnvURLAlias, tt.urlAlias)

			if got := URLFromEnv(); got != tt.want {
				t.Errorf("URLFromEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIKeyFromEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "secret-key")

	if got := APIKeyFromEnv(); got != "secret-key" {
		t.Errorf("APIKeyFromEnv() = %q, want %q", got, "secret-key")
	}
}
