// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: permissions/messages/v1/role.proto

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

type Role struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Role) Reset() {
	*x = Role{}
	mi := &file_permissions_messages_v1_role_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_permissions_messages_v1_role_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_permissions_messages_v1_role_proto_rawDescGZIP(), []int{0}
}

func (x *Role) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Role) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Role) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var File_permissions_messages_v1_role_proto protoreflect.FileDescriptor

const file_permissions_messages_v1_role_proto_rawDesc = "" +
	"\n" +
	"\"permissions/messages/v1/role.proto\x12\x17permissions.messages.v1\"L\n" +
	"\x04Role\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescriptionB\xfa\x01\n" +
	"\x1bcom.permissions.messages.v1B\tRoleProtoP\x01ZRgithub.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1;messagesv1\xa2\x02\x03PMX\xaa\x02\x17Permissions.Messages.V1\xca\x02\x17Permissions\\Messages\\V1\xe2\x02#Permissions\\Messages\\V1\\GPBMetadata\xea\x02\x19Permissions::Messages::V1b\x06proto3"

var (
	file_permissions_messages_v1_role_proto_rawDescOnce sync.Once
	file_permissions_messages_v1_role_proto_rawDescData []byte
)

func file_permissions_messages_v1_role_proto_rawDescGZIP() []byte {
	file_permissions_messages_v1_role_proto_rawDescOnce.Do(func() {
		file_permissions_messages_v1_role_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_permissions_messages_v1_role_proto_rawDesc), len(file_permissions_messages_v1_role_proto_rawDesc)))
	})
	return file_permissions_messages_v1_role_proto_rawDescData
}

var file_permissions_messages_v1_role_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_permissions_messages_v1_role_proto_goTypes = []any{
	(*Role)(nil), // 0: permissions.messages.v1.Role
}
var file_permissions_messages_v1_role_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_permissions_messages_v1_role_proto_init() }
func file_permissions_messages_v1_role_proto_init() {
	if File_permissions_messages_v1_role_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_permissions_messages_v1_role_proto_rawDesc), len(file_permissions_messages_v1_role_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_permissions_messages_v1_role_proto_goTypes,
		DependencyIndexes: file_permissions_messages_v1_role_proto_depIdxs,
		MessageInfos:      file_permissions_messages_v1_role_proto_msgTypes,
	}.Build()
	File_permissions_messages_v1_role_proto = out.File
	file_permissions_messages_v1_role_proto_goTypes = nil
	file_permissions_messages_v1_role_proto_depIdxs = nil
}
