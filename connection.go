package events

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws                 *websocket.Conn
	send               chan Event
	waitingForResponse map[string]chan Event
	wfrMutex           sync.RWMutex
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:                 ws,
		send:               make(chan Event),
		waitingForResponse: make(map[string]chan Event),
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

		c.wfrMutex.RLock()
		ch, ok := c.waitingForResponse[event.EventID]
		c.wfrMutex.RUnlock()

		if ok {
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

func (c *Connection) WaitForResponse(event Event) (Event, error) {

	ch := make(chan Event)
	c.wfrMutex.Lock()
	c.waitingForResponse[event.EventID] = ch
	c.wfrMutex.Unlock()

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(5 * time.Second):

		// Cleanup in the waiting for response map
		c.wfrMutex.Lock()
		delete(c.waitingForResponse, event.EventID)
		c.wfrMutex.Unlock()

		// Close the channel
		close(ch)

		// Return timeout error
		return Event{}, ErrResponseTimeouted
	}
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
	c.wfrMutex.Lock()
	defer c.wfrMutex.Unlock()
	for _, ch := range c.waitingForResponse {
		close(ch)
	}

	return nil
}
