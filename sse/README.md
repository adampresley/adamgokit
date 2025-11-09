# Server-Sent Events (SSE)

This package provides a broker for handling Server-Sent Events (SSE). It allows you to broadcast messages to multiple connected clients.

## How It Works

The `SseBroker` manages client connections. You create a broker, start it, and then it listens for events on a channel. When an event is received, it is broadcast to all connected clients.

The broker needs to be started in a separate goroutine. It also exposes an `http.Handler` to be used with your web server.

## Usage

Here is a basic example of how to set up and use the SSE broker.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adampresley/adamgokit/sse"
)

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel for SSE events
	eventChan := make(chan sse.Event)

	// Configure the SSE broker
	brokerConfig := sse.SseBrokerConfig{
		CancelContext: ctx,
		EventChan:     eventChan,
	}

	// Create a new broker
	broker := sse.NewBroker(brokerConfig)

	// Start the broker in a new goroutine
	go broker.Listen()

	// Set up a simple web server
	http.Handle("/sse", broker)

	// In a separate goroutine, send some events
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				fmt.Println("Sending event...")

				broker.Publish(sse.Event{
					Event: "message",
					Data:  fmt.Sprintf("%s", time.Now().Format(time.RFC1123)),
				})
			}
		}
	}()

	fmt.Println("Server started on :8080. Connect to http://localhost:8080/sse")
	http.ListenAndServe(":8080", nil)
}
```

To test this, you can run the server and then use `curl`:

```bash
curl http://localhost:8080/sse
```

You will see events stream to your console.

## Notes

- If you are using GZip compression, be sure that your SSE handler path is excluded.
- Provide an event to segregate messages so you can have independent handlers in your JavaScript

## JavaScript

Here is a small example of consuming these events.

```html
<p>
	The time is
	<span id="theTime"></span>
</p>

<script>

const evt = new EventSource("/sse");

evt.addEventListener("message", (e) => {
	document.querySelector("#theTime").innerText = e.data;
});

</script>
```

### HTMX

If you want to use HTMX (easy mode), first get the extension at https://cdn.jsdelivr.net/npm/htmx-ext-sse@2.2.4. Then change the publish to look something like this:

```go
broker.Publish(sse.Event{
	Event: "message",
	Data:  fmt.Sprintf(`<span id="theTime">%s</span>`, time.Now().Format(time.RFC1123)),
})
```

Then change your HTML to:

```html
<p>
	The time is
	<span id="theTime" hx-ext="sse" sse-connect="/sse" sse-swap="message"></span>
</p>
```
