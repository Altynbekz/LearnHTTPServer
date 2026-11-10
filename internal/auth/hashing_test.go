package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	tokenString, err := MakeJWT(userID, secret, 2*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		t.Fatalf("invalid signed token: %v", err)
	}
	if claims.Issuer != "chirpy-access" || claims.Subject != userID.String() {
		t.Fatalf("unexpected claims: issuer=%q subject=%q", claims.Issuer, claims.Subject)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil || claims.ExpiresAt.Sub(claims.IssuedAt.Time) != 2*time.Minute {
		t.Fatalf("unexpected token lifetime: issued=%v expires=%v", claims.IssuedAt, claims.ExpiresAt)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	valid, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ValidateJWT(valid, secret)
	if err != nil || got != userID {
		t.Fatalf("valid token: got %s, err %v", got, err)
	}

	expired, err := MakeJWT(userID, secret, -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	badSubject, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "not-a-uuid",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name   string
		token  string
		secret string
	}{
		{name: "wrong signing key", token: valid, secret: "wrong-secret"},
		{name: "expired", token: expired, secret: secret},
		{name: "invalid subject", token: badSubject, secret: secret},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ValidateJWT(tc.token, tc.secret); err == nil {
				t.Fatal("expected token validation error")
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		want          string
		wantErr       bool
	}{
		{name: "valid token", authorization: "Bearer secret-token", want: "secret-token"},
		{name: "token ending in r", authorization: "Bearer server", want: "server"},
		{name: "missing header", wantErr: true},
		{name: "missing Bearer scheme", authorization: "secret-token", wantErr: true},
		{name: "missing token", authorization: "Bearer ", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := http.Header{}
			if tc.authorization != "" {
				headers.Set("Authorization", tc.authorization)
			}
			got, err := GetBearerToken(headers)
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("GetBearerToken() = %q, want %q", got, tc.want)
			}
		})
	}
}
