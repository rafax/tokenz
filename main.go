package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	handler TokenHandler = Base64Handler{}
)

func Decode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	t := ps.ByName("token")
	sd, _ := handler.Decrypt(Token{string: t})
	j, _ := json.Marshal(sd)
	fmt.Fprint(ctx, string(j))
}

func Encode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	validForSeconds, _ := strconv.Atoi(ps.ByName("valid_seconds"))
	sd := SubscriptionData{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(validForSeconds)),
		UserId:    ps.ByName("userId"),
		Platform:  ps.ByName("platform"),
		Level:     ps.ByName("level"),
	}
	t, _ := handler.Encrypt(sd)
	fmt.Fprintf(ctx, "{\"token\": %s}", t.string)
}

func main() {
	router := fasthttprouter.New()
	router.POST("/b64/:userId/:valid_seconds/:level/:platform", Encode)
	router.GET("/b64/:token", Decode)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

type TokenHandler interface {
	Encrypt(SubscriptionData) (Token, error)
	Decrypt(Token) (SubscriptionData, error)
}

type Base64Handler struct {
}

func (b Base64Handler) Encrypt(sd SubscriptionData) (Token, error) {
	s, err := json.Marshal(sd)
	if err != nil {
		return Token{}, err
	}
	return Token{base64.StdEncoding.EncodeToString(s)}, nil
}

func (b Base64Handler) Decrypt(t Token) (SubscriptionData, error) {
	s, err := base64.StdEncoding.DecodeString(t.string)
	var v SubscriptionData
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(s, &v)
	if err != nil {
		return v, err
	}
	return v, nil
}

type SubscriptionData struct {
	UserId    string
	ExpiresAt time.Time
	Level     string
	Platform  string
}

type Token struct {
	string
}
