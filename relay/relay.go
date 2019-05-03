// Package relay implements a rerouted WebSocket connection handler
package relay

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Handler handles a single rerouted WebSocket connection
type Handler struct {
	targetURL   string
	upgrader    websocket.Upgrader
	pumpOptions PumpOptions
}

// DefaultHandler creates a new relay.Handler targeting a given URL with default options
func DefaultHandler(targetURL string) Handler {
	return NewHandler(targetURL, websocket.Upgrader{}, DefaultPumpOptions())
}

// NewHandler creates a new relay.Handler targeting a given URL
func NewHandler(targetURL string, upgrader websocket.Upgrader, pumpOptions PumpOptions) Handler {
	return Handler{targetURL, upgrader, pumpOptions}
}

// ServeHTTP conforms relay.Handler to http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	external, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade to websocket", http.StatusBadRequest)
		return
	}

	internal, response, err := websocket.DefaultDialer.Dial(h.targetURL, nil)
	if err != nil {
		if response != nil {
			for header, value := range response.Header {
				for _, v := range value {
					w.Header().Add(header, v)
				}
			}
			responseBody, err := ioutil.ReadAll(response.Body)
			if err == nil {
				w.Write(responseBody)
			}
		} else {
			http.Error(w, "could not upgrade internal connection", http.StatusInternalServerError)
			return
		}
	}

	pump := NewPump(external, internal, h.pumpOptions)
	go pump.read(true)
	go pump.read(false)
	go pump.write()
}

type message struct {
	messageType int
	body        []byte
}

// Pump directs messages to and from the target WebSocket
type Pump struct {
	external, internal *websocket.Conn
	inbound, outbound  chan message
	options            PumpOptions
}

// PumpOptions includes some
type PumpOptions struct {
	WriteWait, PongWait, PingPeriod time.Duration
	MaxMessageSize                  int
}

// DefaultPumpOptions creates default relay.Pump options
func DefaultPumpOptions() PumpOptions {
	return PumpOptions{
		WriteWait:      10 * time.Second,
		PongWait:       60 * time.Second,
		PingPeriod:     54 * time.Second,
		MaxMessageSize: 1024,
	}
}

// NewPump creates a new relay.Pump between two WebSocket connections
func NewPump(external, internal *websocket.Conn, options PumpOptions) Pump {
	inbound := make(chan message)
	outbound := make(chan message)
	return Pump{external, internal, inbound, outbound, options}
}

func (p *Pump) read(forExternal bool) {
	var conn *websocket.Conn
	var channel chan message
	if forExternal {
		conn = p.external
		channel = p.inbound
	} else {
		conn = p.internal
		channel = p.outbound
	}

	defer func() {
		conn.Close()
	}()

	for {
		messageType, body, err := conn.ReadMessage()
		if err != nil {
			break
		}

		channel <- message{messageType, body}
	}
}

func (p *Pump) write() {
	externalTicker := time.NewTicker(p.options.PingPeriod)
	internalTicker := time.NewTicker(p.options.PingPeriod)
	defer func() {
		externalTicker.Stop()
		internalTicker.Stop()
		p.external.Close()
		p.internal.Close()
	}()

	for {
		select {
		case msg, ok := <-p.outbound:
			p.external.SetWriteDeadline(time.Now().Add(p.options.WriteWait))
			if !ok {
				p.external.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := p.external.WriteMessage(msg.messageType, msg.body)
			if err != nil {
				return
			}

		case msg, ok := <-p.inbound:
			p.internal.SetWriteDeadline(time.Now().Add(p.options.WriteWait))
			if !ok {
				p.internal.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := p.internal.WriteMessage(msg.messageType, msg.body)
			if err != nil {
				return
			}

		case <-externalTicker.C:
			p.external.SetWriteDeadline(time.Now().Add(p.options.WriteWait))
			if err := p.external.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-internalTicker.C:
			p.internal.SetWriteDeadline(time.Now().Add(p.options.WriteWait))
			if err := p.internal.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
