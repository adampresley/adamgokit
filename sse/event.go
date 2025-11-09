package sse

type Event struct {
	Event string `json:"event"`
	ID    string `json:"id"`
	Data  string `json:"data"`
}
