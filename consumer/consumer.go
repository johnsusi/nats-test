package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	subject := "foo.*.bar.*"
	if len(os.Args) > 1 {
		subject = "foo.50.bar.*"
	}
	// Connect to NATS
	nc, e1 := nats.Connect(nats.DefaultURL)
	if e1 != nil {
		log.Fatal(e1)
	}
	// Create JetStream Context
	js, e2 := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if e2 != nil {
		log.Fatal(e1)
	}

	sub, e3 := js.SubscribeSync(subject)
	if e3 != nil {
		log.Fatal(e1)
	}

	start := time.Now()
	count := 0
	for {
		m, err := sub.NextMsg(time.Second * 1)
		if err != nil {
			log.Print(err)
			break
		}
		m.Ack()
		count += 1
		if (count % 1000) == 0 {
			elapsed := time.Since(start).Seconds()
			rate := float64(count) / elapsed
			log.Printf("Processed %d messages in %d seconds. %d msgs / sec.", count, int(elapsed), int(rate))
		}
	}
	sub.Unsubscribe()
	sub.Drain()

}
