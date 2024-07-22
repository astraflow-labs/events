package events

import (
	"context"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws                 *websocket.Conn
	send               chan Event
	ctx                context.Context
	Close              context.CancelFunc
	mu                 sync.Mutex
	waitingForResponse map[string]chan Event
}

func NewConnection(ws *websocket.Conn) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	return &Connection{
		ws:    ws,
		send:  make(chan Event),
		ctx:   ctx,
		Close: cancel,
	}
}

type EventHandler func(Event, *Connection)

func (c *Connection) ReadLoop(eventHandler EventHandler) {
	defer func() {
		c.closeOnce()
	}()

	for {
		select {
		case <-c.ctx.Done():
			// Handle shutdown
			log.Info("ReadLoop shutting down gracefully")
			return
		default:
			var event Event
			err := c.ws.ReadJSON(&event)
			if err != nil {
				log.Error("Error reading from websocket", "err", err)
				c.closeOnce() // Signal other loops to shut down
				return
			}

			if ch, ok := c.waitingForResponse[event.EventID]; ok {
				ch <- event
				delete(c.waitingForResponse, event.EventID)
				close(ch)
				continue
			}

			eventHandler(event, c)
		}
	}
}

func (c *Connection) WriteLoop() {
	defer func() {
		c.ws.Close()
	}()

	//lint:ignore S1000 Måske ska vi ha flere listeners i fremtiden
	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				// Kanalen er lukket, luk også websocket
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.ws.WriteJSON(event)
			if err != nil {
				log.Error("Error writing to websocket", "err", err)
				return
			}
		}
	}
}

// Implement closeOnce to ensure ws.Close() is called only once
func (c *Connection) closeOnce() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ws != nil {
		c.ws.Close()
		c.ws = nil // Prevent further use
	}
}

func (c *Connection) Send() chan<- Event {
	return c.send
}

func (c *Connection) WaitForResponse(event Event) chan Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.waitingForResponse == nil {
		c.waitingForResponse = make(map[string]chan Event)
	}
	ch := make(chan Event, 1)
	c.waitingForResponse[event.EventID] = ch
	return ch
}
