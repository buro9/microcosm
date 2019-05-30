package bpool

// WrapByteSlice wraps a []byte as a ByteSlice
func WrapByteSlice(full []byte, headerLength int) ByteSlice {
	return ByteSlice{
		full:    full,
		current: full[headerLength:],
	}
}

// ByteSlice provides a wrapper around []byte with some added convenience
type ByteSlice struct {
	full    []byte
	current []byte
}

// ResliceTo reslices the end of the current slice.
func (b ByteSlice) ResliceTo(end int) ByteSlice {
	return ByteSlice{
		full:    b.full,
		current: b.current[:end],
	}
}

// Bytes returns the current slice
func (b ByteSlice) Bytes() []byte {
	return b.current
}

// Full returns the full original buffer underlying the ByteSlice
func (b ByteSlice) Full() []byte {
	return b.full
}

// ByteSlicePool is a bool of byte slices
type ByteSlicePool interface {
	// Get gets a byte slice from the pool
	GetSlice() ByteSlice
	// Put returns a byte slice to the pool
	PutSlice(ByteSlice)
	// NumPooled returns the number of currently pooled items
	NumPooled() int
}

// NewByteSlicePool creates a new ByteSlicePool bounded to the
// given maxSize, with new byte arrays sized based on width
func NewByteSlicePool(maxSize int, width int) ByteSlicePool {
	return NewHeaderPreservingByteSlicePool(maxSize, width, 0)
}

// NewHeaderPreservingByteSlicePool creates a new ByteSlicePool bounded to the
// given maxSize, with new byte arrays sized based on width and headerLength
// preserved at the beginning of the slice.
func NewHeaderPreservingByteSlicePool(maxSize int, width int, headerLength int) ByteSlicePool {
	return &BytePool{
		c: make(chan []byte, maxSize),
		w: width + headerLength,
		h: headerLength,
	}
}

// GetSlice implements the method from interface ByteSlicePool
func (bp *BytePool) GetSlice() ByteSlice {
	return WrapByteSlice(bp.Get(), bp.h)
}

// PutSlice implements the method from interface ByteSlicePool
func (bp *BytePool) PutSlice(b ByteSlice) {
	bp.Put(b.Full())
}
