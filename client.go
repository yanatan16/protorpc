// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"io"
	"net"
	"net/rpc"
)

type clientCodec struct {
	c    io.ReadWriteCloser
	req  *bufferPair
	resp *bufferPair
}

func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	req := &bufferPair{proto.NewBuffer(nil), proto.NewBuffer(nil)}
	resp := &bufferPair{proto.NewBuffer(nil), proto.NewBuffer(nil)}

	return &clientCodec{conn, req, resp}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, message interface{}) (err error) {
	c.req.header.Reset()
	c.req.body.Reset()

	h := new(Header)
	h.Seq = &r.Seq
	h.ServiceMethod = &r.ServiceMethod

	err = c.req.header.Marshal(h)
	if err != nil {
		return
	}

	if msg, ok := message.(proto.Message); ok {
		err = c.req.body.Marshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}
	if err != nil {
		return
	}

	_, err = c.c.Write(encodeLen(len(c.req.header.Bytes())))
	if err != nil {
		return
	}

	_, err = c.c.Write(c.req.header.Bytes())
	if err != nil {
		return
	}

	_, err = c.c.Write(encodeLen(len(c.req.body.Bytes())))
	if err != nil {
		return
	}

	_, err = c.c.Write(c.req.body.Bytes())
	if err != nil {
		return
	}

	return
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) (err error) {
	c.resp.header.Reset()

	lbuf := make([]byte, lenSize)
	_, err = io.ReadFull(c.c, lbuf)
	if err != nil {
		return
	}

	pbuf := make([]byte, decodeLen(lbuf))
	_, err = io.ReadFull(c.c, pbuf)
	if err != nil {
		return
	}

	c.resp.header.SetBuf(pbuf)

	h := new(Header)
	err = c.resp.header.Unmarshal(h)
	if err != nil {
		return
	}

	r.Seq = *h.Seq
	r.ServiceMethod = *h.ServiceMethod
	r.Error = *h.Error

	return nil
}

func (c *clientCodec) ReadResponseBody(message interface{}) (err error) {
	c.resp.body.Reset()

	lbuf := make([]byte, lenSize)
	_, err = io.ReadFull(c.c, lbuf)
	if err != nil {
		return
	}

	pbuf := make([]byte, decodeLen(lbuf))
	_, err = io.ReadFull(c.c, pbuf)
	if err != nil {
		return
	}

	c.resp.body.SetBuf(pbuf)

	if msg, ok := message.(proto.Message); ok {
		err = c.resp.body.Unmarshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}
	if err != nil {
		return
	}

	return
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

func Dial(netw, raddr string) (*rpc.Client, error) {
	conn, err := net.Dial(netw, raddr)
	if err != nil {
		return nil, err
	}

	codec := NewClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)

	return client, nil
}
