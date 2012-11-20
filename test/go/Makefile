.PHONY: gotest

test: calc.proto godep gotest

%.proto:
	protoc --go-rpc_out=. $*.proto
	
godep:
	go test -i

gotest:
	go test
