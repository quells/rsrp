package rsrp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quells/rsrp"
)

func TestRedirectRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:3000/test/me", nil)
	request.Header.Add("X-Header", "Some Value")
	deadline := time.Now().Add(time.Second)
	ctx, cancel := context.WithDeadline(request.Context(), deadline)
	request = request.WithContext(ctx)
	request.AddCookie(&http.Cookie{
		Name:  "Cookie",
		Value: "Monster",
	})

	redirected, err := RedirectRequest(request, "http://localhost:3000/new/me")

	cancel()

	if err != nil {
		t.Fatalf("RedirectRequest() expected redirected error to be nil: %s", err.Error())
	}

	if redirected.URL.Path != "/new/me" {
		t.Fatalf("RedirectRequest() expected path to be %s, got %s", "/new/me", redirected.URL.Path)
	}

	headerValue := redirected.Header.Get("X-Header")
	if headerValue != "Some Value" {
		t.Fatalf("RedirectRequest() expected header value to be %s, got %s", "Some Value", headerValue)
	}

	redirectedDeadline, redirectedDeadlineOk := redirected.Context().Deadline()
	if !redirectedDeadlineOk {
		t.Fatalf("RedirectRequest() expected redirected deadline to be OK")
	}
	if !deadline.Equal(redirectedDeadline) {
		t.Fatalf("RedirectRequest() expected redirected deadline to be %s, got %s", deadline, redirectedDeadline)
	}

	redirectedCookie, err := redirected.Cookie("Cookie")
	if err != nil {
		t.Fatalf("RedirectRequest() expected cookie error to be nil: %s", err.Error())
	}
	if redirectedCookie.Value != "Monster" {
		t.Fatalf("RedirectRequest() expected cookie value to be %s, got %s", "Monster", redirectedCookie.Value)
	}
}

func TestRouteAll(t *testing.T) {
	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf("A: %s", r.URL.Path)
		w.Write([]byte(response))
	}))
	defer serverA.Close()

	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf("B: %s", r.URL.Path)
		w.Write([]byte(response))
	}))
	defer serverB.Close()

	configStr := fmt.Sprintf(`
	{
		"routes": [
			{
				"match": "^\/abc\/?.*$",
				"rewrite": {
					"src": "^\/abc(\/?.*)$",
					"dest": "$1"
				},
				"destination": "%s"
			},
			{
				"match": "^\/xyz\/?.*$",
				"rewrite": {
					"src": "^\/xyz(\/?.*)$",
					"dest": "/etc$1"
				},
				"destination": "%s"
			}
		]
	}
	`, serverA.URL, serverB.URL)
	config := &rsrp.Config{}
	_ = json.Unmarshal([]byte(configStr), config)

	routes, _ := rsrp.ConvertRules(config.Routes)

	server := httptest.NewServer(http.HandlerFunc(rsrp.RouteAll(*routes)))
	defer server.Close()

	client := server.Client()

	testCases := []struct {
		path, expected string
	}{
		{"/abc/test", "A: /test"},
		{"/abc", "A: /"},
		{"/xyz/test", "B: /etc/test"},
		{"/xyz", "B: /etc"},
	}

	var resp *http.Response
	var err error
	for _, tc := range testCases {
		resp, err = client.Get(server.URL + tc.path)
		if err != nil {
			t.Fatalf("RouteAll() unexpected error for %s: %s", tc.path, err.Error())
		}

		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		body := string(data)
		if body != tc.expected {
			t.Fatalf("RouteAll() expected %s to yield %s, got %s", tc.path, tc.expected, body)
		}
	}

}
