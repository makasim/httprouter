package radix

import (
	"sync"
	"testing"
)

var t Tree
var key uint64

// compare with https://github.com/fasthttp/router/blob/5c77f27ae28987b4cbb007be06f6ef793cdb062d/radix/tree_test.go#L299
// fasthttp Benchmark_Get-8                       	33380914	        30.3 ns/op
// our      Benchmark_Get-8             	        79250941	        15.1 ns/op
func Benchmark_Get(b *testing.B) {
	tree := NewTree()
	tree, _ = tree.Insert("/", 1)
	tree, _ = tree.Insert("/plaintext", 2)
	tree, _ = tree.Insert("/json", 3)
	tree, _ = tree.Insert("/fortune", 4)
	tree, _ = tree.Insert("/fortune-quick", 5)
	tree, _ = tree.Insert("/db", 6)
	tree, _ = tree.Insert("/queries", 7)
	tree, _ = tree.Insert("/update", 7)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key = tree.Search("/update", func(n string, v interface{}) {
		})
	}
}

func Benchmark_GetLock(b *testing.B) {
	tree := NewTree()
	tree, _ = tree.Insert("/", 1)
	tree, _ = tree.Insert("/plaintext", 2)
	tree, _ = tree.Insert("/json", 3)
	tree, _ = tree.Insert("/fortune", 4)
	tree, _ = tree.Insert("/fortune-quick", 5)
	tree, _ = tree.Insert("/db", 6)
	tree, _ = tree.Insert("/queries", 7)
	tree, _ = tree.Insert("/update", 7)

	mu := &sync.Mutex{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		key = tree.Search("/update", func(n string, v interface{}) {
		})
		mu.Unlock()
	}
}

// compare with https://github.com/fasthttp/router/blob/5c77f27ae28987b4cbb007be06f6ef793cdb062d/radix/tree_test.go#L328
// fasthttp Benchmark_GetWithParams-8   	12547896	        96.2 ns/op
// our      Benchmark_GetWithParams-8   	15075598	        79.1 ns/op	      16 B/op	       1 allocs/op
func Benchmark_GetWithParams(b *testing.B) {
	tree := NewTree()
	tree, _ = tree.Insert("/api/{version}/data", 1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key = tree.Search("/api/v1/data", func(n string, v interface{}) {
		})
	}
}

func Benchmark_Insert(b *testing.B) {
	tree := NewTree()
	for i := 0; i < b.N; i++ {
		var err error
		t, err = tree.Insert("/foo", uint64(b.N))
		if err != nil {
			b.Error("Update failed")
		}
	}
}

func Benchmark_InsertLock(b *testing.B) {
	tree := NewTree()
	mu := &sync.Mutex{}

	for i := 0; i < b.N; i++ {
		var err error
		mu.Lock()
		t, err = tree.Insert("/foo", uint64(b.N))
		mu.Unlock()
		if err != nil {
			b.Error("Update failed")
		}
	}
}

func Benchmark_InsertClone(b *testing.B) {
	tree := NewTree()
	for i := 0; i < b.N; i++ {
		var err error
		ct := tree.Clone()
		t, err = ct.Insert("/foo", uint64(b.N))
		if err != nil {
			b.Error("Update failed")
		}
	}
}
