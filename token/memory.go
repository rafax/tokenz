package token

import (
	"errors"

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
	sd, ok := h.data[t.String()]
	if !ok {
		return SubscriptionData{}, errors.New("No subscription found for token")
	}
	return sd, nil
}

func NewMemoryHandler() *MemoryHandler {
	return &MemoryHandler{data: make(map[string]SubscriptionData)}
}
