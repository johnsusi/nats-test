package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, e1 := nats.Connect(nats.DefaultURL)
	if e1 != nil {
		log.Fatal(e1)
	}

	// Create JetStream Context
	js, e2 := nc.JetStream(nats.PublishAsyncMaxPending(256))

	if e2 != nil {
		log.Fatal(e2)
	}

	n := 10000
	m := 100

	// Simple Async Stream Publisher
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			subject := fmt.Sprintf("foo.%d.bar.%d", j, i)
			_, e3 := js.PublishAsync(subject, []byte("hello, world!"))
			if e1 != nil {
				log.Fatal(e3)
			}
		}
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}

}
