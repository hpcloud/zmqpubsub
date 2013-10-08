package zmqpubsub

import (
	zmq "github.com/alecthomas/gozmq"
)

// Forwarder is a zeromq forwarder device acting as a broker between
// multiple publishers and multiple subscribes.
type Forwarder struct {
	ctx      *zmq.Context
	frontend *zmq.Socket
	backend  *zmq.Socket
	options  Broker
}

func NewForwarder(options Broker) (*Forwarder, error) {
	var err error
	b := new(Forwarder)
	b.options = options

	if b.ctx, err = GetGlobalContext(); err != nil {
		return nil, err
	}

	// Publishers speak to the frontend socket
	if b.frontend, err = b.ctx.NewSocket(zmq.SUB); err != nil {
		b.ctx.Close()
		return nil, err
	}
	if err = b.frontend.Bind(options.PubAddr); err != nil {
		b.ctx.Close()
		return nil, err
	}
	if err = b.frontend.SetSockOptString(
		zmq.SUBSCRIBE, options.SubscribeFilter); err != nil {
		b.ctx.Close()
		return nil, err
	}

	// Subscribers speak to the backend socket
	if b.backend, err = newPubSocket(options.BufferSize); err != nil {
		b.ctx.Close()
		return nil, err
	}
	if err = b.backend.Bind(options.SubAddr); err != nil {
		b.ctx.Close()
		return nil, err
	}

	return b, nil
}

func (b *Forwarder) Run() error {
	return zmq.Device(zmq.FORWARDER, b.frontend, b.backend)
}
