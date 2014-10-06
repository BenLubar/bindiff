// Package bindiff provides a bidirectional binary patch for pairs of []byte.
package bindiff

import (
	"encoding/binary"
	"errors"
	"github.com/mb0/diff"
)

// ErrCorrupt is the only possible error from functions in this package.
var ErrCorrupt = errors.New("bindiff: corrupt patch")

// Diff computes the difference between old and new. A granularity of 1 or more
// combines changes with no greater than that many bytes between them.
func Diff(old, new []byte, granularity int) (patch []byte) {
	changes := diff.Bytes(old, new)
	if granularity > 0 {
		changes = diff.Granular(granularity, changes)
	}

	for i, c := range changes {
		a, b := c.A, c.B
		for _, prev := range changes[:i] {
			if prev.A < c.A {
				a -= prev.Del
				a += prev.Ins
			}
			if prev.B < c.B {
				b -= prev.Ins
				b += prev.Del
			}
		}
		patch = writeUvarint(patch, a)
		patch = writeUvarint(patch, b)
		patch = writeUvarint(patch, c.Del)
		patch = append(patch, old[c.A:c.A+c.Del]...)
		patch = writeUvarint(patch, c.Ins)
		patch = append(patch, new[c.B:c.B+c.Ins]...)
	}

	return
}

// Forward retrieves the second argument to Diff given the first argument and its output.
func Forward(old, patch []byte) (new []byte, err error) {
	return doPatch(old, patch, func(data []byte, a, b int, del, ins []byte) ([]byte, error) {
		return splice(data, a, len(del), ins)
	})
}

// Reverse retrieves the first argument to Diff given the second argument and its output.
func Reverse(new, patch []byte) (old []byte, err error) {
	return doPatch(new, patch, func(data []byte, a, b int, del, ins []byte) ([]byte, error) {
		return splice(data, b, len(ins), del)
	})
}

func doPatch(x, patch []byte, f func(data []byte, a, b int, del, ins []byte) ([]byte, error)) (y []byte, err error) {
	y = make([]byte, len(x))
	copy(y, x)

	for len(patch) > 0 {
		var a, b, l int
		var del, ins []byte

		a, patch, err = readUvarint(patch)
		if err != nil {
			return
		}

		b, patch, err = readUvarint(patch)
		if err != nil {
			return
		}

		l, patch, err = readUvarint(patch)
		if err != nil {
			return
		}
		if len(patch) < l {
			err = ErrCorrupt
			return
		}
		del, patch = patch[:l], patch[l:]

		l, patch, err = readUvarint(patch)
		if err != nil {
			return
		}
		if len(patch) < l {
			err = ErrCorrupt
			return
		}
		ins, patch = patch[:l], patch[l:]

		y, err = f(y, a, b, del, ins)
		if err != nil {
			return
		}
	}

	return
}

func splice(base []byte, anchor, del int, ins []byte) ([]byte, error) {
	if anchor+del > len(base) {
		return base, ErrCorrupt
	}
	l := len(base)
	if del < len(ins) {
		base = append(base, make([]byte, len(ins)-del)...)
	}
	copy(base[anchor+len(ins):], base[anchor+del:])
	copy(base[anchor:], ins)
	return base[:l+len(ins)-del], nil
}

func writeUvarint(buf []byte, i int) []byte {
	var varint [binary.MaxVarintLen64]byte
	return append(buf, varint[:binary.PutUvarint(varint[:], uint64(i))]...)
}

func readUvarint(buf []byte) (int, []byte, error) {
	i, n := binary.Uvarint(buf)
	if n == 0 {
		return 0, buf, ErrCorrupt
	}
	return int(i), buf[n:], nil
}
