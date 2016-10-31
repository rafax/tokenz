package token

import "time"

// SubscriptionData represents the data for a single user subscription
type SubscriptionData struct {
	UserId    string
	ExpiresAt time.Time
	Level     string
	Platform  string
}

// Equal returns true if a two instances of SubscriptionData are equal to each other
// It is necessary due to the fact that ExpiresAt cannot be reliably compared with ==
func (sd SubscriptionData) Equal(u SubscriptionData) bool {
	return sd.UserId == u.UserId && sd.ExpiresAt.Equal(u.ExpiresAt) && sd.Level == u.Level && sd.Platform == u.Platform
}

// Token represents a token issued to the user
type Token interface {
	String() string
}

// StringToken represents the Token as a string
type StringToken struct {
	Token string
}

func (s StringToken) String() string {
	return s.Token
}

// TokenHandler converts between SuscriptionData and Token
type TokenHandler interface {
	Encrypt(SubscriptionData) (Token, error)
	Decrypt(Token) (SubscriptionData, error)
}
