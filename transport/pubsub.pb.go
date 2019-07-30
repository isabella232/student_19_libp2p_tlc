// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pubsub.proto

package transport

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MsgType int32

const (
	MsgType_Raw MsgType = 0
)

var MsgType_name = map[int32]string{
	0: "Raw",
}

var MsgType_value = map[string]int32{
	"Raw": 0,
}

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return proto.EnumName(MsgType_name, int32(x))
}

func (x *MsgType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MsgType_value, data, "MsgType")
	if err != nil {
		return err
	}
	*x = MsgType(value)
	return nil
}

func (MsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_91df006b05e20cf7, []int{0}
}

type Message struct {
	MsgType              *MsgType   `protobuf:"varint,1,req,name=Msg_type,json=MsgType,enum=transport.MsgType" json:"Msg_type,omitempty"`
	Source               *int32     `protobuf:"varint,2,req,name=source" json:"source,omitempty"`
	Step                 *int32     `protobuf:"varint,3,req,name=step" json:"step,omitempty"`
	History              []*Message `protobuf:"bytes,4,rep,name=history" json:"history,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_91df006b05e20cf7, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetMsgType() MsgType {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return MsgType_Raw
}

func (m *Message) GetSource() int32 {
	if m != nil && m.Source != nil {
		return *m.Source
	}
	return 0
}

func (m *Message) GetStep() int32 {
	if m != nil && m.Step != nil {
		return *m.Step
	}
	return 0
}

func (m *Message) GetHistory() []*Message {
	if m != nil {
		return m.History
	}
	return nil
}

func init() {
	proto.RegisterEnum("transport.MsgType", MsgType_name, MsgType_value)
	proto.RegisterType((*Message)(nil), "transport.Message")
}

func init() { proto.RegisterFile("pubsub.proto", fileDescriptor_91df006b05e20cf7) }

var fileDescriptor_91df006b05e20cf7 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x28, 0x4d, 0x2a,
	0x2e, 0x4d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2c, 0x29, 0x4a, 0xcc, 0x2b, 0x2e,
	0xc8, 0x2f, 0x2a, 0x51, 0x9a, 0xcc, 0xc8, 0xc5, 0xee, 0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x2a,
	0xa4, 0xc7, 0xc5, 0xe1, 0x5b, 0x9c, 0x1e, 0x5f, 0x52, 0x59, 0x90, 0x2a, 0xc1, 0xa8, 0xc0, 0xa4,
	0xc1, 0x67, 0x24, 0xac, 0x07, 0x57, 0xa9, 0x07, 0x93, 0x0a, 0x62, 0xf7, 0x2d, 0x4e, 0x0f, 0xa9,
	0x2c, 0x48, 0x15, 0x12, 0xe3, 0x62, 0x2b, 0xce, 0x2f, 0x2d, 0x4a, 0x4e, 0x95, 0x60, 0x52, 0x60,
	0xd2, 0x60, 0x0d, 0x82, 0xf2, 0x84, 0x84, 0xb8, 0x58, 0x8a, 0x4b, 0x52, 0x0b, 0x24, 0x98, 0xc1,
	0xa2, 0x60, 0xb6, 0x90, 0x0e, 0x17, 0x7b, 0x46, 0x66, 0x71, 0x49, 0x7e, 0x51, 0xa5, 0x04, 0x8b,
	0x02, 0xb3, 0x06, 0xb7, 0x91, 0x10, 0xb2, 0xd1, 0x10, 0x07, 0x04, 0xc1, 0x94, 0x68, 0x09, 0x23,
	0x5c, 0x22, 0xc4, 0xce, 0xc5, 0x1c, 0x94, 0x58, 0x2e, 0xc0, 0x00, 0x08, 0x00, 0x00, 0xff, 0xff,
	0x24, 0x52, 0x3a, 0x2d, 0xc4, 0x00, 0x00, 0x00,
}
