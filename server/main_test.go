package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppVersion(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		want   string
	}{
		{"unset", "", "dev"},
		{"full SHA is shortened", "1f94d4169388e2f6a1b0c3d4e5f60718293a4b5c", "1f94d41"},
		{"short value kept", "abc123", "abc123"},
		{"whitespace trimmed", " 1f94d4169388\n", "1f94d41"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SOURCE_COMMIT", tt.commit)

			if got := appVersion(); got != tt.want {
				t.Errorf("appVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	handleHealth("1f94d41")(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding body %q: %v", rec.Body.String(), err)
	}

	if body.Status != "ok" || body.Version != "1f94d41" {
		t.Errorf("body = %+v, want status ok and version 1f94d41", body)
	}
}
