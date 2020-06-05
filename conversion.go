// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
)

const (
	maxWords = 256 / bits.UintSize // number of big.Words in 256-bit

	// The constants below work as compile-time checks: in case evaluated to
	// negative value it cannot be assigned to uint type and compilation fails.
	// These particular expressions check if maxWords either 4 or 8 matching
	// 32-bit and 64-bit architectures.
	_ uint = -(maxWords & (maxWords - 1)) // maxWords is power of two.
	_ uint = -(maxWords & ^(4 | 8))       // maxWords is 4 or 8.
)

// ToBig returns a big.Int version of z.
func (z *Int) ToBig() *big.Int {
	b := new(big.Int)
	switch maxWords { // Compile-time check.
	case 4: // 64-bit architectures.
		words := [4]big.Word{big.Word(z[0]), big.Word(z[1]), big.Word(z[2]), big.Word(z[3])}
		b.SetBits(words[:])
	case 8: // 32-bit architectures.
		words := [8]big.Word{
			big.Word(z[0]), big.Word(z[0] >> 32),
			big.Word(z[1]), big.Word(z[1] >> 32),
			big.Word(z[2]), big.Word(z[2] >> 32),
			big.Word(z[3]), big.Word(z[3] >> 32),
		}
		b.SetBits(words[:])
	}
	return b
}

// FromBig is a convenience-constructor from big.Int.
// Returns a new Int and whether overflow occurred.
func FromBig(b *big.Int) (*Int, bool) {
	z := &Int{}
	overflow := z.SetFromBig(b)
	return z, overflow
}

// SetFromBig converts a big.Int to Int and sets the value to z.
// TODO: Ensure we have sufficient testing, esp for negative bigints.
func (z *Int) SetFromBig(b *big.Int) bool {
	z.Clear()
	words := b.Bits()
	overflow := len(words) > maxWords

	switch maxWords { // Compile-time check.
	case 4: // 64-bit architectures.
		if len(words) > 0 {
			z[0] = uint64(words[0])
			if len(words) > 1 {
				z[1] = uint64(words[1])
				if len(words) > 2 {
					z[2] = uint64(words[2])
					if len(words) > 3 {
						z[3] = uint64(words[3])
					}
				}
			}
		}
	case 8: // 32-bit architectures.
		numWords := len(words)
		if overflow {
			numWords = maxWords
		}
		for i := 0; i < numWords; i++ {
			if i%2 == 0 {
				z[i/2] = uint64(words[i])
			} else {
				z[i/2] |= uint64(words[i]) << 32
			}
		}
	}

	if b.Sign() == -1 {
		z.Neg(z)
	}
	return overflow
}

// Format implements fmt.Formatter. It accepts the formats
// 'b' (binary), 'o' (octal with 0 prefix), 'O' (octal with 0o prefix),
// 'd' (decimal), 'x' (lowercase hexadecimal), and
// 'X' (uppercase hexadecimal).
// Also supported are the full suite of package fmt's format
// flags for integral types, including '+' and ' ' for sign
// control, '#' for leading zero in octal and for hexadecimal,
// a leading "0x" or "0X" for "%#x" and "%#X" respectively,
// specification of minimum digits precision, output field
// width, space or zero padding, and '-' for left or right
// justification.
//
func (z *Int) Format(s fmt.State, ch rune) {
	z.ToBig().Format(s, ch)
}

// SetBytes8 is identical to SetBytes(in[:8]), but panics is input is too short
func (z *Int) SetBytes8(in []byte) *Int {
	_ = in[7] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = binary.BigEndian.Uint64(in[0:8])
	return z
}

// SetBytes16 is identical to SetBytes(in[:16]), but panics is input is too short
func (z *Int) SetBytes16(in []byte) *Int {
	_ = in[15] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = binary.BigEndian.Uint64(in[0:8])
	z[0] = binary.BigEndian.Uint64(in[8:16])
	return z
}

// SetBytes16 is identical to SetBytes(in[:24]), but panics is input is too short
func (z *Int) SetBytes24(in []byte) *Int {
	_ = in[23] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = binary.BigEndian.Uint64(in[0:8])
	z[1] = binary.BigEndian.Uint64(in[8:16])
	z[0] = binary.BigEndian.Uint64(in[16:24])
	return z
}

func (z *Int) SetBytes32(in []byte) *Int {
	_ = in[31] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = binary.BigEndian.Uint64(in[0:8])
	z[2] = binary.BigEndian.Uint64(in[8:16])
	z[1] = binary.BigEndian.Uint64(in[16:24])
	z[0] = binary.BigEndian.Uint64(in[24:32])
	return z
}

func (z *Int) SetBytes1(in []byte) *Int {
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = uint64(in[0])
	return z
}

func (z *Int) SetBytes9(in []byte) *Int {
	_ = in[8] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = uint64(in[0])
	z[0] = binary.BigEndian.Uint64(in[1:9])
	return z
}

func (z *Int) SetBytes17(in []byte) *Int {
	_ = in[16] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = uint64(in[0])
	z[1] = binary.BigEndian.Uint64(in[1:9])
	z[0] = binary.BigEndian.Uint64(in[9:17])
	return z
}

func (z *Int) SetBytes25(in []byte) *Int {
	_ = in[24] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = uint64(in[0])
	z[2] = binary.BigEndian.Uint64(in[1:9])
	z[1] = binary.BigEndian.Uint64(in[9:17])
	z[0] = binary.BigEndian.Uint64(in[17:25])
	return z
}

func (z *Int) SetBytes2(in []byte) *Int {
	_ = in[1] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = uint64(binary.BigEndian.Uint16(in[0:2]))
	return z
}

func (z *Int) SetBytes10(in []byte) *Int {
	_ = in[9] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = uint64(binary.BigEndian.Uint16(in[0:2]))
	z[0] = binary.BigEndian.Uint64(in[2:10])
	return z
}

func (z *Int) SetBytes18(in []byte) *Int {
	_ = in[17] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = uint64(binary.BigEndian.Uint16(in[0:2]))
	z[1] = binary.BigEndian.Uint64(in[2:10])
	z[0] = binary.BigEndian.Uint64(in[10:18])
	return z
}

func (z *Int) SetBytes26(in []byte) *Int {
	_ = in[25] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = uint64(binary.BigEndian.Uint16(in[0:2]))
	z[2] = binary.BigEndian.Uint64(in[2:10])
	z[1] = binary.BigEndian.Uint64(in[10:18])
	z[0] = binary.BigEndian.Uint64(in[18:26])
	return z
}

func (z *Int) SetBytes3(in []byte) *Int {
	_ = in[2] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = uint64(binary.BigEndian.Uint16(in[1:3])) | uint64(in[0])<<16
	return z
}

func (z *Int) SetBytes11(in []byte) *Int {
	_ = in[10] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = uint64(binary.BigEndian.Uint16(in[1:3])) | uint64(in[0])<<16
	z[0] = binary.BigEndian.Uint64(in[3:11])
	return z
}

func (z *Int) SetBytes19(in []byte) *Int {
	_ = in[18] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = uint64(binary.BigEndian.Uint16(in[1:3])) | uint64(in[0])<<16
	z[1] = binary.BigEndian.Uint64(in[3:11])
	z[0] = binary.BigEndian.Uint64(in[11:19])
	return z
}

func (z *Int) SetBytes27(in []byte) *Int {
	_ = in[26] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = uint64(binary.BigEndian.Uint16(in[1:3])) | uint64(in[0])<<16
	z[2] = binary.BigEndian.Uint64(in[3:11])
	z[1] = binary.BigEndian.Uint64(in[11:19])
	z[0] = binary.BigEndian.Uint64(in[19:27])
	return z
}

func (z *Int) SetBytes4(in []byte) *Int {
	_ = in[3] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = uint64(binary.BigEndian.Uint32(in[0:4]))
	return z
}

func (z *Int) SetBytes12(in []byte) *Int {
	_ = in[11] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = uint64(binary.BigEndian.Uint32(in[0:4]))
	z[0] = binary.BigEndian.Uint64(in[4:12])
	return z
}

func (z *Int) SetBytes20(in []byte) *Int {
	_ = in[19] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = uint64(binary.BigEndian.Uint32(in[0:4]))
	z[1] = binary.BigEndian.Uint64(in[4:12])
	z[0] = binary.BigEndian.Uint64(in[12:20])
	return z
}

func (z *Int) SetBytes28(in []byte) *Int {
	_ = in[27] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = uint64(binary.BigEndian.Uint32(in[0:4]))
	z[2] = binary.BigEndian.Uint64(in[4:12])
	z[1] = binary.BigEndian.Uint64(in[12:20])
	z[0] = binary.BigEndian.Uint64(in[20:28])
	return z
}

func (z *Int) SetBytes5(in []byte) *Int {
	_ = in[4] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = bigEndianUint40(in[0:5])
	return z
}

func (z *Int) SetBytes13(in []byte) *Int {
	_ = in[12] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = bigEndianUint40(in[0:5])
	z[0] = binary.BigEndian.Uint64(in[5:13])
	return z
}

func (z *Int) SetBytes21(in []byte) *Int {
	_ = in[20] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = bigEndianUint40(in[0:5])
	z[1] = binary.BigEndian.Uint64(in[5:13])
	z[0] = binary.BigEndian.Uint64(in[13:21])
	return z
}

func (z *Int) SetBytes29(in []byte) *Int {
	_ = in[23] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = bigEndianUint40(in[0:5])
	z[2] = binary.BigEndian.Uint64(in[5:13])
	z[1] = binary.BigEndian.Uint64(in[13:21])
	z[0] = binary.BigEndian.Uint64(in[21:29])
	return z
}

func (z *Int) SetBytes6(in []byte) *Int {
	_ = in[5] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = bigEndianUint48(in[0:6])
	return z
}

func (z *Int) SetBytes14(in []byte) *Int {
	_ = in[13] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = bigEndianUint48(in[0:6])
	z[0] = binary.BigEndian.Uint64(in[6:14])
	return z
}

func (z *Int) SetBytes22(in []byte) *Int {
	_ = in[21] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = bigEndianUint48(in[0:6])
	z[1] = binary.BigEndian.Uint64(in[6:14])
	z[0] = binary.BigEndian.Uint64(in[14:22])
	return z
}

func (z *Int) SetBytes30(in []byte) *Int {
	_ = in[29] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = bigEndianUint48(in[0:6])
	z[2] = binary.BigEndian.Uint64(in[6:14])
	z[1] = binary.BigEndian.Uint64(in[14:22])
	z[0] = binary.BigEndian.Uint64(in[22:30])
	return z
}

func (z *Int) SetBytes7(in []byte) *Int {
	_ = in[6] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2], z[1] = 0, 0, 0
	z[0] = bigEndianUint56(in[0:7])
	return z
}

func (z *Int) SetBytes15(in []byte) *Int {
	_ = in[14] // bounds check hint to compiler; see golang.org/issue/14808
	z[3], z[2] = 0, 0
	z[1] = bigEndianUint56(in[0:7])
	z[0] = binary.BigEndian.Uint64(in[7:15])
	return z
}

func (z *Int) SetBytes23(in []byte) *Int {
	_ = in[22] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = 0
	z[2] = bigEndianUint56(in[0:7])
	z[1] = binary.BigEndian.Uint64(in[7:15])
	z[0] = binary.BigEndian.Uint64(in[15:23])
	return z
}

func (z *Int) SetBytes31(in []byte) *Int {
	_ = in[30] // bounds check hint to compiler; see golang.org/issue/14808
	z[3] = bigEndianUint56(in[0:7])
	z[2] = binary.BigEndian.Uint64(in[7:15])
	z[1] = binary.BigEndian.Uint64(in[15:23])
	z[0] = binary.BigEndian.Uint64(in[23:31])
	return z
}

// Utility methods that are "missing" among the bigEndian.UintXX methods.

func bigEndianUint40(b []byte) uint64 {
	_ = b[4] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 |
		uint64(b[0])<<32
}

func bigEndianUint48(b []byte) uint64 {
	_ = b[5] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[5]) | uint64(b[4])<<8 | uint64(b[3])<<16 | uint64(b[2])<<24 |
		uint64(b[1])<<32 | uint64(b[0])<<40
}

func bigEndianUint56(b []byte) uint64 {
	_ = b[6] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[6]) | uint64(b[5])<<8 | uint64(b[4])<<16 | uint64(b[3])<<24 |
		uint64(b[2])<<32 | uint64(b[1])<<40 | uint64(b[0])<<48
}
