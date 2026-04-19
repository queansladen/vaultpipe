package config

import "testing"

func TestResolveQuote_NilConfig(t *testing.T) {
	got := ResolveQuote(nil)
	if got.Enabled {
		t.Fatal("expected Enabled=false for nil config")
	}
}

func TestResolveQuote_ExplicitDisabled(t *testing.T) {
	got := ResolveQuote(&QuoteConfig{Enabled: false})
	if got.Enabled {
		t.Fatal("expected Enabled=false")
	}
}

func TestResolveQuote_ExplicitEnabled(t *testing.T) {
	got := ResolveQuote(&QuoteConfig{Enabled: true})
	if !got.Enabled {
		t.Fatal("expected Enabled=true")
	}
}
