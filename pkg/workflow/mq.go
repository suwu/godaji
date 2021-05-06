package workflow

import (
	"sync"
)

type Message struct {
	Name    string
	Payload map[string]interface{}
}

type MessageQueue struct {
	messages []*Message
	mu       sync.Mutex
}

func (mq *MessageQueue) Pop() *Message {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if len(mq.messages) == 0 {
		return nil
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]

	return message
}

func (mq *MessageQueue) Push(message *Message) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.messages = append(mq.messages, message)
}

func (mq *MessageQueue) IsEmpty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	l := len(mq.messages)
	return l == 0
}

func (mq *MessageQueue) Size() int {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	l := len(mq.messages)
	return l
}

func Dispatch(message *Message) {
	DefaultQueue.Push(message)
}

var defaultQueue MessageQueue
var DefaultQueue = &defaultQueue
