package radix

import (
	"testing"
)

var nodeKey uint64

func BenchmarkSearchStatic(b *testing.B) {

	n := Node{}
	n = n.Insert("/a", 1)
	n = n.Insert("/foo/bar1", 1)
	n = n.Insert("/foo/bar2", 2)
	n = n.Insert("/foo/bar3", 3)
	n = n.Insert("/foo/bar4", 4)
	n = n.Insert("/foo/bar5", 5)
	n = n.Insert("/foo/bar6", 6)
	n = n.Insert("/foo/bar7", 7)
	n = n.Insert("/foo/bar8", 8)
	n = n.Insert("/foo/bar9", 9)
	n = n.Insert("/foo/bar10", 10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeKey = n.Search("foo/bar5", func(n string, v interface{}) {

		})
		if nodeKey != 5 {
			b.Fatalf("result %d is not equal to 5", nodeKey)
		}
	}
}

func BenchmarkSearchDynamic(b *testing.B) {
	n := Node{}
	n = n.Insert("/a", 1)
	n = n.Insert("/foo/bar1", 1)
	n = n.Insert("/foo/bar2", 2)
	n = n.Insert("/foo/bar3", 3)
	n = n.Insert("/foo/bar4", 4)
	n = n.Insert("/foo/bar5", 5)
	n = n.Insert("/foo/bar6", 6)
	n = n.Insert("/foo/bar7", 7)
	n = n.Insert("/foo/bar8", 8)
	n = n.Insert("/foo/bar9", 9)
	n = n.Insert("/foo/{param}", 10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeKey = n.Search("foo/custom", func(n string, v interface{}) {

		})
		if nodeKey != 10 {
			b.Fatalf("result %d is not equal to 10", nodeKey)
		}
	}
}
