package httprouter

import (
	"fmt"

	"github.com/makasim/httprouter/radix"
	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

var HandlerKeyUserValue = "fasthttprouter.handler_id"

type Router struct {
	PageNotFoundHandler     fasthttp.RequestHandler
	MethodNotAllowedHandler fasthttp.RequestHandler
	GlobalHandler           fasthttp.RequestHandler
	Handlers                map[uint64]fasthttp.RequestHandler

	Trees []radix.Tree
}

func New() *Router {
	return &Router{
		PageNotFoundHandler: func(ctx *fasthttp.RequestCtx) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		},
		MethodNotAllowedHandler: func(ctx *fasthttp.RequestCtx) {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		},
		Handlers: make(map[uint64]fasthttp.RequestHandler),

		Trees: make([]radix.Tree, 9),
	}
}

func (r *Router) Handle(ctx *fasthttp.RequestCtx) {
	i := r.methodIndexOf(gotils.B2S(ctx.Method()))
	if i == -1 {
		r.MethodNotAllowedHandler(ctx)
		return
	}

	hID := r.Trees[i].Search(gotils.B2S(ctx.Path()), ctx.SetUserValue)
	if hID == 0 {
		r.PageNotFoundHandler(ctx)
		return
	}

	ctx.SetUserValue(HandlerKeyUserValue, hID)

	if h, ok := r.Handlers[hID]; ok {
		h(ctx)
		return
	}

	if r.GlobalHandler != nil {
		r.GlobalHandler(ctx)
		return
	}

	r.PageNotFoundHandler(ctx)
}

// Add adds a route for method and path to the router
// It is not safe for concurrent use.
// Add routes before using Handle or protect Add, Remove, Handle with mutex.
func (r *Router) Add(method, path string, handlerID uint64) error {
	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		return fmt.Errorf("method not allowed")
	}
	if len(path) == 0 {
		return fmt.Errorf("path empty")
	}

	var err error
	tree := r.Trees[methodIndex].Clone()

	tree, err = tree.Insert(path, handlerID)
	if err != nil {
		return err
	}

	r.Trees[methodIndex] = tree

	return nil
}

// Remove removes a route for method and path frmo the router
// It is not safe for concurrent use.
//Remove routes before using Handle or protect Add, Remove, Handle with mutex.
func (r *Router) Remove(method, path string) error {
	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		return fmt.Errorf("method not allowed")
	}
	if len(path) == 0 {
		return fmt.Errorf("path empty")
	}

	var err error
	tree := r.Trees[methodIndex].Clone()

	tree, err = tree.Delete(path)
	if err != nil {
		return err
	}

	r.Trees[methodIndex] = tree
	return nil
}

func (r *Router) methodIndexOf(method string) int {
	switch method {
	case fasthttp.MethodGet:
		return 0
	case fasthttp.MethodHead:
		return 1
	case fasthttp.MethodPost:
		return 2
	case fasthttp.MethodPut:
		return 3
	case fasthttp.MethodPatch:
		return 4
	case fasthttp.MethodDelete:
		return 5
	case fasthttp.MethodConnect:
		return 6
	case fasthttp.MethodOptions:
		return 7
	case fasthttp.MethodTrace:
		return 8
	}

	return -1
}
