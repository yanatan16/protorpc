// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"io"
	"net/rpc"
)

type serverCodec struct {
	c    ReadWriteFlushCloser
	req  *bufferPair
	resp *bufferPair
}

type ReadWriteFlushCloser interface {
	io.Reader
	io.Writer
	io.Closer
	Flush() error
}

func NewServerCodec(conn ReadWriteFlushCloser) rpc.ServerCodec {
	req := &bufferPair{proto.NewBuffer(nil), proto.NewBuffer(nil)}
	resp := &bufferPair{proto.NewBuffer(nil), proto.NewBuffer(nil)}

	return &serverCodec{conn, req, resp}
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) (err error) {
	c.req.header.Reset()

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

	c.req.header.SetBuf(pbuf)

	h := new(Header)
	err = c.req.header.Unmarshal(h)
	if err != nil {
		return
	}

	r.Seq = h.GetSeq()
	r.ServiceMethod = h.GetServiceMethod()

	return
}

func (c *serverCodec) ReadRequestBody(message interface{}) (err error) {
	c.req.body.Reset()

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

	c.req.body.SetBuf(pbuf)

	if msg, ok := message.(proto.Message); ok {
		err = c.req.body.Unmarshal(msg)
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
	}
	if err != nil {
		return
	}

	return
}

func (c *serverCodec) WriteResponse(r *rpc.Response, message interface{}) (err error) {
	c.resp.header.Reset()
	c.resp.body.Reset()

	h := new(Header)
	h.Seq = &r.Seq
	h.ServiceMethod = &r.ServiceMethod
	h.Error = &r.Error

	err = c.resp.header.Marshal(h)
	if err != nil {
		return
	}

	_, err = c.c.Write(encodeLen(len(c.resp.header.Bytes())))
	if err != nil {
		return
	}

	_, err = c.c.Write(c.resp.header.Bytes())
	if err != nil {
		return
	}

	if msg, ok := message.(proto.Message); ok {
		err = c.resp.body.Marshal(msg)
		if err != nil {
			return
		}

		_, err = c.c.Write(encodeLen(len(c.resp.body.Bytes())))
		if err != nil {
			return
		}

		_, err = c.c.Write(c.resp.body.Bytes())
		if err != nil {
			return
		}
	} else {
		err = errors.New("Message body does not implement goprotobuf.Message")
		return
	}

	err = c.c.Flush()

	return
}

func (c *serverCodec) Close() error {
	return c.c.Close()
}
