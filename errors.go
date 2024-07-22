package events

import "errors"

var (
	// ErrResponseTimeouted is returned when no response is received in time
	ErrResponseTimeouted = errors.New("no response from other side in time")
)
