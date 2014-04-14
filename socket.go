// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package zmqpubsub

import (
  zmq "github.com/alecthomas/gozmq"
)

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

func newPubSocketBound(bufferSize int, addr string) (*zmq.Socket, error) {
  sock, err := newPubSocket(bufferSize)
  if err != nil {
    return nil, err
  }
  return sock, sock.Bind(addr)
}

func newSubSocketBound(subfilter string, addr string) (*zmq.Socket, error) {
  ctx, err := GetGlobalContext()
  if err != nil {
    return nil, err
  }

  sock, err := ctx.NewSocket(zmq.SUB)
  if err != nil {
    return nil, err
  }

  if err = sock.Bind(addr); err != nil {
    return nil, err
  }

  return sock, sock.SetSockOptString(zmq.SUBSCRIBE, subfilter)
}
