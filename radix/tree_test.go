package radix_test

import (
	"testing"

	"github.com/makasim/httprouter/radix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreeInsertEmptyPath(t *testing.T) {
	tree, err := radix.NewTree().Insert("", 2)
	require.EqualError(t, err, "insert: path empty")
	require.Equal(t, radix.Tree{}, tree)
}

func TestTreeInsertInvalidPath(t *testing.T) {
	tree, err := radix.NewTree().Insert("foo", 2)
	require.EqualError(t, err, "insert: path must start with /")
	require.Equal(t, radix.Tree{}, tree)
}

func TestTreeInsertEmptyKey(t *testing.T) {
	tree, err := radix.NewTree().Insert("/foo", 0)
	require.EqualError(t, err, "insert: key empty")
	require.Equal(t, radix.Tree{}, tree)
}

func TestTree(t *testing.T) {
	tree := radix.NewTree()

	tree, err := tree.Insert("/foo", 1)
	require.NoError(t, err)

	tree, err = tree.Insert("/faa", 2)
	require.NoError(t, err)

	tree, err = tree.Insert("/fab", 3)
	require.NoError(t, err)

	tree, err = tree.Insert("/ccc", 4)
	require.NoError(t, err)

	tree, err = tree.Insert("/cc/bb/aa", 5)
	require.NoError(t, err)

	assert.Equal(t, uint64(0), tree.Search("/", dummyKV()))
	assert.Equal(t, uint64(0), tree.Search("/fo", dummyKV()))
	assert.Equal(t, uint64(0), tree.Search("/fooo", dummyKV()))
	assert.Equal(t, uint64(1), tree.Search("/foo", dummyKV()))
	assert.Equal(t, uint64(2), tree.Search("/faa", dummyKV()))
	assert.Equal(t, uint64(3), tree.Search("/fab", dummyKV()))

	assert.Equal(t, uint64(0), tree.Search("/cc", dummyKV()))
	assert.Equal(t, uint64(0), tree.Search("/cccc", dummyKV()))
	assert.Equal(t, uint64(4), tree.Search("/ccc", dummyKV()))

	assert.Equal(t, uint64(5), tree.Search("/cc/bb/aa", dummyKV()))

	tree, err = tree.Delete("/faa")
	require.NoError(t, err)

	tree, err = tree.Delete("/foooo")
	require.NoError(t, err)

	tree, err = tree.Delete("/fo")
	require.NoError(t, err)

	assert.Equal(t, uint64(0), tree.Search("/faa", dummyKV()))
	assert.Equal(t, uint64(1), tree.Search("/foo", dummyKV()))
}

func TestTreeDelete(t *testing.T) {
	tree := radix.NewTree()

	tree, err := tree.Insert("/foo", 1)
	require.NoError(t, err)

	// guard
	assert.Equal(t, uint64(1), tree.Search("/foo", dummyKV()))

	tree, err = tree.Delete("/foo")
	require.NoError(t, err)

	assert.Equal(t, uint64(0), tree.Search("/foo", dummyKV()))
}

func TestTreeCount(t *testing.T) {
	t0 := radix.Tree{}
	assert.Equal(t, 0, t0.Count())

	var t1 radix.Tree
	assert.Equal(t, 0, t1.Count())

	t2 := radix.NewTree()
	t2, err := t2.Insert("/foo", 123)
	require.NoError(t, err)
	t2, err = t2.Insert("/bar", 124)
	require.NoError(t, err)
	t2, err = t2.Insert("/foo/bar", 125)
	require.NoError(t, err)
	t2, err = t2.Insert("/{name}", 126)
	require.NoError(t, err)
	assert.Equal(t, 4, t2.Count())

	t2, err = t2.Delete("/foo")
	require.NoError(t, err)
	t2, err = t2.Delete("/bar")
	require.NoError(t, err)
	t2, err = t2.Delete("/foo/bar")
	require.NoError(t, err)
	t2, err = t2.Delete("/{name}")
	require.NoError(t, err)
	assert.Equal(t, 0, t2.Count())
}

func dummyKV() func(n string, v interface{}) {
	return func(n string, v interface{}) {
	}
}
