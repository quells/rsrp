package relay

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// EchoServer implements a WebSocket echo server for testing
type EchoServer struct{}

var upgrader = websocket.Upgrader{}

func (s EchoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade to websocket", http.StatusBadRequest)
		return
	}

	echo := make(chan message)

	go func(c *websocket.Conn, e chan message) {
		defer func() {
			conn.Close()
		}()

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			return nil
		})

		for {
			messageType, body, err := conn.ReadMessage()
			if err != nil {
				break
			}

			e <- message{messageType, body}
		}
	}(conn, echo)

	go func(c *websocket.Conn, e chan message) {
		ticker := time.NewTicker(2 * time.Second)
		defer func() {
			ticker.Stop()
			conn.Close()
			close(echo)
		}()

		for {
			select {
			case msg, ok := <-e:
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
				}

				err := conn.WriteMessage(msg.messageType, msg.body)
				if err != nil {
					return
				}

			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}(conn, echo)
}
