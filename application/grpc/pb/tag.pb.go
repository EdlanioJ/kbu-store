// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: protofiles/tag.proto

package pb

import (
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

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag   string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_tag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_tag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_protofiles_tag_proto_rawDescGZIP(), []int{0}
}

func (x *Tag) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Tag) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type TagListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit int32  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Sort  string `protobuf:"bytes,3,opt,name=sort,proto3" json:"sort,omitempty"`
}

func (x *TagListRequest) Reset() {
	*x = TagListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_tag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagListRequest) ProtoMessage() {}

func (x *TagListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_tag_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagListRequest.ProtoReflect.Descriptor instead.
func (*TagListRequest) Descriptor() ([]byte, []int) {
	return file_protofiles_tag_proto_rawDescGZIP(), []int{1}
}

func (x *TagListRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *TagListRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TagListRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

type TagListByCategoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CategoryId string `protobuf:"bytes,1,opt,name=categoryId,proto3" json:"categoryId,omitempty"`
	Page       int32  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Limit      int32  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Sort       string `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
}

func (x *TagListByCategoryRequest) Reset() {
	*x = TagListByCategoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_tag_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagListByCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagListByCategoryRequest) ProtoMessage() {}

func (x *TagListByCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_tag_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagListByCategoryRequest.ProtoReflect.Descriptor instead.
func (*TagListByCategoryRequest) Descriptor() ([]byte, []int) {
	return file_protofiles_tag_proto_rawDescGZIP(), []int{2}
}

func (x *TagListByCategoryRequest) GetCategoryId() string {
	if x != nil {
		return x.CategoryId
	}
	return ""
}

func (x *TagListByCategoryRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *TagListByCategoryRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TagListByCategoryRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

type TagListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tags  []*Tag `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
	Total int64  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *TagListResponse) Reset() {
	*x = TagListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_tag_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagListResponse) ProtoMessage() {}

func (x *TagListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_tag_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagListResponse.ProtoReflect.Descriptor instead.
func (*TagListResponse) Descriptor() ([]byte, []int) {
	return file_protofiles_tag_proto_rawDescGZIP(), []int{3}
}

func (x *TagListResponse) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *TagListResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_protofiles_tag_proto protoreflect.FileDescriptor

var file_protofiles_tag_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x74, 0x61, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6b, 0x62, 0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x22, 0x2d, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x4e, 0x0a, 0x0e, 0x54, 0x61, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74,
	0x22, 0x78, 0x0a, 0x18, 0x54, 0x61, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x79, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x22, 0x4b, 0x0a, 0x0f, 0x54, 0x61,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6b, 0x62,
	0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x32, 0xa6, 0x01, 0x0a, 0x0a, 0x54, 0x61, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x12, 0x19, 0x2e, 0x6b, 0x62, 0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x61, 0x67,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6b, 0x62,
	0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x61, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x6c, 0x42, 0x79, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x23, 0x2e,
	0x6b, 0x62, 0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x61, 0x67, 0x4c, 0x69, 0x73,
	0x74, 0x42, 0x79, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6b, 0x62, 0x75, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x54,
	0x61, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x15, 0x5a, 0x13, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protofiles_tag_proto_rawDescOnce sync.Once
	file_protofiles_tag_proto_rawDescData = file_protofiles_tag_proto_rawDesc
)

func file_protofiles_tag_proto_rawDescGZIP() []byte {
	file_protofiles_tag_proto_rawDescOnce.Do(func() {
		file_protofiles_tag_proto_rawDescData = protoimpl.X.CompressGZIP(file_protofiles_tag_proto_rawDescData)
	})
	return file_protofiles_tag_proto_rawDescData
}

var file_protofiles_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protofiles_tag_proto_goTypes = []interface{}{
	(*Tag)(nil),                      // 0: kbu_store.Tag
	(*TagListRequest)(nil),           // 1: kbu_store.TagListRequest
	(*TagListByCategoryRequest)(nil), // 2: kbu_store.TagListByCategoryRequest
	(*TagListResponse)(nil),          // 3: kbu_store.TagListResponse
}
var file_protofiles_tag_proto_depIdxs = []int32{
	0, // 0: kbu_store.TagListResponse.tags:type_name -> kbu_store.Tag
	1, // 1: kbu_store.TagService.GetAll:input_type -> kbu_store.TagListRequest
	2, // 2: kbu_store.TagService.GetAllByCategory:input_type -> kbu_store.TagListByCategoryRequest
	3, // 3: kbu_store.TagService.GetAll:output_type -> kbu_store.TagListResponse
	3, // 4: kbu_store.TagService.GetAllByCategory:output_type -> kbu_store.TagListResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protofiles_tag_proto_init() }
func file_protofiles_tag_proto_init() {
	if File_protofiles_tag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protofiles_tag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_protofiles_tag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagListRequest); i {
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
		file_protofiles_tag_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagListByCategoryRequest); i {
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
		file_protofiles_tag_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagListResponse); i {
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
			RawDescriptor: file_protofiles_tag_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protofiles_tag_proto_goTypes,
		DependencyIndexes: file_protofiles_tag_proto_depIdxs,
		MessageInfos:      file_protofiles_tag_proto_msgTypes,
	}.Build()
	File_protofiles_tag_proto = out.File
	file_protofiles_tag_proto_rawDesc = nil
	file_protofiles_tag_proto_goTypes = nil
	file_protofiles_tag_proto_depIdxs = nil
}
