// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goeni/epoch/epoch.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Epoch struct {
	// authority defines the custom module authority. If not set, defaults to the governance module.
	Authority               string    `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	GenesisTime             time.Time `protobuf:"bytes,2,opt,name=genesis_time,json=genesisTime,proto3,stdtime" json:"genesis_time" yaml:"genesis_time"`
	EpochInterval           uint64    `protobuf:"varint,3,opt,name=epoch_interval,json=epochInterval,proto3" json:"epoch_interval" yaml:"epoch_interval"`
	CurrentEpoch            uint64    `protobuf:"varint,4,opt,name=current_epoch,json=currentEpoch,proto3" json:"current_epoch" yaml:"current_epoch"`
	CurrentEpochStartHeight uint64    `protobuf:"varint,5,opt,name=current_epoch_start_height,json=currentEpochStartHeight,proto3" json:"current_epoch_start_height" yaml:"current_epoch_start_height"`
	CurrentEpochHeight      int64     `protobuf:"varint,6,opt,name=current_epoch_height,json=currentEpochHeight,proto3" json:"current_epoch_height" yaml:"current_epoch_height"`
}

func (m *Epoch) Reset()         { *m = Epoch{} }
func (m *Epoch) String() string { return proto.CompactTextString(m) }
func (*Epoch) ProtoMessage()    {}
func (*Epoch) Descriptor() ([]byte, []int) {
	return fileDescriptor_88b66597550cc9b8, []int{0}
}
func (m *Epoch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Epoch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Epoch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Epoch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Epoch.Merge(m, src)
}
func (m *Epoch) XXX_Size() int {
	return m.Size()
}
func (m *Epoch) XXX_DiscardUnknown() {
	xxx_messageInfo_Epoch.DiscardUnknown(m)
}

var xxx_messageInfo_Epoch proto.InternalMessageInfo

func (m *Epoch) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *Epoch) GetGenesisTime() time.Time {
	if m != nil {
		return m.GenesisTime
	}
	return time.Time{}
}

func (m *Epoch) GetEpochInterval() uint64 {
	if m != nil {
		return m.EpochInterval
	}
	return 0
}

func (m *Epoch) GetCurrentEpoch() uint64 {
	if m != nil {
		return m.CurrentEpoch
	}
	return 0
}

func (m *Epoch) GetCurrentEpochStartHeight() uint64 {
	if m != nil {
		return m.CurrentEpochStartHeight
	}
	return 0
}

func (m *Epoch) GetCurrentEpochHeight() int64 {
	if m != nil {
		return m.CurrentEpochHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Epoch)(nil), "goeni.epoch.Epoch")
}

func init() { proto.RegisterFile("goeni/epoch/epoch.proto", fileDescriptor_88b66597550cc9b8) }

var fileDescriptor_88b66597550cc9b8 = []byte{
	// 422 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x31, 0x6f, 0xd4, 0x30,
	0x1c, 0xc5, 0xcf, 0xf4, 0x5a, 0xa9, 0x6e, 0xcb, 0x60, 0x0e, 0x11, 0x05, 0x64, 0x87, 0x4c, 0xa9,
	0x50, 0x13, 0x09, 0x06, 0x24, 0xc6, 0x54, 0x48, 0xb0, 0x30, 0x04, 0x26, 0x06, 0xa2, 0x5c, 0x30,
	0x8e, 0xa5, 0x4b, 0x1c, 0x25, 0x0e, 0xe2, 0x36, 0x3e, 0x42, 0x37, 0xbe, 0x52, 0xc7, 0x8e, 0x4c,
	0x06, 0xdd, 0x6d, 0x19, 0xf3, 0x09, 0x50, 0xec, 0x9c, 0x48, 0x4e, 0xa7, 0x2e, 0xa7, 0xfb, 0xbf,
	0xf7, 0xfe, 0xef, 0x27, 0x47, 0x7f, 0x88, 0x69, 0xc1, 0xd3, 0x2c, 0xe1, 0x85, 0xa4, 0xb5, 0x0c,
	0x68, 0x29, 0xd2, 0xcc, 0xfc, 0xfa, 0x65, 0x25, 0xa4, 0x40, 0x68, 0xec, 0xfb, 0xda, 0xb1, 0x17,
	0x4c, 0x30, 0xa1, 0xed, 0xa0, 0xff, 0x67, 0x92, 0x36, 0x61, 0x42, 0xb0, 0x15, 0x0d, 0xf4, 0xb4,
	0x6c, 0xbe, 0x05, 0x92, 0xe7, 0xb4, 0x96, 0x49, 0x5e, 0x0e, 0x01, 0xbc, 0x1f, 0xf8, 0xda, 0x54,
	0x89, 0xe4, 0xa2, 0x30, 0xbe, 0xfb, 0x6b, 0x0e, 0x8f, 0xdf, 0xf6, 0x00, 0xf4, 0x0c, 0x9e, 0x26,
	0x8d, 0xcc, 0x44, 0xc5, 0xe5, 0xda, 0x02, 0x0e, 0xf0, 0x4e, 0xa3, 0xff, 0x02, 0xfa, 0x02, 0xcf,
	0x19, 0x2d, 0x68, 0xcd, 0xeb, 0xb8, 0x47, 0x58, 0x0f, 0x1c, 0xe0, 0x9d, 0xbd, 0xb4, 0x7d, 0x53,
	0xef, 0xef, 0xea, 0xfd, 0x4f, 0x3b, 0x7e, 0x48, 0x6e, 0x15, 0x99, 0x75, 0x8a, 0x3c, 0x5a, 0x27,
	0xf9, 0xea, 0x8d, 0x3b, 0xde, 0x76, 0x6f, 0xfe, 0x10, 0x10, 0x9d, 0x0d, 0x52, 0xbf, 0x82, 0x22,
	0xf8, 0x50, 0xbf, 0x33, 0xee, 0x1f, 0x5d, 0x7d, 0x4f, 0x56, 0xd6, 0x91, 0x03, 0xbc, 0x79, 0xf8,
	0xa2, 0x55, 0x64, 0xcf, 0xe9, 0x14, 0x79, 0x6c, 0x3a, 0xa7, 0xba, 0x1b, 0x5d, 0x68, 0xe1, 0xfd,
	0x30, 0xa3, 0x0f, 0xf0, 0x22, 0x6d, 0xaa, 0x8a, 0x16, 0x32, 0xd6, 0x86, 0x35, 0xd7, 0x95, 0x97,
	0xad, 0x22, 0x53, 0xa3, 0x53, 0x64, 0x61, 0x1a, 0x27, 0xb2, 0x1b, 0x9d, 0x0f, 0xb3, 0xf9, 0x42,
	0x3f, 0x01, 0xb4, 0x27, 0x81, 0xb8, 0x96, 0x49, 0x25, 0xe3, 0x8c, 0x72, 0x96, 0x49, 0xeb, 0x58,
	0xb7, 0x5f, 0xb7, 0x8a, 0xdc, 0x93, 0xea, 0x14, 0x79, 0x7e, 0x00, 0x35, 0xc9, 0xb8, 0xd1, 0x93,
	0x31, 0xf7, 0x63, 0x6f, 0xbd, 0xd3, 0x0e, 0xe2, 0x70, 0x31, 0xdd, 0x1b, 0xd8, 0x27, 0x0e, 0xf0,
	0x8e, 0xc2, 0xd7, 0xad, 0x22, 0x07, 0xfd, 0x4e, 0x91, 0xa7, 0x87, 0xa8, 0x3b, 0x1e, 0x1a, 0xf3,
	0x0c, 0x2a, 0xbc, 0xbe, 0xdd, 0x60, 0x70, 0xb7, 0xc1, 0xe0, 0xef, 0x06, 0x83, 0x9b, 0x2d, 0x9e,
	0xdd, 0x6d, 0xf1, 0xec, 0xf7, 0x16, 0xcf, 0x3e, 0x5f, 0x32, 0x2e, 0xb3, 0x66, 0xe9, 0xa7, 0x22,
	0x0f, 0x68, 0xc1, 0xaf, 0xf4, 0xa9, 0x06, 0x4c, 0x5c, 0xd1, 0x82, 0x07, 0x3f, 0x86, 0x7b, 0x96,
	0xeb, 0x92, 0xd6, 0xcb, 0x13, 0x7d, 0x18, 0xaf, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x6b, 0x15,
	0x27, 0x99, 0xf2, 0x02, 0x00, 0x00,
}

func (m *Epoch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Epoch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Epoch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CurrentEpochHeight != 0 {
		i = encodeVarintEpoch(dAtA, i, uint64(m.CurrentEpochHeight))
		i--
		dAtA[i] = 0x30
	}
	if m.CurrentEpochStartHeight != 0 {
		i = encodeVarintEpoch(dAtA, i, uint64(m.CurrentEpochStartHeight))
		i--
		dAtA[i] = 0x28
	}
	if m.CurrentEpoch != 0 {
		i = encodeVarintEpoch(dAtA, i, uint64(m.CurrentEpoch))
		i--
		dAtA[i] = 0x20
	}
	if m.EpochInterval != 0 {
		i = encodeVarintEpoch(dAtA, i, uint64(m.EpochInterval))
		i--
		dAtA[i] = 0x18
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.GenesisTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.GenesisTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintEpoch(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEpoch(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEpoch(dAtA []byte, offset int, v uint64) int {
	offset -= sovEpoch(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Epoch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEpoch(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.GenesisTime)
	n += 1 + l + sovEpoch(uint64(l))
	if m.EpochInterval != 0 {
		n += 1 + sovEpoch(uint64(m.EpochInterval))
	}
	if m.CurrentEpoch != 0 {
		n += 1 + sovEpoch(uint64(m.CurrentEpoch))
	}
	if m.CurrentEpochStartHeight != 0 {
		n += 1 + sovEpoch(uint64(m.CurrentEpochStartHeight))
	}
	if m.CurrentEpochHeight != 0 {
		n += 1 + sovEpoch(uint64(m.CurrentEpochHeight))
	}
	return n
}

func sovEpoch(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEpoch(x uint64) (n int) {
	return sovEpoch(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Epoch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEpoch
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
			return fmt.Errorf("proto: Epoch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Epoch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEpoch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEpoch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
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
				return ErrInvalidLengthEpoch
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEpoch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.GenesisTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochInterval", wireType)
			}
			m.EpochInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochInterval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentEpoch", wireType)
			}
			m.CurrentEpoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentEpoch |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentEpochStartHeight", wireType)
			}
			m.CurrentEpochStartHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentEpochStartHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentEpochHeight", wireType)
			}
			m.CurrentEpochHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentEpochHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEpoch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEpoch
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
func skipEpoch(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEpoch
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
					return 0, ErrIntOverflowEpoch
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEpoch
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
				return 0, ErrInvalidLengthEpoch
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEpoch
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEpoch
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEpoch        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEpoch          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEpoch = fmt.Errorf("proto: unexpected end of group")
)
