# protorpc

## About

protorpc is an RPC implementation that uses 0-MQ as its transport layer and Google Protobufs as its transport encoding.

UA RPC has (or plans to have) clients implemented in the following languages:

- Go: Integrated with the Go standard library net/rpc
- Java (todo)
- Python (todo)
- Node.js (todo)

Features:

- Asynchronous call ability
- Can have one or more brokers in the RPC call

## History

protorpc is based off of [protorpc](http://github.com/eclark/protorpc) written by Eric Clark.That material is governed by the first license in the LICENSE file. Further modifications are governed by the second license in the LICENSE file.

## Requirements

Go - http://golang.org/doc/install.html
protobuf - http://code.google.com/p/protobuf/
goprotobuf - http://code.google.com/p/goprotobuf/
0mq - http://zeromq.org
gozmq - http://github.com/alecthomas/gozmq

## Installing

`go get github.com/yanatan16/protorpc`
`go install github.com/yanatan16/protorpc/protoc-gen-go-rpc`

Make sure $GOPATH/bin is on your $PATH. `export PATH=$PATH:$GOPATH/bin`

Now install libzmq 3.x from [zeromq](http://zeromq.org).

```
wget http://download.zeromq.org/zeromq-3.2.1-rc2.tar.gz
tar xzvf zeromq-3.2.1-rc2.tar.gz
cd zeromq-3.2.1-rc2
./configure
make
sudo make install
```

Now install gozmq for zmq 3.x

`go get -tags zmq_3_x github.com/alecthomas/gozmq`

_Note_: You will want to test the installation by going to the installation directory `$GOPATH/src/github.com/alecthomas/gozmq` and `go test`.  You may have to install a patch `git pull git@github.com:srid/gozmq patch-1`. If so, then reinstall `go install -tags zmq_3_x`. You may need to make sure libzmq.so.3 is on your $LD_LIBRARY_PATH.


## Usage

To compile `.proto` files, use the protoc command:

`protoc --go-rpc_out=. files.proto`

Then just compile your coe as usual. `go build`

The best way to get started is to review the files in the example
directory.
