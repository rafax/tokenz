package token

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtAssymetricHandler struct {
}

var mySigningKey []byte = []byte("AllYourBase")

type payload struct {
	Data SubscriptionData
	jwt.StandardClaims
}

func (h *JwtAssymetricHandler) Encrypt(sd SubscriptionData) (Token, error) {
	jwt.NewWithClaims(jwt.SigningMethodHS256, payload{Data: sd})
	ss, err := token.SignedString(mySigningKey)
	return StringToken{Token: ss}, err
}

func (h *JwtAssymetricHandler) Decrypt(t Token) (SubscriptionData, error) {
	sd, ok := h.data[t.String()]
	if !ok {
		return SubscriptionData{}, errors.New("No subscription found for token")
	}
	return sd, nil
}

func NewMemoryHandler() *JwtAssymetricHandler {
	return &JwtAssymetricHandler{data: make(map[string]SubscriptionData)}
}
