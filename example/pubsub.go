package main

import (
	"fmt"
	"github.com/ActiveState/zmqpubsub"
	"math/rand"
	"time"
)

var Broker zmqpubsub.Broker

func init() {
	Broker.PubAddr = "tcp://127.0.0.1:4000"
	Broker.SubAddr = "tcp://127.0.0.1:4001"
	Broker.BufferSize = 100
}

func RunPublisher(name string) {
	pub := Broker.NewPublisherMust()
	defer pub.Stop()

	count := 1
	for {
		pub.MustPublish(name, fmt.Sprintf("%d", count))
		count += 1
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func EchoSubscribe() {
	sub := Broker.Subscribe("")
	defer sub.Stop()

	fmt.Println("Monitoring subscription..")
	for msg := range sub.Ch {
		fmt.Printf("%12s => %3s\n", msg.Key, msg.Value)
	}

	// if there was an error, `sub.Ch` will be closed and `sub.Wait`
	// will report the original error.
	err := sub.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Println("End of subscription")
}

func main() {
	// Subscribe
	fmt.Println("Starting subscriber")
	go EchoSubscribe()

	// Setup sample publishers
	fmt.Println("Setting up publishers")
	go RunPublisher("bonobo")
	go RunPublisher("chimpanzee")
	go RunPublisher("lemur")
	go RunPublisher("macaque")

	// Run the broker
	fmt.Printf("Running the broker: %+v\n", Broker)
	Broker.MustRun()
}
