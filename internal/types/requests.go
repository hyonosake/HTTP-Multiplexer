package types

type MultiplyRequest struct {
	URLs []string `json:"urls"`
}

type JsonData struct {
	URL  string `json:"url"`
	Data string `json:"response"`
}

type MultiplyResponse struct {
	Error string     `json:"error,omitempty"`
	Data  []JsonData `json:"data,omitempty"`
}
