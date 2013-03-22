package pubsub

import (
	zmq "github.com/alecthomas/gozmq"
	"sync"
)

// In zeromq, only the context objects are thread-safe, we share a
// single global context between goroutines.

var globalContext zmq.Context
var globalContextErr error
var once sync.Once

func initializeGlobalContext() {
	globalContext, globalContextErr = zmq.NewContext()
}

// GetGlobalContext returns a singleton zmq Context for the current Go
// process. 
func GetGlobalContext() (zmq.Context, error) {
	once.Do(initializeGlobalContext)
	return globalContext, globalContextErr
}
