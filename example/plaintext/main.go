package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handlePing)

	httpServer := &http.Server{
		Addr: ":5001",
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}
