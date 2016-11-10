// +build e2e

package token

import (
	"testing"

	"github.com/twinj/uuid"

	redis "gopkg.in/redis.v5"
)

func TestRedisStore(t *testing.T) {
	k, v := uuid.NewV4().String(), []byte(uuid.NewV4().String())
	s := &redisStore{client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
	err := s.Set(k, v)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	got, err := s.Get(k)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	if string(got) != string(v) {
		t.Errorf("Got %v, expected %v", got, v)
	}
}
