package token

import (
	"encoding/json"

	"github.com/twinj/uuid"

	redis "gopkg.in/redis.v5"
)

type redisHandler struct {
	client *redis.Client
}

func (h *redisHandler) Encrypt(sd SubscriptionData) (Token, error) {
	s, err := json.Marshal(sd)
	if err != nil {
		return nil, err
	}
	t := uuid.NewV4().String()
	h.client.Set(t, s, 0)
	return StringToken{Token: t}, nil
}

func (h *redisHandler) Decrypt(t Token) (SubscriptionData, error) {
	var v SubscriptionData
	bin, err := h.client.Get(t.String()).Bytes()
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bin, &v)
	return v, err
}

func NewRedisHandler() *redisHandler {
	return &redisHandler{client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
}
