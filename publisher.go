package pubsub

import (
	zmq "github.com/alecthomas/gozmq"
)

// Publisher creates a thread/goroutine-unsafe publisher. Only a
// single goroutine must use the created publisher.
type Publisher struct {
	sock *zmq.Socket
}

func newPublisher(sock *zmq.Socket) *Publisher {
	pub := new(Publisher)
	pub.sock = sock
	return pub
}

func (p *Publisher) Publish(key string, value string) error {
	return p.sock.Send([]byte(key+" "+value), 0)
}

func (p *Publisher) MustPublish(key string, value string) {
	if err := p.Publish(key, value); err != nil {
		panic(err)
	}
}

func (p *Publisher) Stop() {
	p.sock.Close()
}
