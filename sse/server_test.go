package sse

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBroker(t *testing.T) {
	t.Run("With Empty Config", func(t *testing.T) {
		broker := NewBroker(SseBrokerConfig{})
		assert.NotNil(t, broker)
		assert.Equal(t, "*", broker.allowedOrigin, "allowedOrigin should default to '*'")
		assert.NotNil(t, broker.clients)
		assert.NotNil(t, broker.newClients)
		assert.NotNil(t, broker.closingClients)
		assert.NotNil(t, broker.lock)
	})

	t.Run("With Specific Config", func(t *testing.T) {
		eventChan := make(chan Event)
		ctx := context.Background()

		broker := NewBroker(SseBrokerConfig{
			AllowedOrigin: "https://example.com",
			EventChan:     eventChan,
			CancelContext: ctx,
		})

		assert.Equal(t, "https://example.com", broker.allowedOrigin)
		assert.Equal(t, eventChan, broker.eventChan)
		assert.Equal(t, ctx, broker.cancelContext)
	})
}

func TestSseBroker_writeEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "Event with Data, ID, and Name",
			event: Event{
				Event: "message",
				ID:    "123",
				Data:  `{"text": "hello"}`,
			},
			expected: "event: message\nid: 123\ndata: {\"text\": \"hello\"}\n\n",
		},
		{
			name: "Event with only Data",
			event: Event{
				Data: "just data",
			},
			expected: "data: just data\n\n",
		},
		{
			name: "Event with no Data",
			event: Event{
				Event: "ping",
			},
			expected: "event: ping\n\n",
		},
		{
			name:     "Empty Event",
			event:    Event{},
			expected: "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			broker := NewBroker(SseBrokerConfig{})

			err := broker.writeEvent(w, tc.event)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, w.Body.String())
		})
	}
}

func TestSseBroker_ServeHTTP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := NewBroker(SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     make(chan Event),
	})

	go broker.Listen()

	req := httptest.NewRequestWithContext(ctx, "GET", "/sse", nil)
	w := httptest.NewRecorder()

	var wg sync.WaitGroup
	wg.Add(1)

	slog.Info("calling serveHTTP")

	go func() {
		defer wg.Done()
		broker.ServeHTTP(w, req)
	}()

	// Give the handler a moment to start and write the initial event
	slog.Info("sleeping for initial event")
	time.Sleep(500 * time.Millisecond)

	// Stop the request to unblock ServeHTTP
	slog.Info("stopping request")
	cancel()
	slog.Info("waiting for handler to finish")
	wg.Wait() // Wait for ServeHTTP to finish

	slog.Info("handler finished")
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Check for the initial connection event
	expectedInitialEvent := "event: connection\ndata: {\"status\": \"connected\"}\n\n"
	assert.True(t, strings.HasPrefix(string(body), expectedInitialEvent), "Response body should start with the connection event")
}

func TestSseBroker_ClientManagement(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	broker := NewBroker(SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     make(chan Event),
	})
	go broker.Listen()

	// Test client connection
	clientChan1 := make(chan Event, 1)
	broker.newClients <- clientChan1
	time.Sleep(50 * time.Millisecond) // Give broker time to process

	broker.lock.Lock()
	assert.Len(t, broker.clients, 1, "Should have 1 client after connection")
	broker.lock.Unlock()

	// Test another client connection
	clientChan2 := make(chan Event, 1)
	broker.newClients <- clientChan2
	time.Sleep(50 * time.Millisecond)

	broker.lock.Lock()
	assert.Len(t, broker.clients, 2, "Should have 2 clients after second connection")
	broker.lock.Unlock()

	// Test client disconnection
	broker.closingClients <- clientChan1
	time.Sleep(50 * time.Millisecond)

	broker.lock.Lock()
	assert.Len(t, broker.clients, 1, "Should have 1 client after disconnection")
	_, exists := broker.clients[clientChan1]
	assert.False(t, exists, "Disconnected client should be removed from map")
	broker.lock.Unlock()

	// Ensure the disconnected channel is closed
	select {
	case _, ok := <-clientChan1:
		assert.False(t, ok, "Disconnected client channel should be closed")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Channel was not closed after disconnect")
	}
}

func TestSseBroker_Broadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	eventChan := make(chan Event, 1)
	broker := NewBroker(SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     eventChan,
	})
	go broker.Listen()

	// Connect two clients
	client1 := make(chan Event, 1)
	client2 := make(chan Event, 1)
	broker.newClients <- client1
	broker.newClients <- client2
	time.Sleep(50 * time.Millisecond) // Allow broker to process new clients

	// Send an event to broadcast
	testEvent := Event{ID: "test-1", Data: "broadcast data"}
	eventChan <- testEvent

	// Verify both clients received the event
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case received := <-client1:
			assert.Equal(t, testEvent, received)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client 1 did not receive event in time")
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case received := <-client2:
			assert.Equal(t, testEvent, received)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client 2 did not receive event in time")
		}
	}()

	wg.Wait()
}

// This test is designed to be run with the -race flag to detect race conditions.
func TestSseBroker_RaceTest(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	eventChan := make(chan Event, 10)
	broker := NewBroker(SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     eventChan,
	})
	go broker.Listen()

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Simulate concurrent connections and disconnections
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(i int) {
			defer wg.Done()
			clientChan := make(chan Event, 1)
			broker.newClients <- clientChan
			time.Sleep(time.Duration(10+i%10) * time.Millisecond) // Stagger operations
			broker.closingClients <- clientChan
		}(i)
	}

	// Simulate concurrent broadcasting
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(i int) {
			defer wg.Done()
			eventChan <- Event{Data: "some data"}
			time.Sleep(time.Duration(10+i%10) * time.Millisecond)
		}(i)
	}

	wg.Wait()
}

func TestSseBroker_ServeHTTP_EventFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := NewBroker(SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     make(chan Event),
	})
	go broker.Listen()

	req := httptest.NewRequest("GET", "/sse", nil)
	w := httptest.NewRecorder()

	var wg sync.WaitGroup

	wg.Go(func() {
		broker.ServeHTTP(w, req)
	})

	// Wait for the client to be registered
	time.Sleep(100 * time.Millisecond)
	broker.lock.Lock()
	require.Len(t, broker.clients, 1)
	broker.lock.Unlock()

	// Get the client channel from the broker
	broker.lock.Lock()
	var clientChan chan Event
	for ch := range broker.clients {
		clientChan = ch
	}
	broker.lock.Unlock()
	require.NotNil(t, clientChan)

	// Send an event to the client
	clientChan <- Event{Event: "test-event", Data: "hello world"}

	// Stop the handler and wait for it to finish
	cancel()
	wg.Wait()

	// Now it's safe to read the entire body
	body := w.Body.String()

	expectedEvents := []string{
		"event: connection\ndata: {\"status\": \"connected\"}\n\n",
		"event: test-event\ndata: hello world\n\n",
	}

	assert.Equal(t, strings.Join(expectedEvents, ""), body)
}
