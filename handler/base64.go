package handler

import (
	"encoding/base64"
	"encoding/json"
)

type base64Handler struct {
}

func (b base64Handler) Encrypt(sd SubscriptionData) (Token, error) {
	s, err := json.Marshal(sd)
	if err != nil {
		return nil, err
	}
	return StringToken{Token: base64.StdEncoding.EncodeToString(s)}, nil
}

func (b base64Handler) Decrypt(t Token) (SubscriptionData, error) {
	var v SubscriptionData
	s, err := base64.StdEncoding.DecodeString(t.String())
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(s, &v)
	return v, err
}

func NewBase64Handler() base64Handler {
	return base64Handler{}
}
