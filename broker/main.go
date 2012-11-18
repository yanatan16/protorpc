package main

import (
	"flag"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/yanatan16/protorpc"
	"log"
)

func main() {
	f := flag.Int("fp", 12001, "Front End port number")
	b := flag.Int("bp", 12002, "Back End port number")
	flag.Parse()

	ctx, err := zmq.NewContext()
	if err != nil {
		log.Fatal("Error creating zmq context:", err)
	}

	front := fmt.Sprintf("*:%d", *f)
	back := fmt.Sprintf("*:%d", *b)

	broker, err := protorpc.NewBroker(ctx, front, back)
	if err != nil {
		log.Fatal("Error creating broker:", err)
	}

	log.Printf("Serving protorpc broker with frontend :%d and backend :%d", *f, *b)
	broker.Serve()
}
