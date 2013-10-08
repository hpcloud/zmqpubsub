package zmqpubsub

import (
	zmq "github.com/alecthomas/gozmq"
)

type Broker struct {
	PubAddr         string // Publisher Endpoint Address
	SubAddr         string // Subscriber Endpoint Address
	BufferSize      int    // Memory buffer size
	SubscribeFilter string
}

// Run runs a broker for this pubsub configuration.
func (z Broker) Run() error {
	f, err := NewForwarder(z)
	if err == nil {
		err = f.Run()
	}
	return err
}

func (z Broker) MustRun() {
	panic(z.Run())
}

// Subscribe returns a subscription (channel) for given filters.
func (z Broker) Subscribe(filters ...string) *Subscription {
	if len(filters) == 0 {
		panic("Subscribe requires at least one filter")
	}
	return newSubscription(z.SubAddr, filters)
}

func (z Broker) NewPublisher() (*Publisher, error) {
	sock, err := newPubSocket(z.BufferSize)
	if err != nil {
		return nil, err
	}
	if err = sock.Connect(z.PubAddr); err != nil {
		sock.Close()
		return nil, err
	}
	// Publisher.Close is responsible for closing `sock`.
	return newPublisher(sock), nil
}

func (z Broker) NewPublisherMust() *Publisher {
	pub, err := z.NewPublisher()
	if err != nil {
		panic(err)
	}
	return pub
}

func newPubSocket(bufferSize int) (*zmq.Socket, error) {
	ctx, err := GetGlobalContext()
	if err != nil {
		return nil, err
	}

	sock, err := ctx.NewSocket(zmq.PUB)
	if err != nil {
		return nil, err
	}

	// prevent 0mq from infinitely buffering messages
	for _, hwm := range []zmq.IntSocketOption{zmq.SNDHWM, zmq.RCVHWM} {
		err = sock.SetSockOptInt(hwm, bufferSize)
		if err != nil {
			sock.Close()
			return nil, err
		}
	}

	return sock, nil
}
