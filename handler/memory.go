package handler

import (
	"errors"

	"github.com/satori/go.uuid"
)

type MemoryHandler struct {
	data map[Token]SubscriptionData
}

func (h *MemoryHandler) Encrypt(sd SubscriptionData) (Token, error) {
	token := uuid.NewV4()
	h.data[token] = sd
	return token, nil
}

func (h *MemoryHandler) Decrypt(t Token) (SubscriptionData, error) {
	sd, ok := h.data[t]
	if !ok {
		return SubscriptionData{}, errors.New("No subscription found for token")
	}
	return sd, nil
}

func NewMemoryHandler() *MemoryHandler {
	return &MemoryHandler{data: make(map[Token]SubscriptionData)}
}
