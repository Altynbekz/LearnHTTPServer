package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateChirpRejectsInvalidBody(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{`},
		{name: "over 140 characters", body: `{"body":"` + strings.Repeat("a", 141) + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/chirps", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			newMux().ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status %d, want 400; body %q", w.Code, w.Body.String())
			}
		})
	}
}

func TestAPIRoutes(t *testing.T) {
	mux := newMux()
	request := func(method, path string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, "http://example.test"+path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		return w
	}

	if got := request(http.MethodGet, "/api/healthz"); got.Code != http.StatusOK || got.Body.String() != "OK" {
		t.Fatalf("healthz: status %d, body %q", got.Code, got.Body.String())
	}
	if got := request(http.MethodGet, "/api/metrics"); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "Chirpy has been visited 0 times!") {
		t.Fatalf("metrics: status %d, body %q", got.Code, got.Body.String())
	}
	// The reset route exists, but the default test config is not in dev mode.
	if got := request(http.MethodPost, "/admin/reset"); got.Code != http.StatusForbidden {
		t.Fatalf("reset: status %d, want 403", got.Code)
	}
	for _, route := range []struct{ method, path string }{
		{http.MethodGet, "/healthz"},
		{http.MethodGet, "/metrics"},
		{http.MethodPost, "/reset"},
	} {
		if got := request(route.method, route.path); got.Code != http.StatusNotFound {
			t.Errorf("old route %s %s: status %d, want 404", route.method, route.path, got.Code)
		}
	}
	if got := request(http.MethodGet, "/about-me"); got.Code != http.StatusOK {
		t.Errorf("website route: status %d, want 200", got.Code)
	}
}
