package queue

type Queue interface {
	SendMessage(interface{}) error
	ReceiveMessage() (interface{}, error)
}

type queue struct {
	messages chan interface{}
}

func NewQueue() Queue {
	return &queue{make(chan interface{})}
}

func (q *queue) SendMessage(message interface{}) error {
	go func(message interface{}) {
		q.messages <- message
	}(message)

	return nil
}

func (q *queue) ReceiveMessage() (interface{}, error) {
	return <-q.messages, nil
}
