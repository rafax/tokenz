package token

import (
	"testing"
	"time"
)

func TestBase64(t *testing.T) {
	testHandler(t, NewBase64Handler())
}

func TestMemory(t *testing.T) {
	testHandler(t, NewMemoryHandler())
}

func testHandler(t *testing.T, h Handler) {
	expected := SubscriptionData{
		UserId:    "uid",
		ExpiresAt: time.Now().Add(time.Second),
		Level:     "all",
		Platform:  "mobile",
	}
	token, err := h.Encrypt(expected)
	if err != nil {
		t.Errorf("Encrypt failed %v", err)
	}
	sd, err := h.Decrypt(token)
	if err != nil {
		t.Errorf("Decrypt failed %v", err)
	}
	if !sd.Equal(expected) {
		t.Errorf("%v != %v, they should be equal", sd, expected)
	}
}
