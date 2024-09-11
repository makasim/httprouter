package stdrouter

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/makasim/httprouter/radix"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, Params)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, p Params) {
	f(rw, r, p)
}

type HandlerID int

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) Get(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

func (ps *Params) Set(name, value string) {
	if value == "" {
		for i, p := range *ps {
			if p.Key == name {
				*ps = append((*ps)[:i], (*ps)[i+1:]...)
			}
		}

		return
	}

	for i, p := range *ps {
		if p.Key == name {
			(*ps)[i].Value = value
			return
		}
	}

	*ps = append(*ps, Param{
		Key:   name,
		Value: value,
	})
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

// ParamsFromContext pulls the URL parameters from a request context,
// or returns nil if none are present.
func ParamsFromContext(ctx context.Context) Params {
	p, ok := ctx.Value(ParamsKey).(Params)
	if !ok {
		return Params{}
	}

	return p
}

var HandlerKeyUserValue = "stdprouter.handler_id"

const MethodAny = "ANY"
const methodAnyIndex = 9

type Router struct {
	PageNotFoundHandler     http.HandlerFunc
	MethodNotAllowedHandler http.HandlerFunc
	GlobalHandler           Handler

	handlers       []Handler
	freeHandlerIds []HandlerID

	Trees []radix.Tree

	paramsPool sync.Pool
}

func New() *Router {
	return &Router{
		PageNotFoundHandler: func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		},
		MethodNotAllowedHandler: func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		},

		Trees: make([]radix.Tree, 10),

		handlers:       make([]Handler, 1), // 0 is nil handler
		freeHandlerIds: make([]HandlerID, 0),

		paramsPool: sync.Pool{
			New: func() interface{} {
				return new(Params)
			},
		},
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	i := methodIndexOf(req.Method)
	if i == -1 {
		r.MethodNotAllowedHandler(rw, req)
		return
	}

	ps := r.getParams()
	defer r.putParams(ps)

	hID := r.Trees[i].Search(req.URL.Path, func(n string, v interface{}) {
		v1, ok := v.(string)
		if !ok {
			return // skip
		}

		*ps = append(*ps, Param{
			Key:   n,
			Value: v1,
		})
	})
	if hID == 0 {
		if ps != nil {
			*ps = (*ps)[:0]
		}

		hID = r.Trees[methodAnyIndex].Search(req.URL.Path, func(n string, v interface{}) {
			v1, ok := v.(string)
			if !ok {
				return // skip
			}

			*ps = append(*ps, Param{
				Key:   n,
				Value: v1,
			})
		})

		if hID == 0 {
			r.PageNotFoundHandler(rw, req)
			return
		}
	}

	maxHID := len(r.handlers) - 1
	if int(hID) <= maxHID {
		if h := r.handlers[int(hID)]; h != nil {
			h.ServeHTTP(rw, req, *ps)
			return
		}
	}

	if r.GlobalHandler != nil {

		r.GlobalHandler.ServeHTTP(rw, req, *ps)
		return
	}

	r.PageNotFoundHandler(rw, req)
}

func (r *Router) AddHandler(handler Handler) HandlerID {
	if handler == nil {
		panic("handler is nil")
	}

	if len(r.freeHandlerIds) > 0 {
		id := r.freeHandlerIds[len(r.freeHandlerIds)-1]
		r.freeHandlerIds = r.freeHandlerIds[:len(r.freeHandlerIds)-1]

		r.handlers[id] = handler

		return id
	}

	id := len(r.handlers)
	r.handlers = append(r.handlers, handler)

	return HandlerID(id)
}

func (r *Router) FindHandler(method, path string) (Handler, error) {
	i := methodIndexOf(method)
	if i == -1 {
		return nil, fmt.Errorf("unsupported method %v", method)
	}

	hID := r.Trees[i].Search(path, func(n string, v interface{}) {})
	if hID == 0 && i != methodAnyIndex {
		hID = r.Trees[methodAnyIndex].Search(path, func(n string, v interface{}) {})

		if hID == 0 {
			return nil, fmt.Errorf("path %v not found", path)
		}
	}

	maxHID := len(r.handlers) - 1
	if int(hID) <= maxHID {
		if h := r.handlers[int(hID)]; h != nil {
			return h, nil
		}
	}

	if r.GlobalHandler != nil {
		return r.GlobalHandler, nil
	}

	return nil, fmt.Errorf("handler not found")
}

func (r *Router) RemoveHandler(hID HandlerID) {
	r.handlers[hID] = nil
	r.freeHandlerIds = append(r.freeHandlerIds, hID)
}

func (r *Router) GetHandler(hID HandlerID) (Handler, error) {
	if slices.Contains(r.freeHandlerIds, hID) {
		return nil, fmt.Errorf("handler not found")
	}
	if hID < 1 || int(hID) >= len(r.handlers) {
		return nil, fmt.Errorf("handler not found")
	}

	return r.handlers[hID], nil
}

func (r *Router) AddStdHandler(handler http.Handler) HandlerID {
	return r.AddHandler(HandlerFunc(func(rw http.ResponseWriter, req *http.Request, _ Params) {
		handler.ServeHTTP(rw, req)
	}))
}

func (r *Router) RegisterHandler(method, path string, handler Handler) error {
	hID := r.AddHandler(handler)
	return r.Add(method, path, hID)
}

func (r *Router) Add(method, path string, handlerID HandlerID) error {
	methodIndex := methodIndexOf(method)
	if methodIndex == -1 {
		return fmt.Errorf("method not allowed")
	}
	if len(path) == 0 {
		return fmt.Errorf("path empty")
	}

	var err error
	tree := r.Trees[methodIndex]

	tree, err = tree.Insert(path, uint64(handlerID))
	if err != nil {
		return err
	}

	r.Trees[methodIndex] = tree

	return nil
}

func (r *Router) Remove(method, path string) error {
	methodIndex := methodIndexOf(method)
	if methodIndex == -1 {
		return fmt.Errorf("method not allowed")
	}
	if len(path) == 0 {
		return fmt.Errorf("path empty")
	}

	var err error
	tree := r.Trees[methodIndex]

	tree, err = tree.Delete(path)
	if err != nil {
		return err
	}

	r.Trees[methodIndex] = tree
	return nil
}

func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

func methodIndexOf(method string) int {
	switch method {
	case http.MethodGet:
		return 0
	case http.MethodHead:
		return 1
	case http.MethodPost:
		return 2
	case http.MethodPut:
		return 3
	case http.MethodPatch:
		return 4
	case http.MethodDelete:
		return 5
	case http.MethodConnect:
		return 6
	case http.MethodOptions:
		return 7
	case http.MethodTrace:
		return 8
	case MethodAny:
		return methodAnyIndex
	}

	return -1
}
