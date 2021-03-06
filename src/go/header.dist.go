package protorpc

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// protorpc imports
// import "net/rpc"
// import "github.com/yanatan16/protorpc"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Header struct {
	Id               *uint64 `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	ServiceMethod    *string `protobuf:"bytes,2,req,name=service_method" json:"service_method,omitempty"`
	Error            *string `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Header) Reset()         { *this = Header{} }
func (this *Header) String() string { return proto.CompactTextString(this) }
func (*Header) ProtoMessage()       {}

func (this *Header) GetId() uint64 {
	if this != nil && this.Id != nil {
		return *this.Id
	}
	return 0
}

func (this *Header) GetServiceMethod() string {
	if this != nil && this.ServiceMethod != nil {
		return *this.ServiceMethod
	}
	return ""
}

func (this *Header) GetError() string {
	if this != nil && this.Error != nil {
		return *this.Error
	}
	return ""
}

func init() {
}

// protorpc code
