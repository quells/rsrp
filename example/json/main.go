package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/nocontent", handleNoContent)
	http.HandleFunc("/echo", handleEcho)

	httpServer := &http.Server{
		Addr: ":5002",
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// PingResponse is the response for a ping request
type PingResponse struct {
	Pong bool `json:"pong"`
}

// NewPingResponse is a constructor for PingResponse
func NewPingResponse() PingResponse {
	return PingResponse{true}
}

// Serialize converts a PingResponse to JSON bytes
func (r PingResponse) Serialize() []byte {
	data, _ := json.Marshal(r)
	return data
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(NewPingResponse().Serialize())
}

func handleNoContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Header", "Nothing to see here; move along.")
	w.WriteHeader(http.StatusNoContent)
}

// EchoResponse is the response for an echo request
type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    *interface{}        `json:"body"`
}

// Serialize converts a EchoResponse to JSON bytes
func (r EchoResponse) Serialize() []byte {
	data, _ := json.Marshal(r)
	return data
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	var body *interface{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		copy := interface{}(err.Error())
		body = &copy
	} else {
		jsonified := &map[string]interface{}{}
		err = json.Unmarshal(data, jsonified)
		if err != nil {
			done := make(chan bool)
			go func() {
				copy := interface{}(string(data))
				body = &copy
				done <- true

				if recover() != nil {
					msg := interface{}(fmt.Sprintf("%v", data))
					body = &msg
					done <- true
				}
			}()
			_ = <-done
		} else {
			copy := interface{}(jsonified)
			body = &copy
		}
	}

	query := make(map[string][]string)
	for k, v := range r.URL.Query() {
		query[k] = v
	}

	response := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   query,
		Headers: r.Header,
		Body:    body,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(response.Serialize())
}
