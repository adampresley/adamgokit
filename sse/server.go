package sse

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/adampresley/adamgokit/httphelpers"
)

type Broker interface {
	Listen()
	Publish(event Event)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Publisher interface {
	Publish(event Event)
}

type SseBrokerConfig struct {
	AllowedOrigin string
	CancelContext context.Context
	EventChan     chan Event
}

type SseBroker struct {
	allowedOrigin  string
	cancelContext  context.Context
	clients        map[chan Event]struct{}
	closingClients chan chan Event
	eventChan      chan Event
	lock           *sync.Mutex
	newClients     chan chan Event
}

func NewBroker(config SseBrokerConfig) *SseBroker {
	result := &SseBroker{
		allowedOrigin:  config.AllowedOrigin,
		eventChan:      config.EventChan,
		newClients:     make(chan chan Event),
		closingClients: make(chan chan Event),
		clients:        make(map[chan Event]struct{}),
		cancelContext:  config.CancelContext,
		lock:           &sync.Mutex{},
	}

	if result.allowedOrigin == "" {
		result.allowedOrigin = "*"
	}

	return result
}

/*
Listen starts the broker's main loop for managing clients and broadcasting events.
This should be run in a separate goroutine after initializing the broker.
*/
func (b *SseBroker) Listen() {
	slog.Info("SSE broker started")
	defer slog.Info("SSE broker stopped")

	for {
		select {
		case <-b.cancelContext.Done():
			slog.Info("shutting down SSE broker")

			b.lock.Lock()

			for client := range b.clients {
				close(client)
			}

			b.lock.Unlock()

			return

		case client := <-b.newClients:
			b.lock.Lock()
			b.clients[client] = struct{}{}
			b.lock.Unlock()

			slog.Info("new SSE client connected", "totalClients", len(b.clients))

		case client := <-b.closingClients:
			b.lock.Lock()
			delete(b.clients, client)
			b.lock.Unlock()

			close(client)
			slog.Info("SSE client disconnected", "totalClients", len(b.clients))

		case event := <-b.eventChan:
			slog.Info("broker received event, broadcasting to clients", "event", event.Event, "id", event.ID, "clients", len(b.clients))

			for client := range b.clients {
				select {
				case client <- event:
				default:
					slog.Warn("client channel full. dropping event for a client.")
				}
			}
		}
	}
}

/*
Publish sends an event to all connected clients.
*/
func (b *SseBroker) Publish(event Event) {
	b.eventChan <- event
}

func (b *SseBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		ok      bool
		flusher http.Flusher
	)

	if flusher, ok = w.(http.Flusher); !ok {
		slog.Error("response writer does not support http.Flusher, SSE not possible")
		httphelpers.TextInternalServerError(w, "SSE streaming unsupported!")
		return
	}

	slog.Info("sse handler: new SSE connection established")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", b.allowedOrigin)

	/*
	 * Create a new channel for this client. Register them, and setup
	 * a a routine to unregister on disconnect. Then send an initial
	 * greeting.
	 */
	messageChan := make(chan Event, 10)
	b.newClients <- messageChan

	defer func() {
		select {
		case b.closingClients <- messageChan:

		case <-b.cancelContext.Done():
			slog.Info("sse handler: broker already shutting down, skip client close notification")
		}

		slog.Info("sse handler: client connection closed")
	}()

	connectionEvent := Event{
		Event: "connection",
		Data:  `{"status": "connected"}`,
	}

	if err = b.writeEvent(w, connectionEvent); err != nil {
		slog.Error("failed to send initial SSE connection message", "error", err)
		return
	}

	flusher.Flush()
	slog.Info("sse handler: sent initial SSE connection message")

	/*
	 * Loop and wait for events. Send them to the client.
	 */
	for {
		select {
		case <-b.cancelContext.Done():
			slog.Info("sse handler: SSE connection closed (in servehttp)")
			return

		case <-r.Context().Done():
			slog.Info("sse handler: client disconnected")
			return

		case event, ok := <-messageChan:
			if !ok {
				slog.Info("sse handler: client channel closed")
				return
			}

			slog.Info("sse handler: received event for client", "event", event.Event, "id", event.ID)

			if err = b.writeEvent(w, event); err != nil {
				slog.Error("sse handler: failed to write SSE event", "error", err)
				return
			}

			flusher.Flush()
			slog.Info("sse handler: sent event to client", "event", event.Event, "id", event.ID)
		}
	}
}

func (b *SseBroker) writeEvent(w http.ResponseWriter, event Event) error {
	var (
		err error
		// bytes []byte
	)

	// if bytes, err = json.Marshal(event.Data); err != nil {
	// 	return fmt.Errorf("failed to marshal SSE event: %w", err)
	// }

	if event.Event != "" {
		if _, err = fmt.Fprintf(w, "event: %s\n", event.Event); err != nil {
			return err
		}
	}

	if event.ID != "" {
		if _, err = fmt.Fprintf(w, "id: %s\n", event.ID); err != nil {
			return err
		}
	}

	// if len(bytes) > 0 && string(bytes) != "null" {
	if len(event.Data) > 0 && event.Data != "null" {
		// if _, err = fmt.Fprintf(w, "data: %s\n", string(bytes)); err != nil {
		if _, err = fmt.Fprintf(w, "data: %s\n", event.Data); err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "\n")
	return err
}
