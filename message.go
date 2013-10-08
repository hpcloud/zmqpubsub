package zmqpubsub

import (
	"strings"
)

// Message represents a zeromq message with two parts, Key and Value
// separated by a single space assuming the convention that Key is
// used to match against subscribe filters.
type Message struct {
	Key   string
	Value string
}

func NewMessage(data []byte) Message {
	parts := strings.SplitN(string(data), " ", 2)
	return Message{parts[0], parts[1]}
}
