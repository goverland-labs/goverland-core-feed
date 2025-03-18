// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: feedpb/subscription.proto

package feedpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubscribeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DaoId         string                 `protobuf:"bytes,2,opt,name=dao_id,json=daoId,proto3" json:"dao_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	mi := &file_feedpb_subscription_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feedpb_subscription_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_feedpb_subscription_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeRequest) GetDaoId() string {
	if x != nil {
		return x.DaoId
	}
	return ""
}

type UnsubscribeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DaoId         string                 `protobuf:"bytes,2,opt,name=dao_id,json=daoId,proto3" json:"dao_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnsubscribeRequest) Reset() {
	*x = UnsubscribeRequest{}
	mi := &file_feedpb_subscription_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnsubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsubscribeRequest) ProtoMessage() {}

func (x *UnsubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feedpb_subscription_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsubscribeRequest.ProtoReflect.Descriptor instead.
func (*UnsubscribeRequest) Descriptor() ([]byte, []int) {
	return file_feedpb_subscription_proto_rawDescGZIP(), []int{1}
}

func (x *UnsubscribeRequest) GetDaoId() string {
	if x != nil {
		return x.DaoId
	}
	return ""
}

var File_feedpb_subscription_proto protoreflect.FileDescriptor

var file_feedpb_subscription_proto_rawDesc = []byte{
	0x0a, 0x19, 0x66, 0x65, 0x65, 0x64, 0x70, 0x62, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x66, 0x65, 0x65,
	0x64, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x29, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x61, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x61, 0x6f, 0x49, 0x64, 0x22, 0x2b, 0x0a, 0x12, 0x55,
	0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x61, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x64, 0x61, 0x6f, 0x49, 0x64, 0x32, 0x90, 0x01, 0x0a, 0x0c, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x09, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x18, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x70, 0x62, 0x2e,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x41, 0x0a, 0x0b, 0x55, 0x6e, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x1a, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x70, 0x62,
	0x2e, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0a, 0x5a, 0x08, 0x2e,
	0x3b, 0x66, 0x65, 0x65, 0x64, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_feedpb_subscription_proto_rawDescOnce sync.Once
	file_feedpb_subscription_proto_rawDescData = file_feedpb_subscription_proto_rawDesc
)

func file_feedpb_subscription_proto_rawDescGZIP() []byte {
	file_feedpb_subscription_proto_rawDescOnce.Do(func() {
		file_feedpb_subscription_proto_rawDescData = protoimpl.X.CompressGZIP(file_feedpb_subscription_proto_rawDescData)
	})
	return file_feedpb_subscription_proto_rawDescData
}

var file_feedpb_subscription_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_feedpb_subscription_proto_goTypes = []any{
	(*SubscribeRequest)(nil),   // 0: feedpb.SubscribeRequest
	(*UnsubscribeRequest)(nil), // 1: feedpb.UnsubscribeRequest
	(*emptypb.Empty)(nil),      // 2: google.protobuf.Empty
}
var file_feedpb_subscription_proto_depIdxs = []int32{
	0, // 0: feedpb.Subscription.Subscribe:input_type -> feedpb.SubscribeRequest
	1, // 1: feedpb.Subscription.Unsubscribe:input_type -> feedpb.UnsubscribeRequest
	2, // 2: feedpb.Subscription.Subscribe:output_type -> google.protobuf.Empty
	2, // 3: feedpb.Subscription.Unsubscribe:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_feedpb_subscription_proto_init() }
func file_feedpb_subscription_proto_init() {
	if File_feedpb_subscription_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_feedpb_subscription_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_feedpb_subscription_proto_goTypes,
		DependencyIndexes: file_feedpb_subscription_proto_depIdxs,
		MessageInfos:      file_feedpb_subscription_proto_msgTypes,
	}.Build()
	File_feedpb_subscription_proto = out.File
	file_feedpb_subscription_proto_rawDesc = nil
	file_feedpb_subscription_proto_goTypes = nil
	file_feedpb_subscription_proto_depIdxs = nil
}
