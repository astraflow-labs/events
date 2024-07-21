package events

import (
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type Event struct {
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
	Data      []byte `json:"data"`
}

func NewEvent(eventType string, eventID string, data []byte) Event {
	return Event{
		EventType: eventType,
		EventID:   eventID,
		Data:      data,
	}
}

