// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package zmqpubsub

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
	return z.NewPublisher2("")
}

func (z Broker) NewPublisher2(addr string) (*Publisher, error) {
	sock, err := newPubSocket(z.BufferSize)
	if err != nil {
		return nil, err
	}
	if addr == "" {
		addr = z.PubAddr
	}
	if err = sock.Connect(addr); err != nil {
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
