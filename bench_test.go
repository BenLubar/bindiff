package bindiff

import (
	"math/rand"
	"testing"
)

func BenchmarkDiff_1B(b *testing.B)  { benchmarkDiff(b, 1<<0) }
func BenchmarkDiff_32B(b *testing.B) { benchmarkDiff(b, 1<<5) }
func BenchmarkDiff_1K(b *testing.B)  { benchmarkDiff(b, 1<<10) }
func BenchmarkDiff_32K(b *testing.B) { benchmarkDiff(b, 1<<15) }

func BenchmarkForward_1B(b *testing.B)  { benchmarkForward(b, 1<<0) }
func BenchmarkForward_32B(b *testing.B) { benchmarkForward(b, 1<<5) }
func BenchmarkForward_1K(b *testing.B)  { benchmarkForward(b, 1<<10) }
func BenchmarkForward_32K(b *testing.B) { benchmarkForward(b, 1<<15) }

func BenchmarkReverse_1B(b *testing.B)  { benchmarkReverse(b, 1<<0) }
func BenchmarkReverse_32B(b *testing.B) { benchmarkReverse(b, 1<<5) }
func BenchmarkReverse_1K(b *testing.B)  { benchmarkReverse(b, 1<<10) }
func BenchmarkReverse_32K(b *testing.B) { benchmarkReverse(b, 1<<15) }

func benchmarkDiff(b *testing.B, size int) {
	x, y := benchmarkSetup(size)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Diff(x, y, 0)
	}
}

func benchmarkForward(b *testing.B, size int) {
	x, y := benchmarkSetup(size)
	d := Diff(x, y, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Forward(x, d)
		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkReverse(b *testing.B, size int) {
	x, y := benchmarkSetup(size)
	d := Diff(x, y, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Reverse(y, d)
		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkSetup(size int) (x, y []byte) {
	r := rand.New(rand.NewSource(0))
	x, y = make([]byte, size), make([]byte, size)
	for i := range x {
		x[i] = byte(r.Intn(1 << 8))
	}
	copy(y, x)
	for i := range y {
		if r.Intn(4) == 0 {
			y[i] = byte(r.Intn(1 << 8))
		}
	}
	return
}
