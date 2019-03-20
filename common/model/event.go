package model

type Event struct {
	Type   string
	Action string
	//le sub Java n'accepte pas d'array
	//Message   []*Dog
	Message   *Dog
	Timestamp int64
}
