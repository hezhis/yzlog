package bufferpool

import (
	"time"
)

type Buffer struct {
	bs []byte
}

// AppendByte writes a single byte to the Buffer.
func (b *Buffer) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

// AppendTime appends the time formatted using the specified layout.
func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.bs = t.AppendFormat(b.bs, layout)
}

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.bs
}

// String returns a string copy of the underlying byte slice.
func (b *Buffer) String() string {
	return string(b.bs)
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

// Free returns the Buffer to its Pool.
//
// Callers must not retain references to the Buffer after calling Free.
func (b *Buffer) Free() {
	_pool.Put(b)
}
