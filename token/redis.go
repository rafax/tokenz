package token

import (
	"encoding/json"

	"github.com/twinj/uuid"

	redis "gopkg.in/redis.v5"
)

type redisHandler struct {
	store KeyValueStore
}

type KeyValueStore interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
}

type redisStore struct {
	client *redis.Client
}

func (h *redisStore) Set(key string, value []byte) error {
	return h.client.Set(key, value, 0).Err()
}

func (h *redisStore) Get(key string) ([]byte, error) {
	return h.client.Get(key).Bytes()
}

func (h *redisHandler) Encrypt(sd SubscriptionData) (Token, error) {
	s, err := json.Marshal(sd)
	if err != nil {
		return nil, err
	}
	t := uuid.NewV4().String()
	h.store.Set(t, s)
	return StringToken{Token: t}, nil
}

func (h *redisHandler) Decrypt(t Token) (SubscriptionData, error) {
	var v SubscriptionData
	bin, err := h.store.Get(t.String())
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bin, &v)
	return v, err
}

func NewRedisHandler() *redisHandler {
	return newRedisHandler(&redisStore{client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})})
}

func newRedisHandler(store KeyValueStore) *redisHandler {
	return &redisHandler{store: store}
}
