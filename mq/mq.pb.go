// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: mq/mq.proto

package mq

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MQRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prefix     string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Topic      string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	ReplyTopic string `protobuf:"bytes,3,opt,name=reply_topic,json=replyTopic,proto3" json:"reply_topic,omitempty"`
}

func (x *MQRule) Reset() {
	*x = MQRule{}
	mi := &file_mq_mq_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MQRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MQRule) ProtoMessage() {}

func (x *MQRule) ProtoReflect() protoreflect.Message {
	mi := &file_mq_mq_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MQRule.ProtoReflect.Descriptor instead.
func (*MQRule) Descriptor() ([]byte, []int) {
	return file_mq_mq_proto_rawDescGZIP(), []int{0}
}

func (x *MQRule) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *MQRule) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *MQRule) GetReplyTopic() string {
	if x != nil {
		return x.ReplyTopic
	}
	return ""
}

var file_mq_mq_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*MQRule)(nil),
		Field:         7229572,
		Name:          "mq.mq",
		Tag:           "bytes,7229572,opt,name=mq",
		Filename:      "mq/mq.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// See `MQRule`.
	//
	// optional mq.MQRule mq = 7229572;
	E_Mq = &file_mq_mq_proto_extTypes[0]
)

var File_mq_mq_proto protoreflect.FileDescriptor

var file_mq_mq_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x71, 0x2f, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x6d,
	0x71, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x06, 0x4d, 0x51, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70,
	0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x72,
	0x65, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x3a, 0x3d, 0x0a, 0x02,
	0x6d, 0x71, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x84, 0xa1, 0xb9, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x6d, 0x71,
	0x2e, 0x4d, 0x51, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x02, 0x6d, 0x71, 0x42, 0x2c, 0x5a, 0x2a, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x79, 0x74, 0x65, 0x63, 0x74,
	0x6c, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67,
	0x6f, 0x6d, 0x71, 0x2f, 0x6d, 0x71, 0x3b, 0x6d, 0x71, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_mq_mq_proto_rawDescOnce sync.Once
	file_mq_mq_proto_rawDescData = file_mq_mq_proto_rawDesc
)

func file_mq_mq_proto_rawDescGZIP() []byte {
	file_mq_mq_proto_rawDescOnce.Do(func() {
		file_mq_mq_proto_rawDescData = protoimpl.X.CompressGZIP(file_mq_mq_proto_rawDescData)
	})
	return file_mq_mq_proto_rawDescData
}

var file_mq_mq_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_mq_mq_proto_goTypes = []any{
	(*MQRule)(nil),                     // 0: mq.MQRule
	(*descriptorpb.MethodOptions)(nil), // 1: google.protobuf.MethodOptions
}
var file_mq_mq_proto_depIdxs = []int32{
	1, // 0: mq.mq:extendee -> google.protobuf.MethodOptions
	0, // 1: mq.mq:type_name -> mq.MQRule
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	1, // [1:2] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mq_mq_proto_init() }
func file_mq_mq_proto_init() {
	if File_mq_mq_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mq_mq_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_mq_mq_proto_goTypes,
		DependencyIndexes: file_mq_mq_proto_depIdxs,
		MessageInfos:      file_mq_mq_proto_msgTypes,
		ExtensionInfos:    file_mq_mq_proto_extTypes,
	}.Build()
	File_mq_mq_proto = out.File
	file_mq_mq_proto_rawDesc = nil
	file_mq_mq_proto_goTypes = nil
	file_mq_mq_proto_depIdxs = nil
}