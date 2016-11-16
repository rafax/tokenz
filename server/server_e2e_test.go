// +build e2e

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/rafax/tokenz/token"
)

const userID string = "userid"

func TestRoundtrip(t *testing.T) {
	bindTo := ":8888"
	s := NewServer(":8888", map[string]token.Handler{"b64": token.NewBase64Handler(), "mem": token.NewMemoryHandler(), "red": token.NewRedisHandler(), "ajw": token.NewJwtAssymetricHandler()})
	go s.Start()
	verifyStarted("http://" + bindTo)
	t.Run("Test base64", handlerTest("b64"))
	t.Run("Test redis", handlerTest("red"))
	t.Run("Test memory", handlerTest("mem"))
	t.Run("Test assymetric jwt", handlerTest("ajw"))
}

func handlerTest(method string) func(t *testing.T) {
	var token *string
	return func(t *testing.T) {
		t.Run("Fetch token", func(t *testing.T) {
			resp, err := http.Post(fmt.Sprintf("http://:8888/%v/%v/1000/all/mobilez", method, userID), "application/json", nil)
			data, err := getJSON(resp, err, t)
			if err != nil {
				t.Errorf("Could not get the token from response: %v\n", err)
				return
			}
			tkn, ok := data.Path("token").Data().(string)
			if !ok {
				t.Errorf("Token not found in response when getting token: %v", data.String())
				return
			}
			token = &tkn
		})
		if token == nil {
			t.Error("Cannot verify data, token is nil")
			return
		}
		t.Run("Fetch data for token", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://:8888/%v/%v", method, *token))
			data, err := getJSON(resp, err, t)
			if err != nil {
				t.Errorf("Could not get subscription data from response: %v", err)
				return
			}
			uid, ok := data.Path("UserId").Data().(string)
			if !ok {
				t.Errorf("UserId not found in response: %v", data.String())
				return
			}
			if uid != userID {
				t.Errorf("Invalid userId %v expected %v", uid, userID)
				return
			}
		})
	}
}

func getJSON(resp *http.Response, err error, t *testing.T) (*gabs.Container, error) {
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func verifyStarted(bindTo string) {
	delay := 1 * time.Nanosecond
	for i := 0; i < 10; i++ {
		_, err := http.Get(bindTo)
		if err == nil {
			return
		}
		fmt.Println(err)
		delay = 10 * delay
		time.Sleep(delay)
	}
}
