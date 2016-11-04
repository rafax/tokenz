package token

import (
	"errors"
	"sync"

	"github.com/satori/go.uuid"
)

type MemoryHandler struct {
	data map[string]SubscriptionData
	lock sync.RWMutex
}

func (h *MemoryHandler) Encrypt(sd SubscriptionData) (Token, error) {
	token := uuid.NewV4()
	h.lock.Lock()
	h.data[token.String()] = sd
	h.lock.Unlock()
	return token, nil
}

func (h *MemoryHandler) Decrypt(t Token) (SubscriptionData, error) {
	h.lock.RLock()
	sd, ok := h.data[t.String()]
	h.lock.RUnlock()
	if !ok {
		return SubscriptionData{}, errors.New("No subscription found for token")
	}
	return sd, nil
}

func NewMemoryHandler() *MemoryHandler {
	return &MemoryHandler{data: make(map[string]SubscriptionData)}
}
