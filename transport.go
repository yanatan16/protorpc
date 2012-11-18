package protorpc

import (
	"bytes"
	"errors"
	zmq "github.com/alecthomas/gozmq"
	"net/rpc"
)

var context zmq.Context

func init() {
	var err error
	context, err = zmq.NewContext()
	if err != nil {
		panic("Couldn't create zmq context!")
	}
}

// Close the zmq context. Make sure to close each connection (or RPC) first.
func Close() {
	context.Close()
}

type ZSocket struct {
	zmq.Socket
	buf          *bytes.Buffer
	syncR, syncW chan bool
}

// Creates a zmq socket, conects it to addr and wraps it in a ReadWriteCloser
func NewZSocket(stype zmq.SocketType, addr string, bind, readFirst bool) (*ZSocket, error) {
	conn, err := context.NewSocket(stype)
	if err != nil {
		return nil, err
	}

	if bind {
		err = conn.Bind("tcp://" + addr)
	} else {
		err = conn.Connect("tcp://" + addr)
	}

	if err != nil {
		return nil, err
	}

	syncR := make(chan bool, 1)
	syncW := make(chan bool, 1)

	if readFirst {
		syncR <- true
	} else {
		syncW <- true
	}

	return &ZSocket{
		Socket: conn,
		buf:    nil,
		syncR:  syncR,
		syncW:  syncW,
	}, nil
}

func (z *ZSocket) Read(p []byte) (n int, err error) {
	if z.buf == nil || z.buf.Len() == 0 {

		<-z.syncR
		msg, err := z.Recv(0)
		z.syncW <- true

		if err != nil {
			return 0, err
		}
		z.buf = bytes.NewBuffer(msg)
	}

	n, err = z.buf.Read(p)
	if z.buf.Len() == 0 {
		z.buf = nil //Remove the reference
	}

	return n, err
}

func (z *ZSocket) Write(p []byte) (n int, err error) {

	if z.buf == nil {
		z.buf = bytes.NewBuffer(p)
		return len(p), err
	}
	return z.buf.Write(p)
}

func (z *ZSocket) Flush() error {
	if z.buf == nil {
		return errors.New("No bytes have been written")
	}

	bytes := z.buf.Bytes()
	z.buf = nil

	<-z.syncW
	err := z.Send(bytes, 0)
	z.syncR <- true

	return err
}

// Serve a port with a protobuf rpc server. Set directConn to false if you are using a broker
func Serve(addr string, directConn bool) error {
	// Socket to talk to clients
	conn, err := NewZSocket(zmq.REP, addr, directConn, true)
	if err != nil {
		return err
	}

	go rpc.ServeCodec(NewServerCodec(conn))
	return nil
}

// Dial a connection to a protobuf rpc server
func Dial(addr string) (*rpc.Client, error) {
	conn, err := NewZSocket(zmq.REQ, addr, false, false)
	if err != nil {
		return nil, err
	}

	client := rpc.NewClientWithCodec(NewClientCodec(conn))

	return client, nil
}
