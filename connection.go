package events

import (
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws   *websocket.Conn
	send chan Event
	// eventLock sync.Mutex
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan Event),
	}
}

func (c *Connection) ReadLoop(handleEvent func(Event, *Connection)) {
	defer func() {
		c.ws.Close()
	}()

	for {
		var event Event
		err := c.ws.ReadJSON(&event)
		if err != nil {
			log.Error("Error reading from websocket", "err", err)
			break
		}

		handleEvent(event, c)
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

func (c *Connection) Close() error {
	close(c.send)
	return c.ws.Close()
}

func (c *Connection) Send() chan<- Event {
	return c.send
}
