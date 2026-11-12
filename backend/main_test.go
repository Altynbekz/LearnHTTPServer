package main

import (
	"httpServerWithGo/internal/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAccessTokenDuration(t *testing.T) {
	for _, tc := range []struct {
		seconds int
		want    time.Duration
	}{
		{seconds: 0, want: time.Hour},
		{seconds: -1, want: time.Hour},
		{seconds: 30, want: 30 * time.Second},
		{seconds: 7200, want: time.Hour},
	} {
		if got := accessTokenDuration(tc.seconds); got != tc.want {
			t.Errorf("accessTokenDuration(%d) = %s, want %s", tc.seconds, got, tc.want)
		}
	}
}

func TestCreateChirpRejectsInvalidBody(t *testing.T) {
	const secret = "test-secret"
	token, err := auth.MakeJWT(uuid.New(), secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	mux := newMuxWithConfig(&apiConfig{secretKey: secret})
	for _, tc := range []struct {
		name        string
		body        string
		bearerToken bool
		wantError   string
	}{
		{name: "malformed JSON", body: `{`, wantError: "Invalid JSON request"},
		{name: "over 140 characters", body: `{"body":"` + strings.Repeat("a", 141) + `"}`, bearerToken: true, wantError: "chirp is too long"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/chirps", strings.NewReader(tc.body))
			if tc.bearerToken {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status %d, want 400; body %q", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tc.wantError) {
				t.Fatalf("body %q does not contain %q", w.Body.String(), tc.wantError)
			}
		})
	}
}

func TestCreateChirpRejectsInvalidToken(t *testing.T) {
	mux := newMuxWithConfig(&apiConfig{secretKey: "test-secret"})
	for _, authorization := range []string{"", "Bearer not-a-jwt"} {
		req := httptest.NewRequest(http.MethodPost, "/api/chirps", strings.NewReader(`{"body":"hello"}`))
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Authorization %q: status %d, want 401", authorization, w.Code)
		}
	}
}

func TestGetChirpRejectsInvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/chirps/not-a-uuid", nil)
	w := httptest.NewRecorder()
	newMux().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400; body %q", w.Code, w.Body.String())
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
