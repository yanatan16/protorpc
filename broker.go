// Simple request-reply broker
//
// Author:  Brendan Mc.
// Requires: http://github.com/alecthomas/gozmq

package protorpc

import (
	"errors"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

func Err(msg string, err error) error {
	return errors.New(msg + ": " + err.Error())
}

// A zmq REQ/REP broker with a Router frontend and a Dealer backend.
type Broker struct {
	frontend, backend zmq.Socket
}

func NewBroker(context zmq.Context, frontendAddr, backendAddr string) (*Broker, error) {
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

func (b *Broker) Serve() error {
	// Initialize poll set
	toPoll := zmq.PollItems{
		zmq.PollItem{Socket: b.frontend, zmq.Events: zmq.POLLIN},
		zmq.PollItem{Socket: b.backend, zmq.Events: zmq.POLLIN},
	}

	for {
		_, err := zmq.Poll(toPoll, -1)
		if err != nil {
			return err
		}

		switch {
		case toPoll[0].REvents&zmq.POLLIN != 0:
			message, err := toPoll[0].Socket.Recv(0)
			if err != nil {
				return err
			}
			b.backend.Send(message, 0)

		case toPoll[1].REvents&zmq.POLLIN != 0:
			message, err := toPoll[1].Socket.Recv(0)
			if err != nil {
				return err
			}
			b.frontend.Send(message, 0)
		}
	}
	return nil
}
