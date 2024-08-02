package events

type NewServerPayload struct {
	ServerID string `json:"server_id"`
	AuthorID string `json:"author_id"`
	Network  string `json:"network"`
}

type NewServerResponsePayload struct {
	Status string `json:"status"`
}

type ProxyRequest struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

type ProxyResponse struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}
