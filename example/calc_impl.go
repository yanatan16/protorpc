// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
)

type MyCalcService int

func (*MyCalcService) Add(req *CalcRequest, resp *CalcResponse) error {
	resp.Result = proto.Int64((*req.A) + (*req.B))
	return nil
}

func (*MyCalcService) Subtract(req *CalcRequest, resp *CalcResponse) error {
	resp.Result = proto.Int64((*req.A) - (*req.B))
	return nil
}

func (*MyCalcService) Multiply(req *CalcRequest, resp *CalcResponse) error {
	resp.Result = proto.Int64((*req.A) * (*req.B))
	return nil
}

func (*MyCalcService) Divide(req *CalcRequest, resp *CalcResponse) (err error) {
	resp.Result = proto.Int64(0)
	defer func() {
		if x := recover(); x != nil {
			if ex, ok := x.(error); ok {
				err = ex
			} else {
				err = errors.New(fmt.Sprint(x))
			}
		}
	}()
	resp.Result = proto.Int64((*req.A) / (*req.B))
	resp.Remainder = proto.Int64((*req.A) % (*req.B))
	return
}
