package journal

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestEncryptDecryptJSONRoundTrip(t *testing.T) {
	key := mustKey(t)
	plaintext := []byte(`{"entry_text":"I felt steady today."}`)

	nonce1, ciphertext1, err := encryptJSON(key, plaintext)
	if err != nil {
		t.Fatalf("encryptJSON() error = %v", err)
	}
	nonce2, ciphertext2, err := encryptJSON(key, plaintext)
	if err != nil {
		t.Fatalf("encryptJSON() second call error = %v", err)
	}
	if nonce1 == nonce2 {
		t.Fatal("expected unique nonces across encryptions")
	}
	if ciphertext1 == ciphertext2 {
		t.Fatal("expected ciphertext to differ when nonce differs")
	}
	if strings.Contains(ciphertext1, "steady") {
		t.Fatalf("ciphertext should not contain plaintext, got %q", ciphertext1)
	}

	decrypted, err := decryptJSON(key, nonce1, ciphertext1)
	if err != nil {
		t.Fatalf("decryptJSON() error = %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected decrypted plaintext %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptJSONRejectsInvalidKey(t *testing.T) {
	if _, _, err := encryptJSON([]byte("short-key"), []byte("{}")); err == nil {
		t.Fatal("expected encryptJSON to reject invalid key length")
	}
}

func TestDecryptJSONRejectsInvalidInputs(t *testing.T) {
	key := mustKey(t)
	nonce, ciphertext, err := encryptJSON(key, []byte(`{"ok":true}`))
	if err != nil {
		t.Fatalf("encryptJSON() error = %v", err)
	}

	tests := []struct {
		name       string
		key        []byte
		nonce      string
		ciphertext string
	}{
		{
			name:       "invalid key",
			key:        []byte("short-key"),
			nonce:      nonce,
			ciphertext: ciphertext,
		},
		{
			name:       "invalid nonce base64",
			key:        key,
			nonce:      "not-base64!",
			ciphertext: ciphertext,
		},
		{
			name:       "wrong nonce length",
			key:        key,
			nonce:      base64.StdEncoding.EncodeToString([]byte("short")),
			ciphertext: ciphertext,
		},
		{
			name:       "invalid ciphertext base64",
			key:        key,
			nonce:      nonce,
			ciphertext: "not-base64!",
		},
		{
			name:       "tampered ciphertext",
			key:        key,
			nonce:      nonce,
			ciphertext: base64.StdEncoding.EncodeToString([]byte("tampered")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := decryptJSON(tt.key, tt.nonce, tt.ciphertext); err == nil {
				t.Fatal("expected decryptJSON to fail")
			}
		})
	}
}
