package handler

import "testing"

func TestNotFound(t *testing.T) {
	h := NewMemoryHandler()
	_, err := h.Decrypt(StringToken{Token: "foo"})
	if err == nil {
		t.Error("Error should not be nil when decrypting nonexistent token")
	}

}
