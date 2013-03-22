package pubsub

import (
	"testing"
	"time"
)

// TODO: split the tests.

// Test a simple pub sub case.
func TestMonolithic(t *testing.T) {
	data := map[string]bool{
		"foo": true, "bar": true, "golang": true}
	z := Broker{
		PubAddr:    "tcp://127.0.0.1:5000",
		SubAddr:    "tcp://127.0.0.1:5001",
		BufferSize: 100}

	go func() {
		t.Fatal(z.Run())
	}()

	doPub := func(key string) {
		pub, err := z.NewPublisher()
		if err != nil {
			t.Fatal(err)
		}
		defer pub.Stop()
		for {
			for value, _ := range data {
				if err := pub.Publish(key, value); err != nil {
					t.Fatal(err)
				}
			}
			// continue to publish the same values as many times as
			// possible -- hence for {} -- as it appears normal for
			// zeromq to drop the initial values. we want the below
			// subscriber to catch some values at least.
		}
	}

	done := make(chan bool)

	go func() {
		sub := z.Subscribe("test")

		// retrieve at least three values from the eternal publishers.
		count := 3
		for count > 0 {
			msg := <-sub.Ch
			if msg.Key == "test" {
				if _, ok := data[msg.Value]; !ok {
					sub.Stop()
					t.Fatalf("got garbage data: %s", msg.Value)
				}
				count -= 1
			}
		}
		err := sub.Stop()
		if err != nil {
			t.Fatal(err)
		}
		done <- true
	}()

	go doPub("test")
	go doPub("ignore")

	select {
	case <-done:
		// all done
	case <-time.After(1 * time.Second):
		t.Fatal("not enough values from subscription")
	}
}
