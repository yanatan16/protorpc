package protorpc

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	zmq "github.com/alecthomas/gozmq"
	"net/rpc"
)

type serverCodec struct {
	sock      zmq.Socket
	req       *proto.Buffer
	packetIds idMap
	clientIds map[uint64]uint64
}

func NewServerCodec(conn zmq.Socket) rpc.ServerCodec {
	req := proto.NewBuffer(nil)
	packetIds := make(idMap)
	clientIds := make(map[uint64]uint64)

	return &serverCodec{conn, req, packetIds, clientIds}
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) (err error) {

	// Read message
	parts, err := c.sock.RecvMultipart(0)
	if err != nil {
		return
	}
	log.Println("Server recieved multipart mesage:", Stringify(parts))

	nId := len(parts) - 2
	ids := parts[:nId]
	hdr := parts[nId]
	data := parts[nId+1]

	// Read header and body into proto buf
	c.req.SetBuf(data)

	// Unmarshal header
	h := new(Header)
	err = proto.Unmarshal(hdr, h)
	if err != nil {
		return
	}

	// Increment sequence id and save ids for sending request
	seq := <-SeqNum
	c.packetIds[seq] = ids
	c.clientIds[seq] = h.GetId()

	r.Seq = seq
	r.ServiceMethod = h.GetServiceMethod()

	log.Println("Server recieved service request:", r)

	return
}

func (c *serverCodec) ReadRequestBody(message interface{}) (err error) {
	// Request body should have already been stored in c.req.body
	// Unmarshal body
	if msg, ok := message.(proto.Message); ok {
		err = c.req.Unmarshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}

	return
}

func (c *serverCodec) WriteResponse(r *rpc.Response, message interface{}) (err error) {
	resp := NewBufferPair()

	// Extract information and delete from map
	seq := r.Seq
	packetIds := c.packetIds[seq]
	clientId, ok := c.clientIds[seq]

	delete(c.packetIds, seq)
	delete(c.clientIds, seq)

	if !ok {
		err = errors.New("Sequence number was not found in id maps.")
		return
	}

	// Create the header
	h := new(Header)
	h.Id = &clientId
	h.ServiceMethod = &r.ServiceMethod
	h.Error = &r.Error

	err = resp.header.Marshal(h)
	if err != nil {
		return
	}

	// Create the body
	if msg, ok := message.(proto.Message); ok {
		err = resp.body.Marshal(msg)
		if err != nil {
			return
		}
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
		return
	}

	// Create and send message
	parts := make([][]byte, 0, len(packetIds)+2)
	parts = append(parts, packetIds...)
	parts = append(parts, resp.header.Bytes(), resp.body.Bytes())

	err = c.sock.SendMultipart(parts, 0)
	// log.Println("Server sent multipart mesage:", Stringify(parts))
	log.Println("Server recieved service response:", r)

	return
}

func (c *serverCodec) Close() error {
	return c.sock.Close()
}
