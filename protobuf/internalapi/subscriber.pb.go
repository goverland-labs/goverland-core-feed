// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: subscriber.proto

package internalapi

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

type CreateSubscriberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriberId string  `protobuf:"bytes,1,opt,name=subscriber_id,json=subscriberId,proto3" json:"subscriber_id,omitempty"`
	WebhookUrl   *string `protobuf:"bytes,2,opt,name=webhook_url,json=webhookUrl,proto3,oneof" json:"webhook_url,omitempty"`
}

func (x *CreateSubscriberRequest) Reset() {
	*x = CreateSubscriberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_subscriber_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSubscriberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSubscriberRequest) ProtoMessage() {}

func (x *CreateSubscriberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_subscriber_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSubscriberRequest.ProtoReflect.Descriptor instead.
func (*CreateSubscriberRequest) Descriptor() ([]byte, []int) {
	return file_subscriber_proto_rawDescGZIP(), []int{0}
}

func (x *CreateSubscriberRequest) GetSubscriberId() string {
	if x != nil {
		return x.SubscriberId
	}
	return ""
}

func (x *CreateSubscriberRequest) GetWebhookUrl() string {
	if x != nil && x.WebhookUrl != nil {
		return *x.WebhookUrl
	}
	return ""
}

type CreateSubscriberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriberId string `protobuf:"bytes,1,opt,name=subscriber_id,json=subscriberId,proto3" json:"subscriber_id,omitempty"`
}

func (x *CreateSubscriberResponse) Reset() {
	*x = CreateSubscriberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_subscriber_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSubscriberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSubscriberResponse) ProtoMessage() {}

func (x *CreateSubscriberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_subscriber_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSubscriberResponse.ProtoReflect.Descriptor instead.
func (*CreateSubscriberResponse) Descriptor() ([]byte, []int) {
	return file_subscriber_proto_rawDescGZIP(), []int{1}
}

func (x *CreateSubscriberResponse) GetSubscriberId() string {
	if x != nil {
		return x.SubscriberId
	}
	return ""
}

type UpdateSubscriberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriberId string  `protobuf:"bytes,1,opt,name=subscriber_id,json=subscriberId,proto3" json:"subscriber_id,omitempty"`
	WebhookUrl   *string `protobuf:"bytes,2,opt,name=webhook_url,json=webhookUrl,proto3,oneof" json:"webhook_url,omitempty"`
}

func (x *UpdateSubscriberRequest) Reset() {
	*x = UpdateSubscriberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_subscriber_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSubscriberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSubscriberRequest) ProtoMessage() {}

func (x *UpdateSubscriberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_subscriber_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSubscriberRequest.ProtoReflect.Descriptor instead.
func (*UpdateSubscriberRequest) Descriptor() ([]byte, []int) {
	return file_subscriber_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateSubscriberRequest) GetSubscriberId() string {
	if x != nil {
		return x.SubscriberId
	}
	return ""
}

func (x *UpdateSubscriberRequest) GetWebhookUrl() string {
	if x != nil && x.WebhookUrl != nil {
		return *x.WebhookUrl
	}
	return ""
}

var File_subscriber_proto protoreflect.FileDescriptor

var file_subscriber_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x61, 0x70, 0x69, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x74, 0x0a, 0x17,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0b,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0a, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x55, 0x72, 0x6c, 0x88,
	0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x75,
	0x72, 0x6c, 0x22, 0x3f, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x74, 0x0a, 0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0b, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x77, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x32, 0xab, 0x01, 0x0a, 0x0a, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x24, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x61, 0x70, 0x69,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x46, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x24, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x3b, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_subscriber_proto_rawDescOnce sync.Once
	file_subscriber_proto_rawDescData = file_subscriber_proto_rawDesc
)

func file_subscriber_proto_rawDescGZIP() []byte {
	file_subscriber_proto_rawDescOnce.Do(func() {
		file_subscriber_proto_rawDescData = protoimpl.X.CompressGZIP(file_subscriber_proto_rawDescData)
	})
	return file_subscriber_proto_rawDescData
}

var file_subscriber_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_subscriber_proto_goTypes = []interface{}{
	(*CreateSubscriberRequest)(nil),  // 0: internalapi.CreateSubscriberRequest
	(*CreateSubscriberResponse)(nil), // 1: internalapi.CreateSubscriberResponse
	(*UpdateSubscriberRequest)(nil),  // 2: internalapi.UpdateSubscriberRequest
	(*emptypb.Empty)(nil),            // 3: google.protobuf.Empty
}
var file_subscriber_proto_depIdxs = []int32{
	0, // 0: internalapi.Subscriber.Create:input_type -> internalapi.CreateSubscriberRequest
	2, // 1: internalapi.Subscriber.Update:input_type -> internalapi.UpdateSubscriberRequest
	1, // 2: internalapi.Subscriber.Create:output_type -> internalapi.CreateSubscriberResponse
	3, // 3: internalapi.Subscriber.Update:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_subscriber_proto_init() }
func file_subscriber_proto_init() {
	if File_subscriber_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_subscriber_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSubscriberRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_subscriber_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSubscriberResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_subscriber_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSubscriberRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_subscriber_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_subscriber_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_subscriber_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_subscriber_proto_goTypes,
		DependencyIndexes: file_subscriber_proto_depIdxs,
		MessageInfos:      file_subscriber_proto_msgTypes,
	}.Build()
	File_subscriber_proto = out.File
	file_subscriber_proto_rawDesc = nil
	file_subscriber_proto_goTypes = nil
	file_subscriber_proto_depIdxs = nil
}
