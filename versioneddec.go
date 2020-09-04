// Implementation of the versioned decoder.

package s2mdec

import (
	"fmt"
	"strconv"

	"github.com/icza/s2prot"
)

// VersionedDec is a versioned decoder.
type VersionedDec struct {
	*BitPackedBuff // Data source: bit-packed buffer
}

// NewVersionedDec creates a new bit-packed decoder.
func NewVersionedDec(contents []byte) *VersionedDec {
	return &VersionedDec{
		BitPackedBuff: &BitPackedBuff{
			contents:  contents,
			bigEndian: true, // All versioned decoder uses big endian order
		},
	}
}

// DataType of a field of struct.
type DataType int

// DataType consts.
const (
	DataTypeArray    DataType = 0x00 // List of elements of the same type
	DataTypeBitArray DataType = 0x01 // List of bits (packed into a byte array)
	DataTypeBlob     DataType = 0x02 // A byte array
	DataTypeChoice   DataType = 0x03 // A choice of multiple types (one of multiple)
	DataTypeOptional DataType = 0x04 // Optionally a value (of a specified type)
	DataTypeStruct   DataType = 0x05 // A structure (list of fields)
	DataTypeUint8    DataType = 0x06 // Bool // A bool value
	DataTypeUint32   DataType = 0x07 // FourCC // 4 bytes data, usually interpreted as string
	DataTypeUint64   DataType = 0x08
	DataTypeVarInt   DataType = 0x09 // Int // An integer number
)

// ReadStruct decodes a value specified by dataType and returns the decoded value.
// ReadStruct reads a nested data structure. If the type is not specified the first byte is used as the type identifier.
func (d *VersionedDec) ReadStruct(dataTypes ...DataType) interface{} {
	if len(dataTypes) < 1 {
		dataTypes = []DataType{DataType(d.ReadBits8())}
	}

	switch dataTypes[0] {
	case DataTypeArray:
		length := d.ReadVarInt()
		arr := make([]interface{}, length)
		for i := range arr {
			arr[i] = d.ReadStruct()
		}
		return arr
	case DataTypeBitArray:
		length := int(d.ReadVarInt())
		barr := s2prot.BitArr{Count: length, Data: d.ReadAligned((length + 7) / 8)}
		return barr
	case DataTypeBlob:
		length := int(d.ReadVarInt())
		return string(d.ReadAligned(length))
	case DataTypeChoice:
		_ = int(d.ReadVarInt()) // flag
		return d.ReadStruct()
	case DataTypeOptional:
		if exists := d.ReadBits8(); exists != 0 {
			return d.ReadStruct()
		}
		return nil
	case DataTypeStruct:
		// TODO order should be preserved! Map does not preserve it!
		s := s2prot.Struct{}
		nEntries := int(d.ReadVarInt()) // length
		for i := 0; i < nEntries; i++ {
			s[strconv.Itoa(int(d.ReadVarInt()))] = d.ReadStruct()
		}
		return s
	case DataTypeUint8:
		return int64(d.ReadBits8()) // This is usually bool and is put int64 to be the same type as VarInt.
	case DataTypeUint32:
		return string(d.ReadAligned(4))
	case DataTypeUint64:
		return string(d.ReadAligned(8))
	case DataTypeVarInt:
		return d.ReadVarInt()
	default:
		panic(fmt.Errorf("unknown data type: %d", dataTypes[0]))
	}

	// return nil
}

// ReadVarInt reads a variable-length int value.
// Format: read from input by 8 bits. Highest bit tells if have to read more bytes,
// lowest bit of the first byte (first 8 bits) is not data but tells if the number is negative.
func (b *BitPackedBuff) ReadVarInt() int64 {
	var data, value int64
	for shift := uint(0); ; shift += 7 {
		data = int64(b.ReadBits8())
		value |= (data & 0x7f) << shift
		if (data & 0x80) == 0 {
			if value&0x01 > 0 {
				return -(value >> 1)
			}
			return value >> 1
		}
	}
}

// SkipInstance reads and discards an instance whose type is deducted from the read Field type.
func (b *BitPackedBuff) SkipInstance() {
	fieldType := b.ReadBits8()
	switch fieldType {
	case 0: // array
		for i := b.ReadVarInt(); i > 0; i-- {
			b.SkipInstance()
		}
	case 1: // bit array
		b.ReadAligned((int(b.ReadVarInt()) + 7) / 8)
	case 2: // blob
		b.ReadAligned(int(b.ReadVarInt()))
	case 3: // choice
		b.ReadVarInt() // tag
		b.SkipInstance()
	case 4: // optional
		if b.ReadBits8() != 0 {
			b.SkipInstance()
		}
	case 5: // struct
		for i := b.ReadVarInt(); i > 0; i-- {
			b.ReadVarInt() // tag
			b.SkipInstance()
		}
	case 6: // uint8
		b.ReadBits8() // b.ReadAligned(1)
	case 7: // uint32
		b.ReadAligned(4)
	case 8: // uint64
		b.ReadAligned(8)
	case 9: // vint
		b.ReadVarInt()
	}
}
