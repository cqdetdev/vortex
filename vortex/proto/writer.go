package proto

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"math"
	"unsafe"

	"github.com/google/uuid"
)

// Writer implements writing methods for data types from Minecraft packets. Each Packet implementation has one
// passed to it when writing.
// Writer implements methods where values are passed using a pointer, so that Reader and Writer have a
// synonymous interface and both implement the IO interface.
type Writer struct {
	w interface {
		io.Writer
		io.ByteWriter
	}
	shieldID int32
}

// NewWriter creates a new initialised Writer with an underlying io.ByteWriter to write to.
func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}, shieldID int32) *Writer {
	return &Writer{w: w, shieldID: shieldID}
}

// Uint8 writes a uint8 to the underlying buffer.
func (w *Writer) Uint8(x *uint8) {
	_ = w.w.WriteByte(*x)
}

// Int8 writes an int8 to the underlying buffer.
func (w *Writer) Int8(x *int8) {
	_ = w.w.WriteByte(byte(*x) & 0xff)
}

// Bool writes a bool as either 0 or 1 to the underlying buffer.
func (w *Writer) Bool(x *bool) {
	_ = w.w.WriteByte(*(*byte)(unsafe.Pointer(x)))
}

// Uint16 writes a little endian uint16 to the underlying buffer.
func (w *Writer) Uint16(x *uint16) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, *x)
	_, _ = w.w.Write(data)
}

// Int16 writes a little endian int16 to the underlying buffer.
func (w *Writer) Int16(x *int16) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(*x))
	_, _ = w.w.Write(data)
}

// Uint32 writes a little endian uint32 to the underlying buffer.
func (w *Writer) Uint32(x *uint32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, *x)
	_, _ = w.w.Write(data)
}

// Int32 writes a little endian int32 to the underlying buffer.
func (w *Writer) Int32(x *int32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(*x))
	_, _ = w.w.Write(data)
}

// BEInt32 writes a big endian int32 to the underlying buffer.
func (w *Writer) BEInt32(x *int32) {
	data := *(*[4]byte)(unsafe.Pointer(x))
	_, _ = w.w.Write(data[:])
}

// Uint64 writes a little endian uint64 to the underlying buffer.
func (w *Writer) Uint64(x *uint64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, *x)
	_, _ = w.w.Write(data)
}

// Int64 writes a little endian int64 to the underlying buffer.
func (w *Writer) Int64(x *int64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(*x))
	_, _ = w.w.Write(data)
}

// Float32 writes a little endian float32 to the underlying buffer.
func (w *Writer) Float32(x *float32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, math.Float32bits(*x))
	_, _ = w.w.Write(data)
}

// StringUTF ...
func (w *Writer) StringUTF(x *string) {
	l := int16(len(*x))
	w.Int16(&l)
	_, _ = w.w.Write([]byte(*x))
}

// String writes a string, prefixed with a varuint32, to the underlying buffer.
func (w *Writer) String(x *string) {
	l := uint32(len(*x))
	w.Varuint32(&l)
	_, _ = w.w.Write([]byte(*x))
}

// ByteSlice writes a []byte, prefixed with a varuint32, to the underlying buffer.
func (w *Writer) ByteSlice(x *[]byte) {
	l := uint32(len(*x))
	w.Varuint32(&l)
	_, _ = w.w.Write(*x)
}

// Bytes appends a []byte to the underlying buffer.
func (w *Writer) Bytes(x *[]byte) {
	_, _ = w.w.Write(*x)
}

// ByteFloat writes a rotational float32 as a single byte to the underlying buffer.
func (w *Writer) ByteFloat(x *float32) {
	_ = w.w.WriteByte(byte(*x / (360.0 / 256.0)))
}

// RGB writes a color.RGBA x as 3 float32s to the underlying buffer.
func (w *Writer) RGB(x *color.RGBA) {
	red := float32(x.R) / 255
	green := float32(x.G) / 255
	blue := float32(x.B) / 255
	w.Float32(&red)
	w.Float32(&green)
	w.Float32(&blue)
}

// RGBA writes a color.RGBA x as a uint32 to the underlying buffer.
func (w *Writer) RGBA(x *color.RGBA) {
	val := uint32(x.R) | uint32(x.G)<<8 | uint32(x.B)<<16 | uint32(x.A)<<24
	w.Uint32(&val)
}

// VarRGBA writes a color.RGBA x as a varuint32 to the underlying buffer.
func (w *Writer) VarRGBA(x *color.RGBA) {
	val := uint32(x.R) | uint32(x.G)<<8 | uint32(x.B)<<16 | uint32(x.A)<<24
	w.Varuint32(&val)
}

// UUID writes a UUID to the underlying buffer.
func (w *Writer) UUID(x *uuid.UUID) {
	b := append((*x)[8:], (*x)[:8]...)
	for i, j := 0, 15; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	_, _ = w.w.Write(b)
}

// Varint64 writes an int64 as 1-10 bytes to the underlying buffer.
func (w *Writer) Varint64(x *int64) {
	u := *x
	ux := uint64(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint64 writes a uint64 as 1-10 bytes to the underlying buffer.
func (w *Writer) Varuint64(x *uint64) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// Varint32 writes an int32 as 1-5 bytes to the underlying buffer.
func (w *Writer) Varint32(x *int32) {
	u := *x
	ux := uint32(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint32 writes a uint32 as 1-5 bytes to the underlying buffer.
func (w *Writer) Varuint32(x *uint32) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// ShieldID returns the shield ID provided to the writer.
func (w *Writer) ShieldID() int32 {
	return w.shieldID
}

// UnknownEnumOption panics with an unknown enum option error.
func (w *Writer) UnknownEnumOption(value any, enum string) {
	w.panicf("unknown value '%v' for enum type '%v'", value, enum)
}

// InvalidValue panics with an invalid value error.
func (w *Writer) InvalidValue(value any, forField, reason string) {
	w.panicf("invalid value '%v' for %v: %v", value, forField, reason)
}

// panicf panics with the format and values passed.
func (w *Writer) panicf(format string, a ...any) {
	panic(fmt.Errorf(format, a...))
}
