package pubsub

import (
	zmq "github.com/alecthomas/gozmq"
	"launchpad.net/tomb"
	"time"
)

// Subscription provides channel abstraction over zmq SUB sockets
type Subscription struct {
	addr    string
	filters []string
	Ch      chan Message // Channel to read messages from
	tomb.Tomb
}

func newSubscription(addr string, filters []string) *Subscription {
	sub := new(Subscription)
	sub.addr = addr
	sub.filters = filters
	sub.Ch = make(chan Message)
	go sub.loop()
	return sub
}

func (sub *Subscription) loop() {
	defer sub.Done()
	defer close(sub.Ch)

	ctx, err := GetGlobalContext()
	if err != nil {
		sub.Kill(err)
		return
	}

	// Establish a connection and subscription filter
	socket, err := ctx.NewSocket(zmq.SUB)
	if err != nil {
		sub.Kill(err)
		return
	}

	for _, filter := range sub.filters {
		err = socket.SetSockOptString(zmq.SUBSCRIBE, filter)
		if err != nil {
			sub.Kill(err)
			return
		}
	}

	err = socket.Connect(sub.addr)
	if err != nil {
		sub.Killf("Couldn't connect to %s: %s", sub.addr, err)
		return
	}

	// Read and stream the results in a channel
	pollItems := []zmq.PollItem{zmq.PollItem{socket, 0, zmq.POLLIN, 0}}

	for {
		n, err := zmq.Poll(pollItems, time.Duration(1)*time.Second)
		if err != nil {
			sub.Kill(err)
			return
		}

		select {
		case <-sub.Dying():
			return
		default:
		}

		if n > 0 {
			data, err := socket.Recv(zmq.DONTWAIT)
			if err != nil {
				sub.Kill(err)
				return
			}

			select {
			case sub.Ch <- NewMessage(data):
			case <-sub.Dying():
				return
			}
		}
	}
}

// Stop stops this Subscription
func (sub *Subscription) Stop() error {
	sub.Kill(nil)
	return sub.Wait()
}
