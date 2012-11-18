// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"log"
    "testing"
    "code.google.com/p/goprotobuf/proto"
    "github.com/yanatan16/protorpc"
    "net"
)

const server_addr string = "127.0.0.1:12345"

func init() {
    calc := new(MyCalcService)
    RegisterCalcService(calc)

    l, e := net.Listen("tcp",server_addr)
    if e != nil {
        log.Fatal("listen error:",e)
    }

    go protorpc.Serve(l)
}

func TestNoClient(t *testing.T) {
    calc := new(MyCalcService)

    doCalc(calc, t)
}

func TestClient(t *testing.T) {
	calc, err := NewCalcServiceClient("MyCalcService", "tcp", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	doCalc(calc, t)
}

func doCalc(calc CalcService, t *testing.T) {
    crq := new(CalcRequest)
    crs := new(CalcResponse)

    // add
    crq.A = proto.Int64(61)
    crq.B = proto.Int64(35)

    err := calc.Add(crq, crs)
    if err != nil {
        t.Error("add error:", err)
    } else if *crs.Result != 61 + 35 {
        t.Error("add result incorrect:", *crs.Result)
    }

    crq.Reset()
    crs.Reset()

    // subtract
    crq.A = proto.Int64(61)
    crq.B = proto.Int64(35)

    err = calc.Subtract(crq, crs)
    if err != nil {
        t.Error("subtract error:", err)
    } else if *crs.Result != 61 - 35 {
        t.Error("subtract result incorrect:", *crs.Result)
    }

    crq.Reset()
    crs.Reset()

    // multiply
    crq.A = proto.Int64(9)
    crq.B = proto.Int64(11)

    err = calc.Multiply(crq, crs)
    if err != nil {
        t.Error("multiply error:", err)
    } else if *crs.Result != 9 * 11 {
        t.Error("multiply result incorrect:", *crs.Result)
    }

    crq.Reset()
    crs.Reset()

    // divide
    crq.A = proto.Int64(20)
    crq.B = proto.Int64(8)

    err = calc.Divide(crq, crs)
    if err != nil {
        t.Error("divide error:", err)
    } else if *crs.Result != 20 / 8 {
        t.Error("divide result incorrect:", *crs.Result)
    } else if *crs.Remainder != 20 % 8 {
        t.Error("divide remainder incorrect:", *crs.Remainder)
    }

    crq.Reset()
    crs.Reset()

}
