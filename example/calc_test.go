// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/yanatan16/protorpc"
	"log"
	"math/rand"
	"testing"
	"time"
)

const server_addr string = "127.0.0.1:12345"
const server2_addr string = "127.0.0.1:12346"
var calcSvr *MyCalcService

func init() {
	calcSvr := new(MyCalcService)
	RegisterCalcService(calcSvr)
	err := protorpc.Serve(server_addr, false)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	slowCalcSvr := new(SlowCalcService)
	RegisterCalcService(slowCalcSvr)
	err = protorpc.Serve(server2_addr, false)
	if err != nil {
		log.Fatal("cant setup calc service 2:", err)
	}

	<-time.After(50 * time.Millisecond)
}

func TestNoClient(t *testing.T) {
	doCalc(calcSvr, t)
}

func TestClient(t *testing.T) {
	calc, err := NewCalcServiceClient("MyCalcService", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	defer calc.Close()

	doCalc(calc, t)
}

func TestBrokeredClient(t *testing.T) {
	broker_front := "127.0.0.1:12001"
	broker_back := "127.0.0.1:12002"

	broker, err := protorpc.NewBroker(broker_front, broker_back)
	if err != nil {
		log.Fatal("Error creating broker:", err)
	}
	go func() {
		broker.Serve()
	}()
	<-time.After(50 * time.Millisecond)

	calcSvr := new(MyCalcService)
	RegisterCalcService(calcSvr)

	err = protorpc.Serve(broker_back, true)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	<-time.After(50 * time.Millisecond)

	calc, err := NewCalcServiceClient("MyCalcService", broker_front)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	defer calc.Close()

	log.Println("TestBrokeredClient now set up, ready to run...")
	doCalc(calc, t)
}

func TestMultiReqClient(t *testing.T) {
	var nReq = 10

	calc, err := NewCalcServiceClient("SlowCalcService", server2_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	defer calc.Close()

	// only add
	crq := make([]*CalcRequest, nReq)
	crs := make([]*CalcResponse, nReq)
	res := make([]int64, nReq)
	errs := make([]chan error, nReq)

	for i := 0; i < nReq; i++ {
		a, b := rand.Int63(), rand.Int63()
		crq[i] = new(CalcRequest)
		crq[i].A = &a
		crq[i].B = &b
		crs[i] = new(CalcResponse)
		res[i] = a+b

		errs[i] = calc.AddAsync(crq[i], crs[i])
	}

	for i, cherr := range errs {
		err := <- cherr
		if err != nil {
			t.Fatal("add error:", err)
		} else if *crs[i].Result != res[i] {
			t.Error("add result incorrect:", *crs[i].Result, "vs", res[i])
		}

	}
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
	} else if *crs.Result != 61+35 {
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
	} else if *crs.Result != 61-35 {
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
	} else if *crs.Result != 9*11 {
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
	} else if *crs.Result != 20/8 {
		t.Error("divide result incorrect:", *crs.Result)
	} else if *crs.Remainder != 20%8 {
		t.Error("divide remainder incorrect:", *crs.Remainder)
	}

	crq.Reset()
	crs.Reset()

}
