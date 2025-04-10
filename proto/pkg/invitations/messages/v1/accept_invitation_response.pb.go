// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: invitations/messages/v1/accept_invitation_response.proto

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

type AcceptInvitationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AcceptInvitationResponse) Reset() {
	*x = AcceptInvitationResponse{}
	mi := &file_invitations_messages_v1_accept_invitation_response_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AcceptInvitationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptInvitationResponse) ProtoMessage() {}

func (x *AcceptInvitationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invitations_messages_v1_accept_invitation_response_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptInvitationResponse.ProtoReflect.Descriptor instead.
func (*AcceptInvitationResponse) Descriptor() ([]byte, []int) {
	return file_invitations_messages_v1_accept_invitation_response_proto_rawDescGZIP(), []int{0}
}

var File_invitations_messages_v1_accept_invitation_response_proto protoreflect.FileDescriptor

const file_invitations_messages_v1_accept_invitation_response_proto_rawDesc = "" +
	"\n" +
	"8invitations/messages/v1/accept_invitation_response.proto\x12\x17invitations.messages.v1\"\x1a\n" +
	"\x18AcceptInvitationResponseB\x8e\x02\n" +
	"\x1bcom.invitations.messages.v1B\x1dAcceptInvitationResponseProtoP\x01ZRgithub.com/FreibergVlad/url-shortener/proto/pkg/invitations/messages/v1;messagesv1\xa2\x02\x03IMX\xaa\x02\x17Invitations.Messages.V1\xca\x02\x17Invitations\\Messages\\V1\xe2\x02#Invitations\\Messages\\V1\\GPBMetadata\xea\x02\x19Invitations::Messages::V1b\x06proto3"

var (
	file_invitations_messages_v1_accept_invitation_response_proto_rawDescOnce sync.Once
	file_invitations_messages_v1_accept_invitation_response_proto_rawDescData []byte
)

func file_invitations_messages_v1_accept_invitation_response_proto_rawDescGZIP() []byte {
	file_invitations_messages_v1_accept_invitation_response_proto_rawDescOnce.Do(func() {
		file_invitations_messages_v1_accept_invitation_response_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_invitations_messages_v1_accept_invitation_response_proto_rawDesc), len(file_invitations_messages_v1_accept_invitation_response_proto_rawDesc)))
	})
	return file_invitations_messages_v1_accept_invitation_response_proto_rawDescData
}

var file_invitations_messages_v1_accept_invitation_response_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_invitations_messages_v1_accept_invitation_response_proto_goTypes = []any{
	(*AcceptInvitationResponse)(nil), // 0: invitations.messages.v1.AcceptInvitationResponse
}
var file_invitations_messages_v1_accept_invitation_response_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_invitations_messages_v1_accept_invitation_response_proto_init() }
func file_invitations_messages_v1_accept_invitation_response_proto_init() {
	if File_invitations_messages_v1_accept_invitation_response_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_invitations_messages_v1_accept_invitation_response_proto_rawDesc), len(file_invitations_messages_v1_accept_invitation_response_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_invitations_messages_v1_accept_invitation_response_proto_goTypes,
		DependencyIndexes: file_invitations_messages_v1_accept_invitation_response_proto_depIdxs,
		MessageInfos:      file_invitations_messages_v1_accept_invitation_response_proto_msgTypes,
	}.Build()
	File_invitations_messages_v1_accept_invitation_response_proto = out.File
	file_invitations_messages_v1_accept_invitation_response_proto_goTypes = nil
	file_invitations_messages_v1_accept_invitation_response_proto_depIdxs = nil
}
