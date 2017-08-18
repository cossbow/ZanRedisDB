// Code generated by protoc-gen-gogo.
// source: index_types.proto
// DO NOT EDIT!

/*
	Package rockredis is a generated protocol buffer package.

	It is generated from these files:
		index_types.proto

	It has these top-level messages:
		HsetIndexInfo
		HsetIndexList
*/
package rockredis

import proto "github.com/gogo/protobuf/proto"
import math "math"

// discarding unused import gogoproto "gogoproto/gogo.pb"

import io "io"
import fmt "fmt"
import github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type IndexPropertyDType int32

const (
	Int64V  IndexPropertyDType = 0
	Int32V  IndexPropertyDType = 1
	StringV IndexPropertyDType = 2
)

var IndexPropertyDType_name = map[int32]string{
	0: "Int64V",
	1: "Int32V",
	2: "StringV",
}
var IndexPropertyDType_value = map[string]int32{
	"Int64V":  0,
	"Int32V":  1,
	"StringV": 2,
}

func (x IndexPropertyDType) Enum() *IndexPropertyDType {
	p := new(IndexPropertyDType)
	*p = x
	return p
}
func (x IndexPropertyDType) String() string {
	return proto.EnumName(IndexPropertyDType_name, int32(x))
}
func (x *IndexPropertyDType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(IndexPropertyDType_value, data, "IndexPropertyDType")
	if err != nil {
		return err
	}
	*x = IndexPropertyDType(value)
	return nil
}

type IndexState int32

const (
	InitIndex      IndexState = 0
	BuildingIndex  IndexState = 1
	BuildDoneIndex IndexState = 2
	ReadyIndex     IndexState = 3
	DeletedIndex   IndexState = 4
)

var IndexState_name = map[int32]string{
	0: "InitIndex",
	1: "BuildingIndex",
	2: "BuildDoneIndex",
	3: "ReadyIndex",
	4: "DeletedIndex",
}
var IndexState_value = map[string]int32{
	"InitIndex":      0,
	"BuildingIndex":  1,
	"BuildDoneIndex": 2,
	"ReadyIndex":     3,
	"DeletedIndex":   4,
}

func (x IndexState) Enum() *IndexState {
	p := new(IndexState)
	*p = x
	return p
}
func (x IndexState) String() string {
	return proto.EnumName(IndexState_name, int32(x))
}
func (x *IndexState) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(IndexState_value, data, "IndexState")
	if err != nil {
		return err
	}
	*x = IndexState(value)
	return nil
}

type HsetIndexInfo struct {
	Name             []byte             `protobuf:"bytes,1,opt,name=name" json:"name"`
	IndexField       []byte             `protobuf:"bytes,2,opt,name=index_field" json:"index_field"`
	PrefixLen        int32              `protobuf:"varint,3,opt,name=prefix_len" json:"prefix_len"`
	Unique           int32              `protobuf:"varint,4,opt,name=unique" json:"unique"`
	ValueType        IndexPropertyDType `protobuf:"varint,5,opt,name=value_type,enum=rockredis.IndexPropertyDType" json:"value_type"`
	State            IndexState         `protobuf:"varint,6,opt,name=state,enum=rockredis.IndexState" json:"state"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *HsetIndexInfo) Reset()         { *m = HsetIndexInfo{} }
func (m *HsetIndexInfo) String() string { return proto.CompactTextString(m) }
func (*HsetIndexInfo) ProtoMessage()    {}

type HsetIndexList struct {
	HsetIndexes      []HsetIndexInfo `protobuf:"bytes,1,rep,name=hset_indexes" json:"hset_indexes"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *HsetIndexList) Reset()         { *m = HsetIndexList{} }
func (m *HsetIndexList) String() string { return proto.CompactTextString(m) }
func (*HsetIndexList) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("rockredis.IndexPropertyDType", IndexPropertyDType_name, IndexPropertyDType_value)
	proto.RegisterEnum("rockredis.IndexState", IndexState_name, IndexState_value)
}
func (m *HsetIndexInfo) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = append([]byte{}, data[index:postIndex]...)
			index = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IndexField", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IndexField = append([]byte{}, data[index:postIndex]...)
			index = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrefixLen", wireType)
			}
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				m.PrefixLen |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unique", wireType)
			}
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				m.Unique |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValueType", wireType)
			}
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				m.ValueType |= (IndexPropertyDType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				m.State |= (IndexState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := github_com_gogo_protobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *HsetIndexList) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HsetIndexes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HsetIndexes = append(m.HsetIndexes, HsetIndexInfo{})
			if err := m.HsetIndexes[len(m.HsetIndexes)-1].Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := github_com_gogo_protobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *HsetIndexInfo) Size() (n int) {
	var l int
	_ = l
	if m.Name != nil {
		l = len(m.Name)
		n += 1 + l + sovIndexTypes(uint64(l))
	}
	if m.IndexField != nil {
		l = len(m.IndexField)
		n += 1 + l + sovIndexTypes(uint64(l))
	}
	n += 1 + sovIndexTypes(uint64(m.PrefixLen))
	n += 1 + sovIndexTypes(uint64(m.Unique))
	n += 1 + sovIndexTypes(uint64(m.ValueType))
	n += 1 + sovIndexTypes(uint64(m.State))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *HsetIndexList) Size() (n int) {
	var l int
	_ = l
	if len(m.HsetIndexes) > 0 {
		for _, e := range m.HsetIndexes {
			l = e.Size()
			n += 1 + l + sovIndexTypes(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovIndexTypes(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIndexTypes(x uint64) (n int) {
	return sovIndexTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HsetIndexInfo) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *HsetIndexInfo) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Name != nil {
		data[i] = 0xa
		i++
		i = encodeVarintIndexTypes(data, i, uint64(len(m.Name)))
		i += copy(data[i:], m.Name)
	}
	if m.IndexField != nil {
		data[i] = 0x12
		i++
		i = encodeVarintIndexTypes(data, i, uint64(len(m.IndexField)))
		i += copy(data[i:], m.IndexField)
	}
	data[i] = 0x18
	i++
	i = encodeVarintIndexTypes(data, i, uint64(m.PrefixLen))
	data[i] = 0x20
	i++
	i = encodeVarintIndexTypes(data, i, uint64(m.Unique))
	data[i] = 0x28
	i++
	i = encodeVarintIndexTypes(data, i, uint64(m.ValueType))
	data[i] = 0x30
	i++
	i = encodeVarintIndexTypes(data, i, uint64(m.State))
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *HsetIndexList) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *HsetIndexList) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.HsetIndexes) > 0 {
		for _, msg := range m.HsetIndexes {
			data[i] = 0xa
			i++
			i = encodeVarintIndexTypes(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeFixed64IndexTypes(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32IndexTypes(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintIndexTypes(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
