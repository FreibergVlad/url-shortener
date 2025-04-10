// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: domains/messages/v1/list_organization_domain_response.proto

package messagesv1

import (
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

type ListOrganizationDomainResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*Domain              `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Total         int64                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListOrganizationDomainResponse) Reset() {
	*x = ListOrganizationDomainResponse{}
	mi := &file_domains_messages_v1_list_organization_domain_response_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListOrganizationDomainResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrganizationDomainResponse) ProtoMessage() {}

func (x *ListOrganizationDomainResponse) ProtoReflect() protoreflect.Message {
	mi := &file_domains_messages_v1_list_organization_domain_response_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrganizationDomainResponse.ProtoReflect.Descriptor instead.
func (*ListOrganizationDomainResponse) Descriptor() ([]byte, []int) {
	return file_domains_messages_v1_list_organization_domain_response_proto_rawDescGZIP(), []int{0}
}

func (x *ListOrganizationDomainResponse) GetData() []*Domain {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListOrganizationDomainResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_domains_messages_v1_list_organization_domain_response_proto protoreflect.FileDescriptor

const file_domains_messages_v1_list_organization_domain_response_proto_rawDesc = "" +
	"\n" +
	";domains/messages/v1/list_organization_domain_response.proto\x12\x13domains.messages.v1\x1a domains/messages/v1/domain.proto\"g\n" +
	"\x1eListOrganizationDomainResponse\x12/\n" +
	"\x04data\x18\x01 \x03(\v2\x1b.domains.messages.v1.DomainR\x04data\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x03R\x05totalB\xfc\x01\n" +
	"\x17com.domains.messages.v1B#ListOrganizationDomainResponseProtoP\x01ZNgithub.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1;messagesv1\xa2\x02\x03DMX\xaa\x02\x13Domains.Messages.V1\xca\x02\x13Domains\\Messages\\V1\xe2\x02\x1fDomains\\Messages\\V1\\GPBMetadata\xea\x02\x15Domains::Messages::V1b\x06proto3"

var (
	file_domains_messages_v1_list_organization_domain_response_proto_rawDescOnce sync.Once
	file_domains_messages_v1_list_organization_domain_response_proto_rawDescData []byte
)

func file_domains_messages_v1_list_organization_domain_response_proto_rawDescGZIP() []byte {
	file_domains_messages_v1_list_organization_domain_response_proto_rawDescOnce.Do(func() {
		file_domains_messages_v1_list_organization_domain_response_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_domains_messages_v1_list_organization_domain_response_proto_rawDesc), len(file_domains_messages_v1_list_organization_domain_response_proto_rawDesc)))
	})
	return file_domains_messages_v1_list_organization_domain_response_proto_rawDescData
}

var file_domains_messages_v1_list_organization_domain_response_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_domains_messages_v1_list_organization_domain_response_proto_goTypes = []any{
	(*ListOrganizationDomainResponse)(nil), // 0: domains.messages.v1.ListOrganizationDomainResponse
	(*Domain)(nil),                         // 1: domains.messages.v1.Domain
}
var file_domains_messages_v1_list_organization_domain_response_proto_depIdxs = []int32{
	1, // 0: domains.messages.v1.ListOrganizationDomainResponse.data:type_name -> domains.messages.v1.Domain
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_domains_messages_v1_list_organization_domain_response_proto_init() }
func file_domains_messages_v1_list_organization_domain_response_proto_init() {
	if File_domains_messages_v1_list_organization_domain_response_proto != nil {
		return
	}
	file_domains_messages_v1_domain_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_domains_messages_v1_list_organization_domain_response_proto_rawDesc), len(file_domains_messages_v1_list_organization_domain_response_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domains_messages_v1_list_organization_domain_response_proto_goTypes,
		DependencyIndexes: file_domains_messages_v1_list_organization_domain_response_proto_depIdxs,
		MessageInfos:      file_domains_messages_v1_list_organization_domain_response_proto_msgTypes,
	}.Build()
	File_domains_messages_v1_list_organization_domain_response_proto = out.File
	file_domains_messages_v1_list_organization_domain_response_proto_goTypes = nil
	file_domains_messages_v1_list_organization_domain_response_proto_depIdxs = nil
}
