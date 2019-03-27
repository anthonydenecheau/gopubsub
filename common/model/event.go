package model

type Event struct {
	Type      string `json:"type"`
	Action    string `json:"action"`
	Message   []*Dog `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
