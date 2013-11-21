# zmqpubsub

zmqpubsub is a simple Go pubsub implementation on top of ZeroMQ.

It abstracts the underlying ZeroMQ machinery to provide a Go-friendly
API for the publish-subscribe messaging pattern.

# Usage

## Broker

First setup a broker:

```Go
var Broker zmqpubsub.Broker

func init() {
	Broker.PubAddr = "tcp://127.0.0.1:4000"
	Broker.SubAddr = "tcp://127.0.0.1:4001"
	Broker.BufferSize = 100
}

func main() {
    ...
    Broker.MustRun()
}
```

The broker specifies the addresses to which publishers/subscribers 
will connect to/from.

## Subscriber

Subscription messages are sent in a Go channel. Thanks to [Tomb](
https://launchpad.net/tomb), subscriptions can be
stopped at any time by calling `Stop`.

```Go
sub := Broker.Subscribe("")
defer sub.Stop()

for msg := range sub.Ch {
    fmt.Printf("%s => %s\n", msg.Key, msg.Value)
}
```

## Publisher

The publisher part is equivalently simple. It should be noted however 
that as publishers are not thread-safe they must be managed from the 
same goroutine that created them.

```Go
pub := Broker.NewPublisherMust()
defer pub.Stop()

pub.MustPublish("key", "hello world")
pub.MustPublish("key", "hello universe")
```

# Example

To run the provided example,

```
go run -tags zmq_3_x example/pubsub.go
```
