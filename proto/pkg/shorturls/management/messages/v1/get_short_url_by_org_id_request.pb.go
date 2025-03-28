// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: shorturls/management/messages/v1/get_short_url_by_org_id_request.proto

package messagesv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetShortURLByOrganizationIDRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	OrganizationId string                 `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	Domain         string                 `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Alias          string                 `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *GetShortURLByOrganizationIDRequest) Reset() {
	*x = GetShortURLByOrganizationIDRequest{}
	mi := &file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetShortURLByOrganizationIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShortURLByOrganizationIDRequest) ProtoMessage() {}

func (x *GetShortURLByOrganizationIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShortURLByOrganizationIDRequest.ProtoReflect.Descriptor instead.
func (*GetShortURLByOrganizationIDRequest) Descriptor() ([]byte, []int) {
	return file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescGZIP(), []int{0}
}

func (x *GetShortURLByOrganizationIDRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *GetShortURLByOrganizationIDRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *GetShortURLByOrganizationIDRequest) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

var File_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto protoreflect.FileDescriptor

var file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDesc = string([]byte{
	0x0a, 0x46, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f,
	0x76, 0x31, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c,
	0x5f, 0x62, 0x79, 0x5f, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x20, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75,
	0x72, 0x6c, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x01, 0x0a, 0x22, 0x47, 0x65, 0x74, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34,
	0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8, 0x01, 0x01, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x06, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x05, 0x61, 0x6c, 0x69,
	0x61, 0x73, 0x42, 0xc6, 0x02, 0x0a, 0x24, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x75, 0x72, 0x6c, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x1e, 0x47, 0x65, 0x74,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x42, 0x79, 0x4f, 0x72, 0x67, 0x49, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x46, 0x72, 0x65, 0x69, 0x62, 0x65,
	0x72, 0x67, 0x56, 0x6c, 0x61, 0x64, 0x2f, 0x75, 0x72, 0x6c, 0x2d, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x53, 0x4d, 0x4d,
	0xaa, 0x02, 0x20, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x73, 0x2e, 0x4d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x20, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x73, 0x5c,
	0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x2c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72,
	0x6c, 0x73, 0x5c, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5c, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x23, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c,
	0x73, 0x3a, 0x3a, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x3a, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescOnce sync.Once
	file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescData []byte
)

func file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescGZIP() []byte {
	file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescOnce.Do(func() {
		file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDesc), len(file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDesc)))
	})
	return file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDescData
}

var file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_goTypes = []any{
	(*GetShortURLByOrganizationIDRequest)(nil), // 0: shorturls.management.messages.v1.GetShortURLByOrganizationIDRequest
}
var file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_init() }
func file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_init() {
	if File_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDesc), len(file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_goTypes,
		DependencyIndexes: file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_depIdxs,
		MessageInfos:      file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_msgTypes,
	}.Build()
	File_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto = out.File
	file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_goTypes = nil
	file_shorturls_management_messages_v1_get_short_url_by_org_id_request_proto_depIdxs = nil
}
