// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.5
// source: shared/protobuf/proto/equipment.proto

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

// ------------ 强化装备 ------------
type C2SEquipmentStrengthen struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                      // 强化装备的唯一ID
	Materials []int64 `protobuf:"varint,2,rep,packed,name=materials,proto3" json:"materials,omitempty"` // 素材装备的唯一ID
	Items     []int32 `protobuf:"varint,3,rep,packed,name=items,proto3" json:"items,omitempty"`         // 经验道具数量 长度4，从低到高
}

func (x *C2SEquipmentStrengthen) Reset() {
	*x = C2SEquipmentStrengthen{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentStrengthen) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentStrengthen) ProtoMessage() {}

func (x *C2SEquipmentStrengthen) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentStrengthen.ProtoReflect.Descriptor instead.
func (*C2SEquipmentStrengthen) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{0}
}

func (x *C2SEquipmentStrengthen) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *C2SEquipmentStrengthen) GetMaterials() []int64 {
	if x != nil {
		return x.Materials
	}
	return nil
}

func (x *C2SEquipmentStrengthen) GetItems() []int32 {
	if x != nil {
		return x.Items
	}
	return nil
}

type S2CEquipmentStrengthen struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equipment      *VOUserEquipment  `protobuf:"bytes,1,opt,name=equipment,proto3" json:"equipment,omitempty"`           // 装备数据
	ResourceResult *VOResourceResult `protobuf:"bytes,2,opt,name=resourceResult,proto3" json:"resourceResult,omitempty"` // 玩家资源
	Materials      []int64           `protobuf:"varint,3,rep,packed,name=materials,proto3" json:"materials,omitempty"`   // 素材装备的唯一ID
}

func (x *S2CEquipmentStrengthen) Reset() {
	*x = S2CEquipmentStrengthen{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentStrengthen) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentStrengthen) ProtoMessage() {}

func (x *S2CEquipmentStrengthen) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentStrengthen.ProtoReflect.Descriptor instead.
func (*S2CEquipmentStrengthen) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{1}
}

func (x *S2CEquipmentStrengthen) GetEquipment() *VOUserEquipment {
	if x != nil {
		return x.Equipment
	}
	return nil
}

func (x *S2CEquipmentStrengthen) GetResourceResult() *VOResourceResult {
	if x != nil {
		return x.ResourceResult
	}
	return nil
}

func (x *S2CEquipmentStrengthen) GetMaterials() []int64 {
	if x != nil {
		return x.Materials
	}
	return nil
}

// ------------ 突破装备 ------------
type C2SEquipmentAdvance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                      // 突破装备的唯一ID
	Materials []int64 `protobuf:"varint,2,rep,packed,name=materials,proto3" json:"materials,omitempty"` // 素材装备的唯一ID
}

func (x *C2SEquipmentAdvance) Reset() {
	*x = C2SEquipmentAdvance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentAdvance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentAdvance) ProtoMessage() {}

func (x *C2SEquipmentAdvance) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentAdvance.ProtoReflect.Descriptor instead.
func (*C2SEquipmentAdvance) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{2}
}

func (x *C2SEquipmentAdvance) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *C2SEquipmentAdvance) GetMaterials() []int64 {
	if x != nil {
		return x.Materials
	}
	return nil
}

type S2CEquipmentAdvance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equipment      *VOUserEquipment  `protobuf:"bytes,1,opt,name=equipment,proto3" json:"equipment,omitempty"`           // 装备数据
	ResourceResult *VOResourceResult `protobuf:"bytes,2,opt,name=resourceResult,proto3" json:"resourceResult,omitempty"` // 玩家资源
	Materials      []int64           `protobuf:"varint,3,rep,packed,name=materials,proto3" json:"materials,omitempty"`   // 素材装备的唯一ID
}

func (x *S2CEquipmentAdvance) Reset() {
	*x = S2CEquipmentAdvance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentAdvance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentAdvance) ProtoMessage() {}

func (x *S2CEquipmentAdvance) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentAdvance.ProtoReflect.Descriptor instead.
func (*S2CEquipmentAdvance) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{3}
}

func (x *S2CEquipmentAdvance) GetEquipment() *VOUserEquipment {
	if x != nil {
		return x.Equipment
	}
	return nil
}

func (x *S2CEquipmentAdvance) GetResourceResult() *VOResourceResult {
	if x != nil {
		return x.ResourceResult
	}
	return nil
}

func (x *S2CEquipmentAdvance) GetMaterials() []int64 {
	if x != nil {
		return x.Materials
	}
	return nil
}

// ------------ 装备阵营重铸 ------------
type C2SEquipmentRecastCamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 装备唯一ID
}

func (x *C2SEquipmentRecastCamp) Reset() {
	*x = C2SEquipmentRecastCamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentRecastCamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentRecastCamp) ProtoMessage() {}

func (x *C2SEquipmentRecastCamp) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentRecastCamp.ProtoReflect.Descriptor instead.
func (*C2SEquipmentRecastCamp) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{4}
}

func (x *C2SEquipmentRecastCamp) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type S2CEquipmentRecastCamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Camp           int32             `protobuf:"varint,1,opt,name=camp,proto3" json:"camp,omitempty"`                    // 装备阵营
	ResourceResult *VOResourceResult `protobuf:"bytes,2,opt,name=resourceResult,proto3" json:"resourceResult,omitempty"` // 玩家资源
}

func (x *S2CEquipmentRecastCamp) Reset() {
	*x = S2CEquipmentRecastCamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentRecastCamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentRecastCamp) ProtoMessage() {}

func (x *S2CEquipmentRecastCamp) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentRecastCamp.ProtoReflect.Descriptor instead.
func (*S2CEquipmentRecastCamp) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{5}
}

func (x *S2CEquipmentRecastCamp) GetCamp() int32 {
	if x != nil {
		return x.Camp
	}
	return 0
}

func (x *S2CEquipmentRecastCamp) GetResourceResult() *VOResourceResult {
	if x != nil {
		return x.ResourceResult
	}
	return nil
}

// ------------ 装备阵营重铸确认 ------------
type C2SEquipmentConfirmRecastCamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`           // 装备唯一ID
	Confirm bool  `protobuf:"varint,2,opt,name=confirm,proto3" json:"confirm,omitempty"` // 是否保留重铸结果
}

func (x *C2SEquipmentConfirmRecastCamp) Reset() {
	*x = C2SEquipmentConfirmRecastCamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentConfirmRecastCamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentConfirmRecastCamp) ProtoMessage() {}

func (x *C2SEquipmentConfirmRecastCamp) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentConfirmRecastCamp.ProtoReflect.Descriptor instead.
func (*C2SEquipmentConfirmRecastCamp) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{6}
}

func (x *C2SEquipmentConfirmRecastCamp) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *C2SEquipmentConfirmRecastCamp) GetConfirm() bool {
	if x != nil {
		return x.Confirm
	}
	return false
}

type S2CEquipmentConfirmRecastCamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equipment *VOUserEquipment `protobuf:"bytes,1,opt,name=equipment,proto3" json:"equipment,omitempty"` // 装备数据
}

func (x *S2CEquipmentConfirmRecastCamp) Reset() {
	*x = S2CEquipmentConfirmRecastCamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentConfirmRecastCamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentConfirmRecastCamp) ProtoMessage() {}

func (x *S2CEquipmentConfirmRecastCamp) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentConfirmRecastCamp.ProtoReflect.Descriptor instead.
func (*S2CEquipmentConfirmRecastCamp) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{7}
}

func (x *S2CEquipmentConfirmRecastCamp) GetEquipment() *VOUserEquipment {
	if x != nil {
		return x.Equipment
	}
	return nil
}

// ------------ 装备加锁和解锁 ------------
type C2SEquipmentLock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`         // 装备唯一ID
	IsLock bool  `protobuf:"varint,2,opt,name=isLock,proto3" json:"isLock,omitempty"` // 是否加锁
}

func (x *C2SEquipmentLock) Reset() {
	*x = C2SEquipmentLock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentLock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentLock) ProtoMessage() {}

func (x *C2SEquipmentLock) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentLock.ProtoReflect.Descriptor instead.
func (*C2SEquipmentLock) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{8}
}

func (x *C2SEquipmentLock) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *C2SEquipmentLock) GetIsLock() bool {
	if x != nil {
		return x.IsLock
	}
	return false
}

type S2CEquipmentLock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`         // 装备唯一ID
	IsLock bool  `protobuf:"varint,2,opt,name=isLock,proto3" json:"isLock,omitempty"` // 是否加锁
}

func (x *S2CEquipmentLock) Reset() {
	*x = S2CEquipmentLock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentLock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentLock) ProtoMessage() {}

func (x *S2CEquipmentLock) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentLock.ProtoReflect.Descriptor instead.
func (*S2CEquipmentLock) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{9}
}

func (x *S2CEquipmentLock) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *S2CEquipmentLock) GetIsLock() bool {
	if x != nil {
		return x.IsLock
	}
	return false
}

// ------------ 获取装备列表 ------------
type C2SEquipmentList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *C2SEquipmentList) Reset() {
	*x = C2SEquipmentList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SEquipmentList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SEquipmentList) ProtoMessage() {}

func (x *C2SEquipmentList) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SEquipmentList.ProtoReflect.Descriptor instead.
func (*C2SEquipmentList) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{10}
}

type S2CEquipmentList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equipments []*VOUserEquipment `protobuf:"bytes,1,rep,name=equipments,proto3" json:"equipments,omitempty"` // 装备列表数据
}

func (x *S2CEquipmentList) Reset() {
	*x = S2CEquipmentList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CEquipmentList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CEquipmentList) ProtoMessage() {}

func (x *S2CEquipmentList) ProtoReflect() protoreflect.Message {
	mi := &file_shared_protobuf_proto_equipment_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CEquipmentList.ProtoReflect.Descriptor instead.
func (*S2CEquipmentList) Descriptor() ([]byte, []int) {
	return file_shared_protobuf_proto_equipment_proto_rawDescGZIP(), []int{11}
}

func (x *S2CEquipmentList) GetEquipments() []*VOUserEquipment {
	if x != nil {
		return x.Equipments
	}
	return nil
}

var File_shared_protobuf_proto_equipment_proto protoreflect.FileDescriptor

var file_shared_protobuf_proto_equipment_proto_rawDesc = []byte{
	0x0a, 0x25, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x16, 0x43, 0x32, 0x53, 0x45, 0x71,
	0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x65,
	0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xa1, 0x01, 0x0a, 0x16, 0x53, 0x32, 0x43, 0x45, 0x71, 0x75,
	0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x65, 0x6e,
	0x12, 0x2e, 0x0a, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x56, 0x4f, 0x55, 0x73, 0x65, 0x72, 0x45, 0x71, 0x75, 0x69,
	0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x39, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x56, 0x4f, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6d,
	0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x03, 0x52, 0x09,
	0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x22, 0x43, 0x0a, 0x13, 0x43, 0x32, 0x53,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x76, 0x61, 0x6e, 0x63, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x22, 0x9e,
	0x01, 0x0a, 0x13, 0x53, 0x32, 0x43, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x41,
	0x64, 0x76, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x56, 0x4f, 0x55, 0x73,
	0x65, 0x72, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x71, 0x75,
	0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x56, 0x4f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x61, 0x74, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x22,
	0x28, 0x0a, 0x16, 0x43, 0x32, 0x53, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61, 0x6d, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x67, 0x0a, 0x16, 0x53, 0x32, 0x43,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x63, 0x61, 0x73, 0x74, 0x43,
	0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x61, 0x6d, 0x70, 0x12, 0x39, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x56, 0x4f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x49, 0x0a, 0x1d, 0x43, 0x32, 0x53, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65, 0x63, 0x61, 0x73, 0x74, 0x43,
	0x61, 0x6d, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x22, 0x4f, 0x0a,
	0x1d, 0x53, 0x32, 0x43, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x72, 0x6d, 0x52, 0x65, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61, 0x6d, 0x70, 0x12, 0x2e,
	0x0a, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x56, 0x4f, 0x55, 0x73, 0x65, 0x72, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3a,
	0x0a, 0x10, 0x43, 0x32, 0x53, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x6f,
	0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x73, 0x4c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4c, 0x6f, 0x63, 0x6b, 0x22, 0x3a, 0x0a, 0x10, 0x53, 0x32,
	0x43, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x73, 0x4c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x73, 0x4c, 0x6f, 0x63, 0x6b, 0x22, 0x12, 0x0a, 0x10, 0x43, 0x32, 0x53, 0x45, 0x71, 0x75,
	0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x44, 0x0a, 0x10, 0x53, 0x32,
	0x43, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x30,
	0x0a, 0x0a, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x56, 0x4f, 0x55, 0x73, 0x65, 0x72, 0x45, 0x71, 0x75, 0x69, 0x70,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_shared_protobuf_proto_equipment_proto_rawDescOnce sync.Once
	file_shared_protobuf_proto_equipment_proto_rawDescData = file_shared_protobuf_proto_equipment_proto_rawDesc
)

func file_shared_protobuf_proto_equipment_proto_rawDescGZIP() []byte {
	file_shared_protobuf_proto_equipment_proto_rawDescOnce.Do(func() {
		file_shared_protobuf_proto_equipment_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_protobuf_proto_equipment_proto_rawDescData)
	})
	return file_shared_protobuf_proto_equipment_proto_rawDescData
}

var file_shared_protobuf_proto_equipment_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_shared_protobuf_proto_equipment_proto_goTypes = []interface{}{
	(*C2SEquipmentStrengthen)(nil),        // 0: C2SEquipmentStrengthen
	(*S2CEquipmentStrengthen)(nil),        // 1: S2CEquipmentStrengthen
	(*C2SEquipmentAdvance)(nil),           // 2: C2SEquipmentAdvance
	(*S2CEquipmentAdvance)(nil),           // 3: S2CEquipmentAdvance
	(*C2SEquipmentRecastCamp)(nil),        // 4: C2SEquipmentRecastCamp
	(*S2CEquipmentRecastCamp)(nil),        // 5: S2CEquipmentRecastCamp
	(*C2SEquipmentConfirmRecastCamp)(nil), // 6: C2SEquipmentConfirmRecastCamp
	(*S2CEquipmentConfirmRecastCamp)(nil), // 7: S2CEquipmentConfirmRecastCamp
	(*C2SEquipmentLock)(nil),              // 8: C2SEquipmentLock
	(*S2CEquipmentLock)(nil),              // 9: S2CEquipmentLock
	(*C2SEquipmentList)(nil),              // 10: C2SEquipmentList
	(*S2CEquipmentList)(nil),              // 11: S2CEquipmentList
	(*VOUserEquipment)(nil),               // 12: VOUserEquipment
	(*VOResourceResult)(nil),              // 13: VOResourceResult
}
var file_shared_protobuf_proto_equipment_proto_depIdxs = []int32{
	12, // 0: S2CEquipmentStrengthen.equipment:type_name -> VOUserEquipment
	13, // 1: S2CEquipmentStrengthen.resourceResult:type_name -> VOResourceResult
	12, // 2: S2CEquipmentAdvance.equipment:type_name -> VOUserEquipment
	13, // 3: S2CEquipmentAdvance.resourceResult:type_name -> VOResourceResult
	13, // 4: S2CEquipmentRecastCamp.resourceResult:type_name -> VOResourceResult
	12, // 5: S2CEquipmentConfirmRecastCamp.equipment:type_name -> VOUserEquipment
	12, // 6: S2CEquipmentList.equipments:type_name -> VOUserEquipment
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_shared_protobuf_proto_equipment_proto_init() }
func file_shared_protobuf_proto_equipment_proto_init() {
	if File_shared_protobuf_proto_equipment_proto != nil {
		return
	}
	file_shared_protobuf_proto_vo_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_shared_protobuf_proto_equipment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentStrengthen); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentStrengthen); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentAdvance); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentAdvance); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentRecastCamp); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentRecastCamp); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentConfirmRecastCamp); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentConfirmRecastCamp); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentLock); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentLock); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SEquipmentList); i {
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
		file_shared_protobuf_proto_equipment_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CEquipmentList); i {
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
			RawDescriptor: file_shared_protobuf_proto_equipment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_protobuf_proto_equipment_proto_goTypes,
		DependencyIndexes: file_shared_protobuf_proto_equipment_proto_depIdxs,
		MessageInfos:      file_shared_protobuf_proto_equipment_proto_msgTypes,
	}.Build()
	File_shared_protobuf_proto_equipment_proto = out.File
	file_shared_protobuf_proto_equipment_proto_rawDesc = nil
	file_shared_protobuf_proto_equipment_proto_goTypes = nil
	file_shared_protobuf_proto_equipment_proto_depIdxs = nil
}