// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: straw-hat-llc/onepiece/extensions.proto

package onepiece

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_straw_hat_llc_onepiece_extensions_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "onepiece.protobuf.stream_id",
		Tag:           "varint,50001,opt,name=stream_id",
		Filename:      "straw-hat-llc/onepiece/extensions.proto",
	},
}

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool stream_id = 50001;
	E_StreamId = &file_straw_hat_llc_onepiece_extensions_proto_extTypes[0]
)

var File_straw_hat_llc_onepiece_extensions_proto protoreflect.FileDescriptor

var file_straw_hat_llc_onepiece_extensions_proto_rawDesc = []byte{
	0x0a, 0x27, 0x73, 0x74, 0x72, 0x61, 0x77, 0x2d, 0x68, 0x61, 0x74, 0x2d, 0x6c, 0x6c, 0x63, 0x2f,
	0x6f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6f, 0x6e, 0x65, 0x70, 0x69,
	0x65, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x20, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3c,
	0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x12, 0x1d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x42, 0xc3, 0x01, 0x0a,
	0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x6f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x42, 0x0f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x75, 0x6e, 0x73, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2f, 0x70, 0x6c, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70,
	0x6c, 0x61, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x72, 0x61, 0x77, 0x2d, 0x68,
	0x61, 0x74, 0x2d, 0x6c, 0x6c, 0x63, 0x2f, 0x6f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65, 0xa2,
	0x02, 0x03, 0x4f, 0x50, 0x58, 0xaa, 0x02, 0x11, 0x4f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0xca, 0x02, 0x11, 0x4f, 0x6e, 0x65, 0x70,
	0x69, 0x65, 0x63, 0x65, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0xe2, 0x02, 0x1d,
	0x4f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12,
	0x4f, 0x6e, 0x65, 0x70, 0x69, 0x65, 0x63, 0x65, 0x3a, 0x3a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_straw_hat_llc_onepiece_extensions_proto_goTypes = []interface{}{
	(*descriptorpb.FieldOptions)(nil), // 0: google.protobuf.FieldOptions
}
var file_straw_hat_llc_onepiece_extensions_proto_depIdxs = []int32{
	0, // 0: onepiece.protobuf.stream_id:extendee -> google.protobuf.FieldOptions
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_straw_hat_llc_onepiece_extensions_proto_init() }
func file_straw_hat_llc_onepiece_extensions_proto_init() {
	if File_straw_hat_llc_onepiece_extensions_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_straw_hat_llc_onepiece_extensions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_straw_hat_llc_onepiece_extensions_proto_goTypes,
		DependencyIndexes: file_straw_hat_llc_onepiece_extensions_proto_depIdxs,
		ExtensionInfos:    file_straw_hat_llc_onepiece_extensions_proto_extTypes,
	}.Build()
	File_straw_hat_llc_onepiece_extensions_proto = out.File
	file_straw_hat_llc_onepiece_extensions_proto_rawDesc = nil
	file_straw_hat_llc_onepiece_extensions_proto_goTypes = nil
	file_straw_hat_llc_onepiece_extensions_proto_depIdxs = nil
}
