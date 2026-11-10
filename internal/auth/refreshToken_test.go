package auth

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	first := MakeRefreshToken()
	if len(first) != 64 {
		t.Fatalf("token length = %d, want 64 hex characters", len(first))
	}
	log.Println(first)
	decoded, err := hex.DecodeString(first)
	if err != nil {
		t.Fatalf("token is not valid hex: %v", err)
	}
	if len(decoded) != 32 {
		t.Fatalf("decoded token length = %d, want 32 bytes", len(decoded))
	}
	log.Println(decoded)
	if second := MakeRefreshToken(); second == first {
		t.Fatal("two refresh tokens are identical")
	}
}
