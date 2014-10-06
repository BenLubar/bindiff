package bindiff

import (
	"testing"
	"testing/quick"
)

func TestForward(t *testing.T) {
	if err := quick.CheckEqual(func(a, b []byte) ([]byte, error) {
		diff := Diff(a, b, 0)
		return Forward(a, diff)
	}, func(a, b []byte) ([]byte, error) {
		return b, nil
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestReverse(t *testing.T) {
	if err := quick.CheckEqual(func(a, b []byte) ([]byte, error) {
		diff := Diff(a, b, 0)
		return Reverse(b, diff)
	}, func(a, b []byte) ([]byte, error) {
		return a, nil
	}, nil); err != nil {
		t.Error(err)
	}
}
