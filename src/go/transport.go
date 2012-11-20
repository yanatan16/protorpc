// Package protorpc implements a fast, scalable, and schema'd RPC framework using ZMQ sockets and Protobufs serialization
package protorpc

import (
	"code.google.com/p/goprotobuf/proto"
	zmq "github.com/alecthomas/gozmq"
	"log"
	"io"
	"net/rpc"
)

var context zmq.Context

var SeqNum chan uint64

type bufferPair struct {
	header *proto.Buffer
	body   *proto.Buffer
}

func NewBufferPair() *bufferPair {
	return &bufferPair{proto.NewBuffer(nil), proto.NewBuffer(nil)}
}
func (b *bufferPair) Reset() {
	b.header.Reset()
	b.body.Reset()
}

type idMap map[uint64][][]byte

func init() {
	var err error
	context, err = zmq.NewContext()
	if err != nil {
		log.Fatal("Couldn't create zmq context!")
	}

	// Create sequence numbers
	SeqNum = make(chan uint64, 1)
	go func() {
		count := uint64(0)
		for {
			count += 1
			SeqNum <- count
		}
	}()
}

// Close the zmq context. Make sure to close each connection (or RPC) first.
func Close() {
	context.Close()
}

// Serve a port with a protobuf rpc server. Set brokered to true if you are using a broker
func Serve(addr string, brokered bool) (io.Closer, error) {
	// Socket to talk to clients
	conn, err := context.NewSocket(zmq.DEALER)
	if err != nil {
		return nil, err
	}

	if brokered {
		err = conn.Connect("tcp://" + addr)
	} else {
		err = conn.Bind("tcp://" + addr)
	}
	if err != nil {
		return nil, err
	}

	// log.Println("Created Server", conn, addr)
	server := NewServerCodec(conn)
	go rpc.ServeCodec(server)
	return server, nil
}

// Dial a connection to a protobuf rpc server
func Dial(addr string) (*rpc.Client, error) {
	// Socket to talk to server
	conn, err := context.NewSocket(zmq.DEALER)
	if err != nil {
		return nil, err
	}

	err = conn.Connect("tcp://" + addr)
	if err != nil {
		return nil, err
	}

	// log.Println("Created Client", conn, addr)

	client := rpc.NewClientWithCodec(NewClientCodec(conn))

	return client, nil
}

func Stringify(parts [][]byte) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = string(p)
	}
	return out
}
