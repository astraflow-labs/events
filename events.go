package events

import (
	"strings"

	"github.com/google/uuid"
)

type Event struct {
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
	Data      []byte `json:"data"`
}

func NewEvent(eventType string, data []byte) Event {
	uid := uuid.New().String()
	return Event{
		EventType: strings.ToUpper(eventType),
		EventID:   uid,
		Data:      data,
	}
}

func (e Event) Respond(data []byte) Event {
	return Event{
		EventType: e.EventType + "RSP",
		EventID:   e.EventID,
		Data:      data,
	}
}
