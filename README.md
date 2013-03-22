# pubsub

pubsub is a simple Go pubsub implementation on top of ZeroMQ.

It abstracts the underlying ZeroMQ machinery to provide a Go-friendly
API for the publish-subscribe messaging pattern.

# Usage

## Broker

First setup a broker:

```Go
var Broker pubsub.Broker

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

The broker specifies the addresses to which publishers and subscribers
(whether from the same process or a different node) will connect to.

## Subscriber

Subscription messages are sent in a Go channel. Subscriptions can be
stopped at any time by calling `Stop`.

```Go
sub := Broker.Subscribe("")
defer sub.Stop()

for msg := range sub.Ch {
    fmt.Printf("%s => %s\n", msg.Key, msg.Value)
}
```

## Publisher

The publisher part is equivalently simple. Do remember, however, that
as publishers are not thread-safe they must be managed from the same
goroutine that created them.

```Go
pub := Broker.NewPublisherMust()
defer pub.Stop()

pub.MustPublish("key", "hello world")
pub.MustPublish("key", "hello universe")
```

# Example

See `pubsub-example/pubsub.go` for a complete example.