package httprouter

import (
	"errors"
	"fmt"
	"github.com/makasim/httprouter/radix"
	"github.com/savsgio/gotils"
)

var ErrPathNotFound = fmt.Errorf("path not found")

type Router struct {
	tree radix.Tree
}

func (r Router) Route(method, path []byte, kv func(n string, v interface{})) (uint64, error) {
	if len(method) == 0 {
		return 0, fmt.Errorf("method empty")
	}
	if len(path) == 0 {
		return 0, fmt.Errorf("path empty")
	}

	path1 := "/" + gotils.B2S(method) + "/" + gotils.B2S(path)

	handlerID := r.tree.Search(path1, kv)
	if handlerID == 0 {
		return 0, ErrPathNotFound
	}

	return handlerID, nil
}

func (r Router) Insert(method, path []byte, handlerID uint64) (Router, error) {
	if len(method) == 0 {
		return Router{}, fmt.Errorf("method empty")
	}
	if len(path) == 0 {
		return Router{}, fmt.Errorf("path empty")
	}

	path1 := "/" + gotils.B2S(method) + "/" + gotils.B2S(path)

	var err error
	tree := r.tree.Clone()

	tree, err = tree.Insert(path1, handlerID)
	if err != nil {
		return Router{}, err
	}

	return Router{tree: tree}, nil
}

func (r Router) Delete(method, path []byte) (Router, error) {
	if len(method) == 0 {
		return Router{}, fmt.Errorf("method empty")
	}
	if len(path) == 0 {
		return Router{}, fmt.Errorf("path empty")
	}

	path1 := "/" + gotils.B2S(method) + "/" + gotils.B2S(path)

	var err error
	tree := r.tree.Clone()

	tree, err = tree.Delete(path1)
	if err != nil && !errors.Is(err, ErrPathNotFound) {
		return Router{}, err
	}

	return Router{tree: tree}, nil
}

// Tree could be used for debugging purpose
func (r Router) Tree() radix.Tree {
	return r.tree.Clone()
}

func (r Router) Count() int {
	return r.tree.Count()
}
