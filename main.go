package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/rafax/tokenz/handler"
	"github.com/valyala/fasthttp"
)

var (
	h handler.TokenHandler
)

func decode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	t := ps.ByName("token")
	sd, err := h.Decrypt(handler.StringToken{Token: t})
	if err != nil {
		ctx.Error(fmt.Sprintf("Error when decrypting token: %s", err), 500)
	}
	j, err := json.Marshal(sd)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error when marshalling response: %s", err), 500)
	}
	fmt.Fprint(ctx, string(j))
}

func encode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	validForSeconds, err := strconv.Atoi(ps.ByName("valid_seconds"))
	if err != nil {
		ctx.Error(fmt.Sprintf("Error when parsing valid_seconds: %s", err), 500)
	}
	sd := handler.SubscriptionData{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(validForSeconds)),
		UserId:    ps.ByName("userId"),
		Platform:  ps.ByName("platform"),
		Level:     ps.ByName("level"),
	}
	log.Println(sd)
	t, err := h.Encrypt(sd)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error when encrypting to token: %s", err), 500)
	}
	fmt.Fprintf(ctx, "{\"token\": %s}", t.String())
}

func main() {
	h = handler.NewBase64Handler()

	router := fasthttprouter.New()
	router.POST("/b64/:userId/:valid_seconds/:level/:platform", encode)
	router.GET("/b64/:token", decode)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
