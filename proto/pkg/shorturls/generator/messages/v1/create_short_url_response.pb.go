// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: shorturls/generator/messages/v1/create_short_url_response.proto

package messagesv1

import (
	v1 "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
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

type CreateShortURLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      *v1.ShortURL           `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateShortURLResponse) Reset() {
	*x = CreateShortURLResponse{}
	mi := &file_shorturls_generator_messages_v1_create_short_url_response_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateShortURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLResponse) ProtoMessage() {}

func (x *CreateShortURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shorturls_generator_messages_v1_create_short_url_response_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLResponse.ProtoReflect.Descriptor instead.
func (*CreateShortURLResponse) Descriptor() ([]byte, []int) {
	return file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescGZIP(), []int{0}
}

func (x *CreateShortURLResponse) GetShortUrl() *v1.ShortURL {
	if x != nil {
		return x.ShortUrl
	}
	return nil
}

var File_shorturls_generator_messages_v1_create_short_url_response_proto protoreflect.FileDescriptor

const file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDesc = "" +
	"\n" +
	"?shorturls/generator/messages/v1/create_short_url_response.proto\x12\x1fshorturls.generator.messages.v1\x1a0shorturls/management/messages/v1/short_url.proto\"a\n" +
	"\x16CreateShortURLResponse\x12G\n" +
	"\tshort_url\x18\x01 \x01(\v2*.shorturls.management.messages.v1.ShortURLR\bshortUrlB\xbf\x02\n" +
	"#com.shorturls.generator.messages.v1B\x1bCreateShortUrlResponseProtoP\x01ZZgithub.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1;messagesv1\xa2\x02\x03SGM\xaa\x02\x1fShorturls.Generator.Messages.V1\xca\x02 Shorturls\\Generator_\\Messages\\V1\xe2\x02,Shorturls\\Generator_\\Messages\\V1\\GPBMetadata\xea\x02\"Shorturls::Generator::Messages::V1b\x06proto3"

var (
	file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescOnce sync.Once
	file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescData []byte
)

func file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescGZIP() []byte {
	file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescOnce.Do(func() {
		file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDesc), len(file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDesc)))
	})
	return file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDescData
}

var file_shorturls_generator_messages_v1_create_short_url_response_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shorturls_generator_messages_v1_create_short_url_response_proto_goTypes = []any{
	(*CreateShortURLResponse)(nil), // 0: shorturls.generator.messages.v1.CreateShortURLResponse
	(*v1.ShortURL)(nil),            // 1: shorturls.management.messages.v1.ShortURL
}
var file_shorturls_generator_messages_v1_create_short_url_response_proto_depIdxs = []int32{
	1, // 0: shorturls.generator.messages.v1.CreateShortURLResponse.short_url:type_name -> shorturls.management.messages.v1.ShortURL
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_shorturls_generator_messages_v1_create_short_url_response_proto_init() }
func file_shorturls_generator_messages_v1_create_short_url_response_proto_init() {
	if File_shorturls_generator_messages_v1_create_short_url_response_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDesc), len(file_shorturls_generator_messages_v1_create_short_url_response_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shorturls_generator_messages_v1_create_short_url_response_proto_goTypes,
		DependencyIndexes: file_shorturls_generator_messages_v1_create_short_url_response_proto_depIdxs,
		MessageInfos:      file_shorturls_generator_messages_v1_create_short_url_response_proto_msgTypes,
	}.Build()
	File_shorturls_generator_messages_v1_create_short_url_response_proto = out.File
	file_shorturls_generator_messages_v1_create_short_url_response_proto_goTypes = nil
	file_shorturls_generator_messages_v1_create_short_url_response_proto_depIdxs = nil
}
