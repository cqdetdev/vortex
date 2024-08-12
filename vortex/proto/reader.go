package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"io"
	"math"
	"unsafe"

	"github.com/google/uuid"
)


type Reader struct {
	r interface {
		io.Reader
		io.ByteReader
	}
	shieldID      int32
	limitsEnabled bool
}

// NewReader creates a new Reader using the io.ByteReader passed as underlying source to read bytes from.
func NewReader(r interface {
	io.Reader
	io.ByteReader
}, shieldID int32, enableLimits bool) *Reader {
	return &Reader{r: r, shieldID: shieldID, limitsEnabled: enableLimits}
}

// Uint8 reads a uint8 from the underlying buffer.
func (r *Reader) Uint8(x *uint8) {
	var err error
	*x, err = r.r.ReadByte()
	if err != nil {
		r.panic(err)
	}
}

// Int8 reads an int8 from the underlying buffer.
func (r *Reader) Int8(x *int8) {
	var b uint8
	r.Uint8(&b)
	*x = int8(b)
}

// Uint16 reads a little endian uint16 from the underlying buffer.
func (r *Reader) Uint16(x *uint16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = binary.BigEndian.Uint16(b)
}

// Int16 reads a little endian int16 from the underlying buffer.
func (r *Reader) Int16(x *int16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = int16(binary.BigEndian.Uint16(b))
}

// Uint32 reads a little endian uint32 from the underlying buffer.
func (r *Reader) Uint32(x *uint32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = binary.BigEndian.Uint32(b)
}

// Int32 reads a little endian int32 from the underlying buffer.
func (r *Reader) Int32(x *int32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = int32(binary.BigEndian.Uint32(b))
}

// BEInt32 reads a big endian int32 from the underlying buffer.
func (r *Reader) BEInt32(x *int32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = *(*int32)(unsafe.Pointer(&b[0]))
}

// Uint64 reads a little endian uint64 from the underlying buffer.
func (r *Reader) Uint64(x *uint64) {
	b := make([]byte, 8)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = binary.BigEndian.Uint64(b)
}

// Int64 reads a little endian int64 from the underlying buffer.
func (r *Reader) Int64(x *int64) {
	b := make([]byte, 8)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = int64(binary.BigEndian.Uint64(b))
}

// Float32 reads a little endian float32 from the underlying buffer.
func (r *Reader) Float32(x *float32) {
	b := make([]byte, 4)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}
	*x = math.Float32frombits(binary.BigEndian.Uint32(b))
}

// Bool reads a bool from the underlying buffer.
func (r *Reader) Bool(x *bool) {
	u, err := r.r.ReadByte()
	if err != nil {
		r.panic(err)
	}
	*x = *(*bool)(unsafe.Pointer(&u))
}

// errStringTooLong is an error set if a string decoded using the String method has a length that is too long.
var errStringTooLong = errors.New("string length overflows a 32-bit integer")

// StringUTF ...
func (r *Reader) StringUTF(x *string) {
	var length int16
	r.Int16(&length)
	l := int(length)
	if l > math.MaxInt16 {
		r.panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		r.panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// String reads a string from the underlying buffer.
func (r *Reader) String(x *string) {
	var length uint32
	r.Varuint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		r.panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		r.panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// ByteSlice reads a byte slice from the underlying buffer, similarly to String.
func (r *Reader) ByteSlice(x *[]byte) {
	var length uint32
	r.Varuint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		r.panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		r.panic(err)
	}
	*x = data
}

// ByteFloat reads a rotational float32 from a single byte.
func (r *Reader) ByteFloat(x *float32) {
	var v uint8
	r.Uint8(&v)
	*x = float32(v) * (360.0 / 256.0)
}

// RGB reads a color.RGBA x from three float32s.
func (r *Reader) RGB(x *color.RGBA) {
	var red, green, blue float32
	r.Float32(&red)
	r.Float32(&green)
	r.Float32(&blue)
	*x = color.RGBA{
		R: uint8(red * 255),
		G: uint8(green * 255),
		B: uint8(blue * 255),
	}
}

// RGBA reads a color.RGBA x from a uint32.
func (r *Reader) RGBA(x *color.RGBA) {
	var v uint32
	r.Uint32(&v)
	*x = color.RGBA{
		R: byte(v),
		G: byte(v >> 8),
		B: byte(v >> 16),
		A: byte(v >> 24),
	}
}

// VarRGBA reads a color.RGBA x from a varuint32.
func (r *Reader) VarRGBA(x *color.RGBA) {
	var v uint32
	r.Varuint32(&v)
	*x = color.RGBA{
		R: byte(v),
		G: byte(v >> 8),
		B: byte(v >> 16),
		A: byte(v >> 24),
	}
}

// Bytes reads the leftover bytes into a byte slice.
func (r *Reader) Bytes(p *[]byte) {
	var err error
	*p, err = io.ReadAll(r.r)
	if err != nil {
		r.panic(err)
	}
}

// UUID reads a uuid.UUID from the underlying buffer.
func (r *Reader) UUID(x *uuid.UUID) {
	b := make([]byte, 16)
	if _, err := r.r.Read(b); err != nil {
		r.panic(err)
	}

	// The UUIDs we read are Little Endian, but the uuid library is based on Big Endian UUIDs, so we need to
	// reverse the two int64s the UUID is composed of, then reverse their bytes too.
	b = append(b[8:], b[:8]...)
	var arr [16]byte
	for i, j := 0, 15; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = b[j], b[i]
	}
	*x = arr
}

// LimitUint32 checks if the value passed is lower than the limit passed. If not, the Reader panics.
func (r *Reader) LimitUint32(value uint32, max uint32) {
	if max == math.MaxUint32 {
		// Account for 0-1 overflowing into max.
		max = 0
	}
	if value > max {
		r.panicf("uint32 %v exceeds maximum of %v", value, max)
	}
}

// LimitInt32 checks if the value passed is lower than the limit passed and higher than the minimum. If not,
// the Reader panics.
func (r *Reader) LimitInt32(value int32, min, max int32) {
	if value < min {
		r.panicf("int32 %v exceeds minimum of %v", value, min)
	} else if value > max {
		r.panicf("int32 %v exceeds maximum of %v", value, max)
	}
}

// ShieldID returns the shield ID provided to the reader.
func (r *Reader) ShieldID() int32 {
	return r.shieldID
}

// UnknownEnumOption panics with an unknown enum option error.
func (r *Reader) UnknownEnumOption(value any, enum string) {
	r.panicf("unknown value '%v' for enum type '%v'", value, enum)
}

// InvalidValue panics with an error indicating that the value passed is not valid for a specific field.
func (r *Reader) InvalidValue(value any, forField, reason string) {
	r.panicf("invalid value '%v' for %v: %v", value, forField, reason)
}

// errVarIntOverflow is an error set if one of the Varint methods encounters a varint that does not terminate
// after 5 or 10 bytes, depending on the data type read into.
var errVarIntOverflow = errors.New("varint overflows integer")

// Varint64 reads up to 10 bytes from the underlying buffer into an int64.
func (r *Reader) Varint64(x *int64) {
	var ux uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			r.panic(err)
		}

		ux |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int64(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	r.panic(errVarIntOverflow)
}

// Varuint64 reads up to 10 bytes from the underlying buffer into a uint64.
func (r *Reader) Varuint64(x *uint64) {
	var v uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			r.panic(err)
		}

		v |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	r.panic(errVarIntOverflow)
}

// Varint32 reads up to 5 bytes from the underlying buffer into an int32.
func (r *Reader) Varint32(x *int32) {
	var ux uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			r.panic(err)
		}

		ux |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int32(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	r.panic(errVarIntOverflow)
}

// Varuint32 reads up to 5 bytes from the underlying buffer into a uint32.
func (r *Reader) Varuint32(x *uint32) {
	var v uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			r.panic(err)
		}

		v |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	r.panic(errVarIntOverflow)
}

// panicf panics with the format and values passed and assigns the error created to the Reader.
func (r *Reader) panicf(format string, a ...any) {
	panic(fmt.Errorf(format, a...))
}

// panic panics with the error passed, similarly to panicf.
func (r *Reader) panic(err error) {
	panic(err)
}
