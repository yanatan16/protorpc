PROTOC=protoc

PROTO_FILES = header.proto

all: $(PROTO_FILES)

$(PROTO_FILES): %.proto:
	$(PROTOC) --go-rpc_out=./go $*.proto