// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/yanatan16/protorpc"
	"log"
	"testing"
	"time"
)

func TestNoClient(t *testing.T) {
	calc := new(MyCalcService)

	doCalc(calc, t)
}

func TestClient(t *testing.T) {
	server_addr := "127.0.0.1:12345"

	calcSvr := new(MyCalcService)
	RegisterCalcService(calcSvr)

	err := protorpc.Serve(server_addr, false)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	<-time.After(50 * time.Millisecond)

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
