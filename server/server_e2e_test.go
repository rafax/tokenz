// +build e2e

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/rafax/tokenz/token"
)

const userID string = "userid"

func TestRoundtrip(t *testing.T) {
	s := NewServer(":8888", map[string]token.Handler{"b64": token.NewBase64Handler(), "mem": token.NewMemoryHandler()})
	go s.Start()
	var token *string
	t.Run("Test memory", handlerTest("mem", token))
	t.Run("Test base64", handlerTest("b64", token))
	//t.Run("Test redis", handlerTest("redis", token))
}

func handlerTest(method string, token *string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Fetch token", func(t *testing.T) {
			resp, err := http.Post(fmt.Sprintf("http://:8888/%v/%v/1000/all/mobilez", method, userID), "application/json", nil)
			data := getJSON(resp, err, t)
			tkn, ok := data.Path("token").Data().(string)
			if !ok {
				t.Errorf("Token not found in response when getting token: %v", data.String())
			}
			token = &tkn
		})
		t.Run("Fetch data for token", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://:8888/%v/%v", method, *token))
			data := getJSON(resp, err, t)
			uid, ok := data.Path("UserId").Data().(string)
			if !ok {
				t.Errorf("UserId not found in response: %v", data.String())
			}
			if uid != userID {
				t.Errorf("Invalid userId %v expected %v", uid, userID)
			}
		})
	}
}

func getJSON(resp *http.Response, err error, t *testing.T) *gabs.Container {
	if err != nil {
		t.Errorf("Error when getting token: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	data, err := gabs.ParseJSON(body)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}
	return data
}
