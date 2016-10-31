package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/rafax/tokenz/token"
	"github.com/valyala/fasthttp"
)

type Server struct {
	bindTo string
	router *fasthttprouter.Router
}

func NewServer(b64Handler, memHandler token.TokenHandler, bindTo string) *Server {
	router := fasthttprouter.New()
	router.POST("/b64/:userId/:valid_seconds/:level/:platform", encode(b64Handler))
	router.GET("/b64/:token", decode(b64Handler))
	router.POST("/mem/:userId/:valid_seconds/:level/:platform", encode(memHandler))
	router.GET("/mem/:token", decode(memHandler))

	return &Server{bindTo: bindTo, router: router}
}

func (s *Server) Start() {
	fmt.Printf("Listening on %s\n", s.bindTo)
	log.Fatal(fasthttp.ListenAndServe(s.bindTo, s.router.Handler))
}

func decode(h token.TokenHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := ctx.UserValue("token").(string)
		sd, err := h.Decrypt(token.StringToken{Token: t})
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when decrypting token: %s", err), 500)
			return
		}
		j, err := json.Marshal(sd)
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when marshalling response: %s", err), 500)
			return
		}
		fmt.Fprint(ctx, string(j))
	}
}

func encode(h token.TokenHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		validForSeconds, err := strconv.Atoi(ctx.UserValue("valid_seconds").(string))
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when parsing valid_seconds: %s", err), 500)
			return
		}
		sd := token.SubscriptionData{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(validForSeconds)),
			UserId:    ctx.UserValue("userId").(string),
			Platform:  ctx.UserValue("platform").(string),
			Level:     ctx.UserValue("level").(string),
		}
		log.Println(sd)
		t, err := h.Encrypt(sd)
		if err != nil {
			ctx.Error(fmt.Sprintf("Error when encrypting to token: %s", err), 500)
			return
		}
		fmt.Fprintf(ctx, "{\"token\": \"%s\"}", t.String())
	}
}