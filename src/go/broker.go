package protorpc

import (
	"errors"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"log"
)

func Err(msg string, err error) error {
	return errors.New(msg + ": " + err.Error())
}

// A zmq REQ/REP broker with a Router frontend and a Dealer backend.
type Broker struct {
	frontend, backend zmq.Socket
}

func NewBroker(frontendAddr, backendAddr string) (*Broker, error) {
	frontend, err := context.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, Err("Couldn't create new ROUTER socket", err)
	}

	backend, err := context.NewSocket(zmq.DEALER)
	if err != nil {
		return nil, Err("Couldn't create new DEALER socket", err)
	}

	err = frontend.Bind(fmt.Sprintf("tcp://%s", frontendAddr))
	if err != nil {
		return nil, Err("Error binding front end ("+frontendAddr+")", err)
	}

	err = backend.Bind(fmt.Sprintf("tcp://%s", backendAddr))
	if err != nil {
		return nil, Err("Error binding front end ("+backendAddr+")", err)
	}

	return &Broker{
		frontend, backend,
	}, nil
}

func (b *Broker) Close() error {
	b.frontend.Close()
	b.backend.Close()
	return nil
}

// A blocking function that will infinitely forward multi-part messages between two zmq.Sockets
func Forward(a, b zmq.Socket) {
	for {
		parts, err := a.RecvMultipart(0)
		if err != nil {
			log.Println("Error receiving message on frontend broker", err)
		}

		err = b.SendMultipart(parts, 0)
		if err != nil {
			log.Println("Error sending message on backend broker", err)
		}

		// log.Println("Brokered message:", Stringify(parts))
	}
}

func (b *Broker) Serve() {
	go Forward(b.frontend, b.backend)
	Forward(b.backend, b.frontend)
}
