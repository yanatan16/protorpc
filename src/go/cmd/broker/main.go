package main

import (
	"flag"
	"fmt"
	"log"
	"protorpc"
)

func main() {
	f := flag.Int("fp", 12001, "Front End port number")
	b := flag.Int("bp", 12002, "Back End port number")
	flag.Parse()

	front := fmt.Sprintf("*:%d", *f)
	back := fmt.Sprintf("*:%d", *b)

	broker, err := protorpc.NewBroker(front, back)
	if err != nil {
		log.Fatal("Error creating broker:", err)
	}

	log.Printf("Serving protorpc broker with frontend :%d and backend :%d", *f, *b)
	broker.Serve()
}
