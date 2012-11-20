// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	. "code.google.com/p/goprotobuf/protoc-gen-go/generator"
	"log"
)

func init() {
	rpc := new(RpcPlugin)
	RegisterPlugin(rpc)
}

type RpcPlugin struct {
	*Generator
}

func (*RpcPlugin) Name() string {
	return "protorpc"
}

func (g *RpcPlugin) Init(ng *Generator) {
	g.Generator = ng
}

func (g *RpcPlugin) Generate(file *FileDescriptor) {
	g.P()
	g.P("// protorpc code")

	for _, sd := range file.Service {
		serviceName := *sd.Name
		if serviceName == "" {
			log.Println("no service name")
			continue
		}

		// build the interface
		g.P("type ", serviceName, " interface {")
		g.In()
		for _, m := range sd.Method {
			name := *m.Name
			input_type := CamelCaseSlice(g.ObjectNamed(*m.InputType).TypeName())
			output_type := CamelCaseSlice(g.ObjectNamed(*m.OutputType).TypeName())

			g.P(name, "(*", input_type, ", *", output_type, ") error")

		}
		g.Out()
		g.P("}")

		// build server registration helper
		g.P("func Register", serviceName, "(s ", serviceName, ") error {")
		g.In()
		g.P("return rpc.Register(s)")
		g.Out()
		g.P("}")

		// build the concrete client
		g.P("type ", serviceName, "Client struct {")
		g.In()
		g.P("*rpc.Client")
		g.P("remoteName string")
		g.Out()
		g.P("}")

		// client constructor
		g.P("func New", serviceName, "Client(rname, raddr string) (csc *", serviceName, "Client, err error) {")
		g.In()
		g.P("client, err := protorpc.Dial(raddr)")
		g.P("if err != nil {")
		g.In()
		g.P("return")
		g.Out()
		g.P("}")
		g.P("csc = new(", serviceName, "Client)")
		g.P("csc.Client = client")
		g.P("csc.remoteName = rname")
		g.P("return")
		g.Out()
		g.P("}")

		// build methods on client
		for _, m := range sd.Method {
			name := *m.Name
			input_type := CamelCaseSlice(g.ObjectNamed(*m.InputType).TypeName())
			output_type := CamelCaseSlice(g.ObjectNamed(*m.OutputType).TypeName())

			g.P("func (self *", serviceName, "Client) ", name, "(request *", input_type, ", response *", output_type, ") error {")
			g.In()
			g.P("return self.Call(self.remoteName + ", Quote("."+name), ", request, response)")
			g.Out()
			g.P("}")

			// Build Asynchronous Request
			asyncName := name + "Async"
			g.P("func (self *", serviceName, "Client) ", asyncName, "(request *", input_type, ", response *", output_type, ") (chan error) {")
			g.In()
			g.P("ret := make(chan error, 0)")
			g.P("done := self.Go(self.remoteName + ", Quote("."+name), ", request, response, make(chan *rpc.Call, 1)).Done")
			g.P("go func() {")
			g.In()
			g.P("call := <- done")
			g.P("ret <- call.Error")
			g.Out()
			g.P("}()")
			g.P("return ret")
			g.Out()
			g.P("}")
		}

	}
}

func (g *RpcPlugin) GenerateImports(file *FileDescriptor) {
	g.P()
	g.P("// protorpc imports")
	g.P("import ", Quote("net/rpc"))
	g.P("import ", Quote("protorpc"))

	g.P("// Reference rpc and protorpc")
	g.P("var _ = rpc.DefaultRPCPath")
	g.P("var _ = protorpc.NewBufferPair")
}
