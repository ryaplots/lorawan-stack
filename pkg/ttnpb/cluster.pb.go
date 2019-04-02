// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/cluster.proto

package ttnpb

import (
	fmt "fmt"
	io "io"
	math "math"
	reflect "reflect"
	strconv "strconv"
	strings "strings"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type PeerInfo_Role int32

const (
	PeerInfo_NONE               PeerInfo_Role = 0
	PeerInfo_ENTITY_REGISTRY    PeerInfo_Role = 1
	PeerInfo_ACCESS             PeerInfo_Role = 2
	PeerInfo_GATEWAY_SERVER     PeerInfo_Role = 3
	PeerInfo_NETWORK_SERVER     PeerInfo_Role = 4
	PeerInfo_APPLICATION_SERVER PeerInfo_Role = 5
	PeerInfo_JOIN_SERVER        PeerInfo_Role = 6
	PeerInfo_CRYPTO_SERVER      PeerInfo_Role = 7
)

var PeerInfo_Role_name = map[int32]string{
	0: "NONE",
	1: "ENTITY_REGISTRY",
	2: "ACCESS",
	3: "GATEWAY_SERVER",
	4: "NETWORK_SERVER",
	5: "APPLICATION_SERVER",
	6: "JOIN_SERVER",
	7: "CRYPTO_SERVER",
}

var PeerInfo_Role_value = map[string]int32{
	"NONE":               0,
	"ENTITY_REGISTRY":    1,
	"ACCESS":             2,
	"GATEWAY_SERVER":     3,
	"NETWORK_SERVER":     4,
	"APPLICATION_SERVER": 5,
	"JOIN_SERVER":        6,
	"CRYPTO_SERVER":      7,
}

func (PeerInfo_Role) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5716c3fcd711eefd, []int{0, 0}
}

// PeerInfo
type PeerInfo struct {
	// Port on which the gRPC server is exposed.
	GRPCPort uint32 `protobuf:"varint,1,opt,name=grpc_port,json=grpcPort,proto3" json:"grpc_port,omitempty"`
	// Indicates whether the gRPC server uses TLS.
	TLS bool `protobuf:"varint,2,opt,name=tls,proto3" json:"tls,omitempty"`
	// Roles of the peer.
	Roles []PeerInfo_Role `protobuf:"varint,3,rep,packed,name=roles,proto3,enum=ttn.lorawan.v3.PeerInfo_Role" json:"roles,omitempty"`
	// Tags of the peer
	Tags                 map[string]string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PeerInfo) Reset()      { *m = PeerInfo{} }
func (*PeerInfo) ProtoMessage() {}
func (*PeerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5716c3fcd711eefd, []int{0}
}
func (m *PeerInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PeerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PeerInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PeerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerInfo.Merge(m, src)
}
func (m *PeerInfo) XXX_Size() int {
	return m.Size()
}
func (m *PeerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PeerInfo proto.InternalMessageInfo

func (m *PeerInfo) GetGRPCPort() uint32 {
	if m != nil {
		return m.GRPCPort
	}
	return 0
}

func (m *PeerInfo) GetTLS() bool {
	if m != nil {
		return m.TLS
	}
	return false
}

func (m *PeerInfo) GetRoles() []PeerInfo_Role {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *PeerInfo) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func init() {
	proto.RegisterEnum("ttn.lorawan.v3.PeerInfo_Role", PeerInfo_Role_name, PeerInfo_Role_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.PeerInfo_Role", PeerInfo_Role_name, PeerInfo_Role_value)
	proto.RegisterType((*PeerInfo)(nil), "ttn.lorawan.v3.PeerInfo")
	golang_proto.RegisterType((*PeerInfo)(nil), "ttn.lorawan.v3.PeerInfo")
	proto.RegisterMapType((map[string]string)(nil), "ttn.lorawan.v3.PeerInfo.TagsEntry")
	golang_proto.RegisterMapType((map[string]string)(nil), "ttn.lorawan.v3.PeerInfo.TagsEntry")
}

func init() { proto.RegisterFile("lorawan-stack/api/cluster.proto", fileDescriptor_5716c3fcd711eefd) }
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/cluster.proto", fileDescriptor_5716c3fcd711eefd)
}

var fileDescriptor_5716c3fcd711eefd = []byte{
	// 556 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x31, 0x4c, 0xdb, 0x4c,
	0x14, 0xc7, 0xef, 0x70, 0x80, 0xe4, 0xf8, 0x80, 0x7c, 0xd7, 0xaa, 0x4a, 0x91, 0xfa, 0x88, 0x98,
	0x52, 0xa9, 0xb1, 0x25, 0x90, 0xda, 0xaa, 0x9d, 0x42, 0x64, 0x21, 0xb7, 0x28, 0xb1, 0x2e, 0x56,
	0x51, 0xba, 0x20, 0x27, 0x75, 0x9c, 0x28, 0xae, 0xcf, 0xb2, 0x2f, 0x41, 0xd9, 0x18, 0x99, 0xaa,
	0x2e, 0x95, 0x3a, 0x56, 0x9d, 0x18, 0x19, 0x19, 0x19, 0x19, 0x19, 0x99, 0x10, 0x3e, 0x2f, 0x8c,
	0x8c, 0x8c, 0x95, 0x9d, 0x06, 0x15, 0x55, 0xdd, 0xfe, 0xef, 0xe7, 0xdf, 0x93, 0xff, 0x27, 0x3d,
	0xb2, 0xee, 0xf1, 0xd0, 0x3e, 0xb0, 0xfd, 0x6a, 0x24, 0xec, 0xee, 0x50, 0xb3, 0x83, 0x81, 0xd6,
	0xf5, 0x46, 0x91, 0x70, 0x42, 0x35, 0x08, 0xb9, 0xe0, 0x74, 0x45, 0x08, 0x5f, 0xfd, 0x2d, 0xa9,
	0xe3, 0xad, 0xb5, 0xaa, 0x3b, 0x10, 0xfd, 0x51, 0x47, 0xed, 0xf2, 0xcf, 0x9a, 0xcb, 0x5d, 0xae,
	0x65, 0x5a, 0x67, 0xd4, 0xcb, 0xa6, 0x6c, 0xc8, 0xd2, 0x74, 0x7d, 0xed, 0xed, 0x1f, 0xba, 0x37,
	0xe9, 0x89, 0xa9, 0xde, 0xad, 0xba, 0x8e, 0x5f, 0x1d, 0xdb, 0xde, 0xe0, 0x93, 0x2d, 0x1c, 0xed,
	0xaf, 0x30, 0x5d, 0xde, 0xf8, 0xa2, 0x90, 0xbc, 0xe9, 0x38, 0xa1, 0xe1, 0xf7, 0x38, 0x7d, 0x4e,
	0x0a, 0x6e, 0x18, 0x74, 0xf7, 0x03, 0x1e, 0x8a, 0x12, 0x2e, 0xe3, 0xca, 0xf2, 0xf6, 0x7f, 0xf2,
	0x6a, 0x3d, 0xbf, 0xc3, 0xcc, 0xba, 0xc9, 0x43, 0xc1, 0xf2, 0xe9, 0xe7, 0x34, 0xd1, 0xa7, 0x44,
	0x11, 0x5e, 0x54, 0x9a, 0x2b, 0xe3, 0x4a, 0x7e, 0x7b, 0x51, 0x5e, 0xad, 0x2b, 0xd6, 0x6e, 0x8b,
	0xa5, 0x8c, 0x6e, 0x91, 0xf9, 0x90, 0x7b, 0x4e, 0x54, 0x52, 0xca, 0x4a, 0x65, 0x65, 0xf3, 0x99,
	0xfa, 0xf0, 0x79, 0xea, 0xec, 0x77, 0x2a, 0xe3, 0x9e, 0xc3, 0xa6, 0x2e, 0x7d, 0x49, 0x72, 0xc2,
	0x76, 0xa3, 0x52, 0xae, 0xac, 0x54, 0x96, 0x36, 0x37, 0xfe, 0xb9, 0x63, 0xd9, 0x6e, 0xa4, 0xfb,
	0x22, 0x9c, 0xb0, 0xcc, 0x5f, 0x7b, 0x45, 0x0a, 0xf7, 0x88, 0x16, 0x89, 0x32, 0x74, 0x26, 0x59,
	0xf3, 0x02, 0x4b, 0x23, 0x7d, 0x4c, 0xe6, 0xc7, 0xb6, 0x37, 0x72, 0xb2, 0xa2, 0x05, 0x36, 0x1d,
	0xde, 0xcc, 0xbd, 0xc6, 0x1b, 0xdf, 0x30, 0xc9, 0xa5, 0x05, 0x68, 0x9e, 0xe4, 0x1a, 0xcd, 0x86,
	0x5e, 0x44, 0xf4, 0x11, 0x59, 0xd5, 0x1b, 0x96, 0x61, 0xb5, 0xf7, 0x99, 0xbe, 0x63, 0xb4, 0x2c,
	0xd6, 0x2e, 0x62, 0x4a, 0xc8, 0x42, 0xad, 0x5e, 0xd7, 0x5b, 0xad, 0xe2, 0x1c, 0xa5, 0x64, 0x65,
	0xa7, 0x66, 0xe9, 0x7b, 0xb5, 0xf6, 0x7e, 0x4b, 0x67, 0x1f, 0x74, 0x56, 0x54, 0x52, 0xd6, 0xd0,
	0xad, 0xbd, 0x26, 0x7b, 0x3f, 0x63, 0x39, 0xfa, 0x84, 0xd0, 0x9a, 0x69, 0xee, 0x1a, 0xf5, 0x9a,
	0x65, 0x34, 0x1b, 0x33, 0x3e, 0x4f, 0x57, 0xc9, 0xd2, 0xbb, 0xa6, 0x71, 0x0f, 0x16, 0xe8, 0xff,
	0x64, 0xb9, 0xce, 0xda, 0xa6, 0xd5, 0x9c, 0xa1, 0xc5, 0xed, 0x9f, 0xf8, 0x3c, 0x06, 0x7c, 0x11,
	0x03, 0xbe, 0x8c, 0x01, 0x5d, 0xc7, 0x80, 0x6e, 0x62, 0x40, 0xb7, 0x31, 0xa0, 0xbb, 0x18, 0xf0,
	0xa1, 0x04, 0x7c, 0x24, 0x01, 0x1d, 0x4b, 0xc0, 0x27, 0x12, 0xd0, 0xa9, 0x04, 0x74, 0x26, 0x01,
	0x9d, 0x4b, 0xc0, 0x17, 0x12, 0xf0, 0xa5, 0x04, 0x74, 0x2d, 0x01, 0xdf, 0x48, 0x40, 0xb7, 0x12,
	0xf0, 0x9d, 0x04, 0x74, 0x98, 0x00, 0x3a, 0x4a, 0x00, 0x7f, 0x4d, 0x00, 0x7d, 0x4f, 0x00, 0xff,
	0x48, 0x00, 0x1d, 0x27, 0x80, 0x4e, 0x12, 0xc0, 0xa7, 0x09, 0xe0, 0xb3, 0x04, 0xf0, 0xc7, 0x17,
	0x2e, 0x57, 0x45, 0xdf, 0x11, 0xfd, 0x81, 0xef, 0x46, 0xaa, 0xef, 0x88, 0x03, 0x1e, 0x0e, 0xb5,
	0x87, 0x97, 0x1b, 0x0c, 0x5d, 0x4d, 0x08, 0x3f, 0xe8, 0x74, 0x16, 0xb2, 0xe3, 0xd9, 0xfa, 0x15,
	0x00, 0x00, 0xff, 0xff, 0x81, 0x7d, 0x44, 0x14, 0xdb, 0x02, 0x00, 0x00,
}

func (x PeerInfo_Role) String() string {
	s, ok := PeerInfo_Role_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *PeerInfo) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PeerInfo)
	if !ok {
		that2, ok := that.(PeerInfo)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.GRPCPort != that1.GRPCPort {
		return false
	}
	if this.TLS != that1.TLS {
		return false
	}
	if len(this.Roles) != len(that1.Roles) {
		return false
	}
	for i := range this.Roles {
		if this.Roles[i] != that1.Roles[i] {
			return false
		}
	}
	if len(this.Tags) != len(that1.Tags) {
		return false
	}
	for i := range this.Tags {
		if this.Tags[i] != that1.Tags[i] {
			return false
		}
	}
	return true
}
func (m *PeerInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PeerInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.GRPCPort != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCluster(dAtA, i, uint64(m.GRPCPort))
	}
	if m.TLS {
		dAtA[i] = 0x10
		i++
		if m.TLS {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.Roles) > 0 {
		dAtA2 := make([]byte, len(m.Roles)*10)
		var j1 int
		for _, num := range m.Roles {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCluster(dAtA, i, uint64(j1))
		i += copy(dAtA[i:], dAtA2[:j1])
	}
	if len(m.Tags) > 0 {
		for k := range m.Tags {
			dAtA[i] = 0x22
			i++
			v := m.Tags[k]
			mapSize := 1 + len(k) + sovCluster(uint64(len(k))) + 1 + len(v) + sovCluster(uint64(len(v)))
			i = encodeVarintCluster(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintCluster(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintCluster(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	return i, nil
}

func encodeVarintCluster(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedPeerInfo(r randyCluster, easy bool) *PeerInfo {
	this := &PeerInfo{}
	this.GRPCPort = r.Uint32()
	this.TLS = bool(r.Intn(2) == 0)
	v1 := r.Intn(10)
	this.Roles = make([]PeerInfo_Role, v1)
	for i := 0; i < v1; i++ {
		this.Roles[i] = PeerInfo_Role([]int32{0, 1, 2, 3, 4, 5, 6, 7}[r.Intn(8)])
	}
	if r.Intn(10) != 0 {
		v2 := r.Intn(10)
		this.Tags = make(map[string]string)
		for i := 0; i < v2; i++ {
			this.Tags[randStringCluster(r)] = randStringCluster(r)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyCluster interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneCluster(r randyCluster) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringCluster(r randyCluster) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneCluster(r)
	}
	return string(tmps)
}
func randUnrecognizedCluster(r randyCluster, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldCluster(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldCluster(dAtA []byte, r randyCluster, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(v4))
	case 1:
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateCluster(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateCluster(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *PeerInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GRPCPort != 0 {
		n += 1 + sovCluster(uint64(m.GRPCPort))
	}
	if m.TLS {
		n += 2
	}
	if len(m.Roles) > 0 {
		l = 0
		for _, e := range m.Roles {
			l += sovCluster(uint64(e))
		}
		n += 1 + sovCluster(uint64(l)) + l
	}
	if len(m.Tags) > 0 {
		for k, v := range m.Tags {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovCluster(uint64(len(k))) + 1 + len(v) + sovCluster(uint64(len(v)))
			n += mapEntrySize + 1 + sovCluster(uint64(mapEntrySize))
		}
	}
	return n
}

func sovCluster(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCluster(x uint64) (n int) {
	return sovCluster((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *PeerInfo) String() string {
	if this == nil {
		return "nil"
	}
	keysForTags := make([]string, 0, len(this.Tags))
	for k := range this.Tags {
		keysForTags = append(keysForTags, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForTags)
	mapStringForTags := "map[string]string{"
	for _, k := range keysForTags {
		mapStringForTags += fmt.Sprintf("%v: %v,", k, this.Tags[k])
	}
	mapStringForTags += "}"
	s := strings.Join([]string{`&PeerInfo{`,
		`GRPCPort:` + fmt.Sprintf("%v", this.GRPCPort) + `,`,
		`TLS:` + fmt.Sprintf("%v", this.TLS) + `,`,
		`Roles:` + fmt.Sprintf("%v", this.Roles) + `,`,
		`Tags:` + mapStringForTags + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringCluster(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *PeerInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCluster
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PeerInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PeerInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GRPCPort", wireType)
			}
			m.GRPCPort = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCluster
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GRPCPort |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TLS", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCluster
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.TLS = bool(v != 0)
		case 3:
			if wireType == 0 {
				var v PeerInfo_Role
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCluster
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= PeerInfo_Role(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Roles = append(m.Roles, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCluster
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthCluster
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthCluster
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				if elementCount != 0 && len(m.Roles) == 0 {
					m.Roles = make([]PeerInfo_Role, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v PeerInfo_Role
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCluster
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= PeerInfo_Role(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Roles = append(m.Roles, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Roles", wireType)
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCluster
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCluster
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCluster
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Tags == nil {
				m.Tags = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCluster
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCluster
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthCluster
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthCluster
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCluster
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthCluster
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthCluster
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCluster(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCluster
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Tags[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCluster(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCluster
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCluster
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCluster(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCluster
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCluster
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCluster
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthCluster
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthCluster
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCluster
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCluster(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthCluster
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCluster = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCluster   = fmt.Errorf("proto: integer overflow")
)