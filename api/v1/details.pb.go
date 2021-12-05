// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: details.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type BadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Violations []*BadRequest_FieldViolation `protobuf:"bytes,3,rep,name=violations,proto3" json:"violations,omitempty"`
}

func (x *BadRequest) Reset() {
	*x = BadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_details_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BadRequest) ProtoMessage() {}

func (x *BadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_details_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BadRequest.ProtoReflect.Descriptor instead.
func (*BadRequest) Descriptor() ([]byte, []int) {
	return file_details_proto_rawDescGZIP(), []int{0}
}

func (x *BadRequest) GetViolations() []*BadRequest_FieldViolation {
	if x != nil {
		return x.Violations
	}
	return nil
}

type PreconditionFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Violations []*PreconditionFailure_PreconditionViolation `protobuf:"bytes,3,rep,name=violations,proto3" json:"violations,omitempty"`
}

func (x *PreconditionFailure) Reset() {
	*x = PreconditionFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_details_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreconditionFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreconditionFailure) ProtoMessage() {}

func (x *PreconditionFailure) ProtoReflect() protoreflect.Message {
	mi := &file_details_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreconditionFailure.ProtoReflect.Descriptor instead.
func (*PreconditionFailure) Descriptor() ([]byte, []int) {
	return file_details_proto_rawDescGZIP(), []int{1}
}

func (x *PreconditionFailure) GetViolations() []*PreconditionFailure_PreconditionViolation {
	if x != nil {
		return x.Violations
	}
	return nil
}

type BadRequest_FieldViolation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field       string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *BadRequest_FieldViolation) Reset() {
	*x = BadRequest_FieldViolation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_details_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BadRequest_FieldViolation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BadRequest_FieldViolation) ProtoMessage() {}

func (x *BadRequest_FieldViolation) ProtoReflect() protoreflect.Message {
	mi := &file_details_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BadRequest_FieldViolation.ProtoReflect.Descriptor instead.
func (*BadRequest_FieldViolation) Descriptor() ([]byte, []int) {
	return file_details_proto_rawDescGZIP(), []int{0, 0}
}

func (x *BadRequest_FieldViolation) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *BadRequest_FieldViolation) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type PreconditionFailure_PreconditionViolation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *PreconditionFailure_PreconditionViolation) Reset() {
	*x = PreconditionFailure_PreconditionViolation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_details_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreconditionFailure_PreconditionViolation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreconditionFailure_PreconditionViolation) ProtoMessage() {}

func (x *PreconditionFailure_PreconditionViolation) ProtoReflect() protoreflect.Message {
	mi := &file_details_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreconditionFailure_PreconditionViolation.ProtoReflect.Descriptor instead.
func (*PreconditionFailure_PreconditionViolation) Descriptor() ([]byte, []int) {
	return file_details_proto_rawDescGZIP(), []int{1, 0}
}

func (x *PreconditionFailure_PreconditionViolation) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PreconditionFailure_PreconditionViolation) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var File_details_proto protoreflect.FileDescriptor

var file_details_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0f, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x63, 0x6e, 0x63, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x22, 0xa2, 0x01, 0x0a, 0x0a, 0x42, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x4a, 0x0a, 0x0a, 0x76, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x63, 0x6e, 0x63, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x76, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x48, 0x0a, 0x0e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xc0, 0x01, 0x0a, 0x13, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x5a, 0x0a,
	0x0a, 0x76, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x3a, 0x2e, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x63, 0x6e, 0x63, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x76,
	0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x4d, 0x0a, 0x15, 0x50, 0x72, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x26, 0x5a, 0x24, 0x74, 0x68, 0x65, 0x74,
	0x65, 0x63, 0x68, 0x6e, 0x69, 0x63, 0x6b, 0x2e, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x6c, 0x69,
	0x6e, 0x75, 0x78, 0x63, 0x6e, 0x63, 0x2d, 0x75, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_details_proto_rawDescOnce sync.Once
	file_details_proto_rawDescData = file_details_proto_rawDesc
)

func file_details_proto_rawDescGZIP() []byte {
	file_details_proto_rawDescOnce.Do(func() {
		file_details_proto_rawDescData = protoimpl.X.CompressGZIP(file_details_proto_rawDescData)
	})
	return file_details_proto_rawDescData
}

var file_details_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_details_proto_goTypes = []interface{}{
	(*BadRequest)(nil),                                // 0: linuxcnc.api.v1.BadRequest
	(*PreconditionFailure)(nil),                       // 1: linuxcnc.api.v1.PreconditionFailure
	(*BadRequest_FieldViolation)(nil),                 // 2: linuxcnc.api.v1.BadRequest.FieldViolation
	(*PreconditionFailure_PreconditionViolation)(nil), // 3: linuxcnc.api.v1.PreconditionFailure.PreconditionViolation
}
var file_details_proto_depIdxs = []int32{
	2, // 0: linuxcnc.api.v1.BadRequest.violations:type_name -> linuxcnc.api.v1.BadRequest.FieldViolation
	3, // 1: linuxcnc.api.v1.PreconditionFailure.violations:type_name -> linuxcnc.api.v1.PreconditionFailure.PreconditionViolation
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_details_proto_init() }
func file_details_proto_init() {
	if File_details_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_details_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BadRequest); i {
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
		file_details_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreconditionFailure); i {
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
		file_details_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BadRequest_FieldViolation); i {
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
		file_details_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreconditionFailure_PreconditionViolation); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_details_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_details_proto_goTypes,
		DependencyIndexes: file_details_proto_depIdxs,
		MessageInfos:      file_details_proto_msgTypes,
	}.Build()
	File_details_proto = out.File
	file_details_proto_rawDesc = nil
	file_details_proto_goTypes = nil
	file_details_proto_depIdxs = nil
}