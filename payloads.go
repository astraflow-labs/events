package events

type NewServerPayload struct {
	ServerID string `json:"server_id"`
	AuthorID string `json:"author_id"`
	Network  string `json:"network"`
}

type NewServerResponsePayload struct {
	Status string `json:"status"`
}
