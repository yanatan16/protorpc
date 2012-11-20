package protorpc

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
	"log"
	zmq "github.com/alecthomas/gozmq"
	"net/rpc"
)

type clientCodec struct {
	sock zmq.Socket
	bufs *bufferPair
}

func NewClientCodec(conn zmq.Socket) rpc.ClientCodec {
	bufs := NewBufferPair()

	return &clientCodec{conn, bufs}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, message interface{}) (err error) {
	c.bufs.Reset()

	h := new(Header)
	h.Id = &r.Seq
	h.ServiceMethod = &r.ServiceMethod

	err = c.bufs.header.Marshal(h)
	if err != nil {
		return
	}

	if msg, ok := message.(proto.Message); ok {
		err = c.bufs.body.Marshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}
	if err != nil {
		return
	}

	parts := [][]byte{c.bufs.header.Bytes(), c.bufs.body.Bytes()}
	err = c.sock.SendMultipart(parts, 0)
	// log.Println("Client sent multipart mesage:", Stringify(parts))

	log.Println("Client sent service request:", r)

	return
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) (err error) {
	c.bufs.Reset()

	// Read message
	parts, err := c.sock.RecvMultipart(0)
	if err != nil {
		return
	} else if len(parts) != 2 {
		err = errors.New("Unexpected number of parts to reply: " + fmt.Sprintf("%d", len(parts)))
		return
	}
	// log.Println("Client recieved multipart mesage:", Stringify(parts))

	hdr := parts[0]
	data := parts[1]

	c.bufs.body.SetBuf(data)

	h := new(Header)
	c.bufs.header.SetBuf(hdr)
	err = c.bufs.header.Unmarshal(h)
	if err != nil {
		return
	}

	r.Seq = h.GetId()
	r.ServiceMethod = h.GetServiceMethod()
	r.Error = h.GetError()

	log.Println("Client recieved service response:", r)

	return nil
}

func (c *clientCodec) ReadResponseBody(message interface{}) (err error) {

	if msg, ok := message.(proto.Message); ok {
		err = c.bufs.body.Unmarshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}
	return
}

func (c *clientCodec) Close() error {
	return c.sock.Close()
}
