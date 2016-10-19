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
	h handler.TokenHandler = handler.NewBase64Handler()
)

func Decode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	t := ps.ByName("token")
	sd, _ := h.Decrypt(handler.StringToken{Token: t})
	j, _ := json.Marshal(sd)
	fmt.Fprint(ctx, string(j))
}

func Encode(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	validForSeconds, _ := strconv.Atoi(ps.ByName("valid_seconds"))
	sd := handler.SubscriptionData{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(validForSeconds)),
		UserId:    ps.ByName("userId"),
		Platform:  ps.ByName("platform"),
		Level:     ps.ByName("level"),
	}
	log.Println(sd)
	t, _ := h.Encrypt(sd)
	fmt.Fprintf(ctx, "{\"token\": %s}", t.String())
}

func main() {
	router := fasthttprouter.New()
	router.POST("/b64/:userId/:valid_seconds/:level/:platform", Encode)
	router.GET("/b64/:token", Decode)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
