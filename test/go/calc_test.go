// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"code.google.com/p/goprotobuf/proto"
	"protorpc"
	"log"
	"math/rand"
	"testing"
	"time"
	"io"
	"fmt"
)

var (
	server_port int = 12345
	server_addr, broker_front string
	calcSvr *MyCalcService
	slowCalcSvr *SlowCalcService

	servers []io.Closer
	_ = rand.Int31
)

func init() {
	servers = make([]io.Closer,0)
	calcSvr := new(MyCalcService)
	RegisterCalcService(calcSvr)
	slowCalcSvr := new(SlowCalcService)
	RegisterCalcService(slowCalcSvr)
}

func StopAll(clients ...io.Closer) {
	log.Println("Finalizing all servers")
	for _, cli := range clients {
		err := cli.Close()
		if err != nil {
			log.Println("ERROR: ", err)
		}
	}
	for _, svr := range servers {
		err := svr.Close()
		if err != nil {
			log.Println("ERROR: ", err)
		}
	}
	servers = []io.Closer{}
}

func GetAddr() string {
	addr := fmt.Sprintf("127.0.0.1:%d", server_port)
	server_port += 1
	return addr
}

func StartBasicRPC() {
	server_addr = GetAddr()
	svr, err := protorpc.Serve(server_addr, false)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	servers = append(servers, svr)
}

func StartBrokeredRPC() {
	broker_front = GetAddr()
	broker_back := GetAddr()
	broker, err := protorpc.NewBroker(broker_front, broker_back)
	if err != nil {
		log.Fatal("Error creating broker:", err)
	}
	go func() {
		broker.Serve()
	}()
	servers = append(servers, broker)
	<-time.After(50 * time.Millisecond)

	svr, err := protorpc.Serve(broker_back, true)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}
	servers = append(servers, svr)
}

func TestNoClient(t *testing.T) {
	log.Println("TestNoClient now set up, ready to run...")
	doCalc(calcSvr, t)
}

// func TestStartStop(t *testing.T) {
// 	StartBasicRPC()
// 	StopAll()
// 	StartBrokeredRPC()
// 	StopAll()
// }

func TestClient(t *testing.T) {
	StartBasicRPC()

	calc, err := NewCalcServiceClient("MyCalcService", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	log.Println("TestClient now set up, ready to run...")
	doCalc(calc, t)

	servers = append(servers, calc)
	// StopAll(calc)
}

func TestBrokeredClient(t *testing.T) {
	StartBrokeredRPC()

	calc, err := NewCalcServiceClient("MyCalcService", broker_front)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	log.Println("TestBrokeredClient now set up, ready to run...")
	doCalc(calc, t)

	servers = append(servers, calc)
	// StopAll(calc)
}

func TestMultiReqClient(t *testing.T) {
	StartBasicRPC()

	var nReq = 10

	calc, err := NewCalcServiceClient("MyCalcService", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	// only add
	crq := make([]*CalcRequest, nReq)
	crs := make([]*CalcResponse, nReq)
	res := make([]int64, nReq)
	errs := make([]chan error, nReq)

	for i := 0; i < nReq; i++ {
		a, b := rand.Int63n(10000000), rand.Int63n(1000000000)
		crq[i] = new(CalcRequest)
		crq[i].A = &a
		crq[i].B = &b
		crs[i] = new(CalcResponse)

		if rand.Int31n(2)%2 > 0 {
			errs[i] = calc.AddAsync(crq[i], crs[i])
			res[i] = a + b
		} else {
			errs[i] = calc.MultiplyAsync(crq[i], crs[i])
			res[i] = a * b
		}
	}

	for i, cherr := range errs {
		err := <-cherr
		if err != nil {
			t.Fatal("add error:", err)
		} else if *crs[i].Result != res[i] {
			t.Error("add result incorrect:", *crs[i].Result, "vs", res[i])
		}
	}
	
	servers = append(servers, calc)
// 	StopAll(calc)
}

func TestMultiClient(t *testing.T) {
	StartBasicRPC()

	calc1, err := NewCalcServiceClient("MyCalcService", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	calc2, err := NewCalcServiceClient("MyCalcService", server_addr)
	if err != nil {
		log.Fatal("cant setup calc service:", err)
	}

	log.Println("TestMultiClient now set up, ready to run...")
	ret := make(chan bool, 0)
	go func() {
		doCalc(calc1, t)
		ret <- true
	}()
	go func() {
		doCalc(calc2, t)
		ret <- true
	}()

	<-ret
	<-ret
	
	log.Println("TestMultiClient done...")
	
	servers = append(servers, calc1, calc2)
	// StopAll(calc1, calc2)
}

func TestClose(t *testing.T) {
	StopAll()
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
