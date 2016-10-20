package handler

import "time"

type SubscriptionData struct {
	UserId    string
	ExpiresAt time.Time
	Level     string
	Platform  string
}

func (sd SubscriptionData) Equal(u SubscriptionData) bool {
	return sd.UserId == u.UserId && sd.ExpiresAt.Equal(u.ExpiresAt) && sd.Level == u.Level && sd.Platform == u.Platform
}

type Token interface {
	String() string
}

type StringToken struct {
	Token string
}

func (s StringToken) String() string {
	return s.Token
}

type TokenHandler interface {
	Encrypt(SubscriptionData) (Token, error)
	Decrypt(Token) (SubscriptionData, error)
}
