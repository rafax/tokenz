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

func NewServer(bindTo string, handlers map[string]token.Handler) *Server {
	router := fasthttprouter.New()
	for k, h := range handlers {
		router.POST(fmt.Sprintf("/%v/:userId/:valid_seconds/:level/:platform", k), encode(h))
		router.GET(fmt.Sprintf("/%v/:token", k), decode(h))
	}
	return &Server{bindTo: bindTo, router: router}
}

func (s *Server) Start() {
	fmt.Printf("Listening on %s\n", s.bindTo)
	log.Fatal(fasthttp.ListenAndServe(s.bindTo, s.router.Handler))
}

func decode(h token.Handler) fasthttp.RequestHandler {
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

func encode(h token.Handler) fasthttp.RequestHandler {
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
