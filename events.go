package events

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

func (e Event) Respond(data []byte) Event {
	return Event{
		EventType: e.EventType + "RSP",
		EventID:   e.EventID,
		Data:      data,
	}
}
