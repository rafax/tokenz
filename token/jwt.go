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
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload{Data: sd}).SignedString(mySigningKey)
	return StringToken{Token: ss}, err
}

func (h *JwtAssymetricHandler) Decrypt(t Token) (SubscriptionData, error) {
	token, err := jwt.ParseWithClaims(t.String(), &payload{}, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return SubscriptionData{}, err
	}
	if !token.Valid {
		return SubscriptionData{}, errors.New("Invalid token")
	}
	return token.Claims.(*payload).Data, nil
}

func NewJwtAssymetricHandler() *JwtAssymetricHandler {
	return &JwtAssymetricHandler{}
}
