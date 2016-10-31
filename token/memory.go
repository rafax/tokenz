package token

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
)

type MemoryHandler struct {
	data map[string]SubscriptionData
}

func (h *MemoryHandler) Encrypt(sd SubscriptionData) (Token, error) {
	token := uuid.NewV4()
	h.data[token.String()] = sd
	return token, nil
}

func (h *MemoryHandler) Decrypt(t Token) (SubscriptionData, error) {
	fmt.Printf("Looking for %v in %v", t, h.data)
	sd, ok := h.data[t.String()]
	if !ok {
		return SubscriptionData{}, errors.New("Token invalid")
	}
	return sd, nil
}

func NewMemoryHandler() *MemoryHandler {
	return &MemoryHandler{data: make(map[string]SubscriptionData)}
}
