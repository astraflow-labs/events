package events

import "net/http"

type NewServerPayload struct {
	ServerID string `json:"server_id"`
	AuthorID string `json:"author_id"`
	Network  string `json:"network"`
}

type NewServerResponsePayload struct {
	Status string `json:"status"`
}

type ProxyRequest struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    []byte      `json:"body"`
}

type ProxyResponse struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
}
