package handler

import "time"

type SubscriptionData struct {
	UserId    string
	ExpiresAt time.Time
	Level     string
	Platform  string
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
