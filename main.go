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
	b64h   handler.TokenHandler
	mh     *handler.MemoryHandler
	bindTo string = ":8080"
)

func decode(h handler.TokenHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := ctx.UserValue("token").(string)
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
}

func encode(h handler.TokenHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		validForSeconds, err := strconv.Atoi(ctx.UserValue("valid_seconds").(string))
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when parsing valid_seconds: %s", err), 500)
		}
		sd := handler.SubscriptionData{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(validForSeconds)),
			UserId:    ctx.UserValue("userId").(string),
			Platform:  ctx.UserValue("platform").(string),
			Level:     ctx.UserValue("level").(string),
		}
		log.Println(sd)
		t, err := h.Encrypt(sd)
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when encrypting to token: %s", err), 500)
		}
		fmt.Fprintf(ctx, "{\"token\": \"%s\"}", t.String())
	}
}

func main() {
	b64h = handler.NewBase64Handler()
	mh = handler.NewMemoryHandler()

	router := fasthttprouter.New()
	router.POST("/b64/:userId/:valid_seconds/:level/:platform", encode(b64h))
	router.GET("/b64/:token", decode(b64h))
	router.POST("/mem/:userId/:valid_seconds/:level/:platform", encode(mh))
	router.GET("/mem/:token", decode(mh))

	fmt.Printf("Listening on %s\n", bindTo)

	log.Fatal(fasthttp.ListenAndServe(bindTo, router.Handler))
}
