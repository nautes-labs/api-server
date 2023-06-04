// Copyright 2023 Nautes Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: common/v1/error_reason.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type ErrorReason int32

const (
	ErrorReason_PROJECT_NOT_FOUND                 ErrorReason = 0
	ErrorReason_GROUP_NOT_FOUND                   ErrorReason = 1
	ErrorReason_NODE_NOT_FOUND                    ErrorReason = 2
	ErrorReason_RESOURCE_NOT_FOUND                ErrorReason = 3
	ErrorReason_RESOURCE_NOT_MATCH                ErrorReason = 4
	ErrorReason_NO_AUTHORIZATION                  ErrorReason = 5
	ErrorReason_DEPLOYKEY_NOT_FOUND               ErrorReason = 6
	ErrorReason_SECRET_NOT_FOUND                  ErrorReason = 7
	ErrorReason_ACCESSTOKEN_NOT_FOUND             ErrorReason = 8
	ErrorReason_REFRESH_PERMISSIONS_ACCESS_DENIED ErrorReason = 9
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0: "PROJECT_NOT_FOUND",
		1: "GROUP_NOT_FOUND",
		2: "NODE_NOT_FOUND",
		3: "RESOURCE_NOT_FOUND",
		4: "RESOURCE_NOT_MATCH",
		5: "NO_AUTHORIZATION",
		6: "DEPLOYKEY_NOT_FOUND",
		7: "SECRET_NOT_FOUND",
		8: "ACCESSTOKEN_NOT_FOUND",
		9: "REFRESH_PERMISSIONS_ACCESS_DENIED",
	}
	ErrorReason_value = map[string]int32{
		"PROJECT_NOT_FOUND":                 0,
		"GROUP_NOT_FOUND":                   1,
		"NODE_NOT_FOUND":                    2,
		"RESOURCE_NOT_FOUND":                3,
		"RESOURCE_NOT_MATCH":                4,
		"NO_AUTHORIZATION":                  5,
		"DEPLOYKEY_NOT_FOUND":               6,
		"SECRET_NOT_FOUND":                  7,
		"ACCESSTOKEN_NOT_FOUND":             8,
		"REFRESH_PERMISSIONS_ACCESS_DENIED": 9,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_error_reason_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_common_v1_error_reason_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_error_reason_proto_rawDescGZIP(), []int{0}
}

var File_common_v1_error_reason_proto protoreflect.FileDescriptor

var file_common_v1_error_reason_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xc6,
	0x02, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1b,
	0x0a, 0x11, 0x50, 0x52, 0x4f, 0x4a, 0x45, 0x43, 0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f,
	0x55, 0x4e, 0x44, 0x10, 0x00, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x19, 0x0a, 0x0f, 0x47,
	0x52, 0x4f, 0x55, 0x50, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01,
	0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x18, 0x0a, 0x0e, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x4e,
	0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03,
	0x12, 0x1c, 0x0a, 0x12, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x4e, 0x4f, 0x54,
	0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x03, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x1c,
	0x0a, 0x12, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x4d,
	0x41, 0x54, 0x43, 0x48, 0x10, 0x04, 0x1a, 0x04, 0xa8, 0x45, 0xf4, 0x03, 0x12, 0x1a, 0x0a, 0x10,
	0x4e, 0x4f, 0x5f, 0x41, 0x55, 0x54, 0x48, 0x4f, 0x52, 0x49, 0x5a, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x05, 0x1a, 0x04, 0xa8, 0x45, 0x93, 0x03, 0x12, 0x1d, 0x0a, 0x13, 0x44, 0x45, 0x50, 0x4c,
	0x4f, 0x59, 0x4b, 0x45, 0x59, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10,
	0x06, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x1a, 0x0a, 0x10, 0x53, 0x45, 0x43, 0x52, 0x45,
	0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x07, 0x1a, 0x04, 0xa8,
	0x45, 0x94, 0x03, 0x12, 0x1f, 0x0a, 0x15, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x54, 0x4f, 0x4b,
	0x45, 0x4e, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x08, 0x1a, 0x04,
	0xa8, 0x45, 0x94, 0x03, 0x12, 0x2b, 0x0a, 0x21, 0x52, 0x45, 0x46, 0x52, 0x45, 0x53, 0x48, 0x5f,
	0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x41, 0x43, 0x43, 0x45,
	0x53, 0x53, 0x5f, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44, 0x10, 0x09, 0x1a, 0x04, 0xa8, 0x45, 0x93,
	0x03, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x75, 0x74, 0x65, 0x73, 0x2d, 0x6c, 0x61, 0x62,
	0x73, 0x2f, 0x61, 0x70, 0x69, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_v1_error_reason_proto_rawDescOnce sync.Once
	file_common_v1_error_reason_proto_rawDescData = file_common_v1_error_reason_proto_rawDesc
)

func file_common_v1_error_reason_proto_rawDescGZIP() []byte {
	file_common_v1_error_reason_proto_rawDescOnce.Do(func() {
		file_common_v1_error_reason_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_error_reason_proto_rawDescData)
	})
	return file_common_v1_error_reason_proto_rawDescData
}

var file_common_v1_error_reason_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_v1_error_reason_proto_goTypes = []interface{}{
	(ErrorReason)(0), // 0: common.v1.ErrorReason
}
var file_common_v1_error_reason_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_v1_error_reason_proto_init() }
func file_common_v1_error_reason_proto_init() {
	if File_common_v1_error_reason_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_v1_error_reason_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1_error_reason_proto_goTypes,
		DependencyIndexes: file_common_v1_error_reason_proto_depIdxs,
		EnumInfos:         file_common_v1_error_reason_proto_enumTypes,
	}.Build()
	File_common_v1_error_reason_proto = out.File
	file_common_v1_error_reason_proto_rawDesc = nil
	file_common_v1_error_reason_proto_goTypes = nil
	file_common_v1_error_reason_proto_depIdxs = nil
}
