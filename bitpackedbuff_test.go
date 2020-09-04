// Edited from github.com/icza/s2prot

package s2mdec

import (
	"bytes"
	"testing"
)

func TestEOFD(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{}, bigEndian: true}
	if !bb.EOF() {
		t.Error("EOF falsely NOT reported.")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3}, bigEndian: true}

	if bb.EOF() {
		t.Error("EOF falsely reported.")
	}
	bb.ReadBits(1)
	if bb.EOF() {
		t.Error("EOF falsely reported.")
	}
	bb.ReadBits(7)
	if bb.EOF() {
		t.Error("EOF falsely reported.")
	}
	bb.ReadBits(1)
	if bb.EOF() {
		t.Error("EOF falsely reported.")
	}
	bb.ReadBits(12)
	if bb.EOF() {
		t.Error("EOF falsely reported.")
	}
	bb.ReadBits(3)
	if !bb.EOF() {
		t.Error("EOF falsely NOT reported.")
	}
}

func TestByteAlign(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3}, bigEndian: true}

	bb.ByteAlign()
	if bb.ReadBits(8) != 1 {
		t.Error("Unexpected value!")
	}

	bb.ReadBits(1)
	bb.ByteAlign()
	if bb.ReadBits(8) != 3 {
		t.Error("Unexpected value!")
	}
}

func TestReadBits1(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{0xaa, 0xaa}, bigEndian: true}

	for expected := false; !bb.EOF(); expected = !expected {
		if bb.ReadBits1() != expected {
			t.Error("Unexpected value!")
		}
	}
}

func TestReadBits8(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3, 4}, bigEndian: true}

	if bb.ReadBits8() != 1 {
		t.Error("Unexpected value!")
	}
	bb.ReadBits(3)
	if bb.ReadBits8() != 3 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits8() != 4 {
		t.Error("Unexpected value!")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3, 4}, bigEndian: false}

	if bb.ReadBits8() != 1 {
		t.Error("Unexpected value!")
	}
	bb.ReadBits(3)
	if bb.ReadBits8() != 0x60 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits8() != 0x80 {
		t.Error("Unexpected value!")
	}
}

func TestReadBits(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3, 4}, bigEndian: true}
	if bb.ReadBits(0) != 0 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(8) != 1 {
		t.Error("Unexpected value!")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3, 4}, bigEndian: true}

	if bb.ReadBits(3) != 1 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(13) != 2 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(1) != 1 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(15) != 0x0104 {
		t.Error("Unexpected value!")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3, 4}, bigEndian: false}

	if bb.ReadBits(3) != 1 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(13) != 0x40 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(1) != 1 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(15) != 0x0201 {
		t.Error("Unexpected value!")
	}
}

func TestReadBitsSpecial(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3, 4, 5, 6, 7, 8}, bigEndian: true}

	bb.ReadBits(3) // Non-empty cache
	if bb.ReadBits(8) != 2 {
		t.Error("Unexpected value!")
	}
	// 111111111 222222222 333333333 444444444 555555555 666666666 777777777 888888888
	// 0000 0001 0000 0010 0000 0011 0000 0100 0000 0101 0000 0110 0000 0111 0000 1000
	// 0000 0000 0001 1100
	if bb.ReadBits(16) != 0x001c {
		t.Error("Unexpected value!")
	}
	// 0000 0000 0010 1000 0011 0000 0011 1000
	if bb.ReadBits(32) != 0x00283038 {
		t.Error("Unexpected value!")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3, 4, 5, 6, 7}, bigEndian: false}

	// Empty cache
	if bb.ReadBits(8) != 1 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(16) != 0x0302 {
		t.Error("Unexpected value!")
	}
	if bb.ReadBits(32) != 0x07060504 {
		t.Error("Unexpected value!")
	}
}

func TestReadAligned(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3, 4, 5, 6, 7, 8}, bigEndian: true}

	if !bytes.Equal([]byte{}, bb.ReadAligned(0)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{1}, bb.ReadAligned(1)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{2, 3}, bb.ReadAligned(2)) {
		t.Error("Unexpected value!")
	}
	bb.ReadBits(3)
	if !bytes.Equal([]byte{5, 6, 7, 8}, bb.ReadAligned(4)) {
		t.Error("Unexpected value!")
	}
}

func TestReadUnaligned(t *testing.T) {
	bb := &BitPackedBuff{contents: []byte{1, 2, 3, 4, 5, 6, 7, 8}, bigEndian: true}

	if !bytes.Equal([]byte{}, bb.ReadUnaligned(0)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{1, 2}, bb.ReadUnaligned(2)) {
		t.Error("Unexpected value!")
	}
	bb.ReadBits(3)
	if !bytes.Equal([]byte{0x04, 0x05}, bb.ReadUnaligned(2)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{0x06, 0x07, 0x00}, bb.ReadUnaligned(3)) {
		t.Error("Unexpected value!")
	}

	bb = &BitPackedBuff{contents: []byte{1, 2, 3, 4, 5, 6, 7, 8}, bigEndian: false}

	if !bytes.Equal([]byte{}, bb.ReadUnaligned(0)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{1, 2}, bb.ReadUnaligned(2)) {
		t.Error("Unexpected value!")
	}
	bb.ReadBits(3)
	// 111111111 222222222 333333333 444444444 555555555 666666666 777777777 888888888
	// 0000 0001 0000 0010 0000 0011 0000 0100 0000 0101 0000 0110 0000 0111 0000 1000
	if !bytes.Equal([]byte{0x80, 0xa0}, bb.ReadUnaligned(2)) {
		t.Error("Unexpected value!")
	}
	if !bytes.Equal([]byte{0xc0, 0xe0, 0x00}, bb.ReadUnaligned(3)) {
		t.Error("Unexpected value!")
	}
}
