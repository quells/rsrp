package relay_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/quells/rsrp/relay"

	"github.com/gorilla/websocket"
)

func TestHandler(t *testing.T) {
	hiddenListener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("could not create TCP listener: %v", err)
	}
	defer hiddenListener.Close()

	hiddenServer := &http.Server{
		Handler:      relay.EchoServer{},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	defer hiddenServer.Close()

	go func() {
		if err := hiddenServer.Serve(hiddenListener); err != http.ErrServerClosed {
			t.Fatalf("failed to start echo server: %v", err)
		}
	}()

	hiddenURL := "ws://" + hiddenListener.Addr().String()

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("could not create TCP listener: %v", err)
	}
	defer l.Close()

	s := &http.Server{
		Handler:      relay.DefaultHandler(hiddenURL),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	defer s.Close()

	go func() {
		if err := s.Serve(l); err != http.ErrServerClosed {
			t.Fatalf("failed to start relay server: %v", err)
		}
	}()

	url := "ws://" + l.Addr().String()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not connect to echo server: %v", err)
	}
	defer conn.Close()

	for i := 0; i < 100; i++ {
		conn.SetWriteDeadline(time.Now().Add(time.Second))
		payload := fmt.Sprintf("%d", i)
		err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			t.Fatalf("could not send message: %v", err)
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("could not read message: %v", err)
		}
		response := string(msg)
		if response != payload {
			t.Fatalf("expected %s, got %s", payload, response)
		}
	}
}
