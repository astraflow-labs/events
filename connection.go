package events

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws                 *websocket.Conn
	send               chan Event
	waitingForResponse map[string]chan Event
	mu                 sync.Mutex
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan Event),
	}
}

type EventHandler func(Event, *Connection)

func (c *Connection) ReadLoop(eventHandler EventHandler) error {
	for {
		var event Event
		err := c.ws.ReadJSON(&event)
		if err != nil {
			return err
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

func (c *Connection) WriteLoop() error {
	defer func() {
		c.ws.Close()
	}()

	//lint:ignore S1000 Måske ska vi ha flere listeners i fremtiden
	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				// Kanalen er lukket, luk også websocket
				c.Close()
				return errors.New("send channel closed")
			}
			err := c.ws.WriteJSON(event)
			if err != nil {
				return err
			}
		}
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

func (c *Connection) Close() error {
	err := c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	err = c.ws.Close()
	if err != nil {
		return err
	}

	close(c.send)
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.waitingForResponse {
		close(ch)
	}

	return nil
}
