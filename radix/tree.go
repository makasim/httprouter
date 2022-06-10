package radix

import (
	"fmt"
)

type Tree struct {
	root Node
}

func NewTree() Tree {
	return Tree{}
}

func (t Tree) Insert(path string, key uint64) (tree Tree, err error) {
	if path == "" {
		return Tree{}, fmt.Errorf("insert: path empty")
	}
	if string(path[0]) != "/" {
		return Tree{}, fmt.Errorf("insert: path must start with /")
	}

	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		if recErr, ok := rec.(error); ok {
			err = recErr
		} else {
			err = fmt.Errorf("%v", rec)
		}
	}()

	t.root = t.root.Insert(path, key)
	return t, err
}

func (t Tree) Delete(path string) (tree Tree, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	if path == "" {
		return Tree{}, fmt.Errorf("delete: path empty")
	}
	if string(path[0]) != "/" {
		return Tree{}, fmt.Errorf("delete: path must start with /")
	}

	i := longestCommonPrefix(path, t.root.path)
	if i >= 0 {
		path = path[i:]

		if path == "" && len(t.root.children) == 0 {
			t.root.path = ""
			t.root.key = 0
			return t, nil
		}

		if path == "" && len(t.root.children) > 0 {
			t.root.key = 0
			return t, nil
		}

		t.root = t.root.Delete(path)
		return t, nil
	}

	return t, nil
}
func (t Tree) Search(path string, kv func(n string, v interface{})) uint64 {
	if path == "" {
		return 0
	}
	if kv == nil {
		kv = func(n string, v interface{}) {}
	}

	if len(path) > len(t.root.path) {
		if path[:len(t.root.path)] != t.root.path {
			return 0
		}

		return t.root.Search(path[len(t.root.path):], kv)
	} else if len(path) == len(t.root.path) {
		return t.root.key
	}

	return 0
}

func (t Tree) Count() int {
	return t.root.Count()
}

func (t Tree) Clone() Tree {
	cloneTree := t
	cloneTree.root = t.root.Clone()

	return cloneTree
}

func (t Tree) String() string {
	return t.root.String()
}
