// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

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
	b.frontend, err = newSubSocketBound(options.SubscribeFilter, options.PubAddr)
	if err != nil {
		b.ctx.Close()
		return nil, err
	}

	// Subscribers speak to the backend socket
	if b.backend, err = newPubSocketBound(options.BufferSize, options.SubAddr); err != nil {
		b.ctx.Close()
		return nil, err
	}

	return b, nil
}

func (b *Forwarder) Run() error {
	return zmq.Device(zmq.FORWARDER, b.frontend, b.backend)
}
