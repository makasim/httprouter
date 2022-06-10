package httprouter_test

import (
	"github.com/makasim/httprouter"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertEmptyMethod(t *testing.T) {
	r := httprouter.Router{}

	rr, err := r.Insert(nil, []byte("apath"), 123)
	require.EqualError(t, err, "method empty")
	require.Equal(t, httprouter.Router{}, rr)
}

func TestInsertEmptyPath(t *testing.T) {
	r := httprouter.Router{}

	rr, err := r.Insert([]byte("POST"), nil, 123)
	require.EqualError(t, err, "path empty")
	require.Equal(t, httprouter.Router{}, rr)
}

func TestInsert(t *testing.T) {
	var err error

	r := httprouter.Router{}

	r, err = r.Insert([]byte("GET"), []byte("/get"), 1)
	require.NoError(t, err)

	r, err = r.Insert([]byte("HEAD"), []byte("/head"), 2)
	require.NoError(t, err)

	r, err = r.Insert([]byte("POST"), []byte("/post"), 3)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PUT"), []byte("/put"), 4)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PATCH"), []byte("/patch"), 5)
	require.NoError(t, err)

	r, err = r.Insert([]byte("DELETE"), []byte("/delete"), 6)
	require.NoError(t, err)

	r, err = r.Insert([]byte("CONNECT"), []byte("/connect"), 7)
	require.NoError(t, err)

	r, err = r.Insert([]byte("OPTIONS"), []byte("/options"), 8)
	require.NoError(t, err)

	_, err = r.Insert([]byte("TRACE"), []byte("/trace"), 9)
	require.NoError(t, err)
}

func TestRoute(t *testing.T) {
	var err error

	r := httprouter.Router{}

	r, err = r.Insert([]byte("GET"), []byte("/get0"), 1)
	require.NoError(t, err)

	r, err = r.Insert([]byte("GET"), []byte("/get1"), 2)
	require.NoError(t, err)

	r, err = r.Insert([]byte("GET"), []byte("/get1/{param}"), 3)
	require.NoError(t, err)

	r, err = r.Insert([]byte("HEAD"), []byte("/head0"), 11)
	require.NoError(t, err)

	r, err = r.Insert([]byte("HEAD"), []byte("/head1"), 12)
	require.NoError(t, err)

	r, err = r.Insert([]byte("HEAD"), []byte("/head1/{param}"), 13)
	require.NoError(t, err)

	r, err = r.Insert([]byte("POST"), []byte("/post0"), 21)
	require.NoError(t, err)

	r, err = r.Insert([]byte("POST"), []byte("/post1"), 22)
	require.NoError(t, err)

	r, err = r.Insert([]byte("POST"), []byte("/post1/{param}"), 23)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PUT"), []byte("/put0"), 31)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PUT"), []byte("/put1"), 32)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PUT"), []byte("/put1/{param}"), 33)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PATCH"), []byte("/patch0"), 41)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PATCH"), []byte("/patch1"), 42)
	require.NoError(t, err)

	r, err = r.Insert([]byte("PATCH"), []byte("/patch1/{param}"), 43)
	require.NoError(t, err)

	r, err = r.Insert([]byte("DELETE"), []byte("/delete0"), 51)
	require.NoError(t, err)

	r, err = r.Insert([]byte("DELETE"), []byte("/delete1"), 52)
	require.NoError(t, err)

	r, err = r.Insert([]byte("DELETE"), []byte("/delete1/{param}"), 53)
	require.NoError(t, err)

	r, err = r.Insert([]byte("CONNECT"), []byte("/connect0"), 61)
	require.NoError(t, err)

	r, err = r.Insert([]byte("CONNECT"), []byte("/connect1"), 62)
	require.NoError(t, err)

	r, err = r.Insert([]byte("CONNECT"), []byte("/connect1/{param}"), 63)
	require.NoError(t, err)

	r, err = r.Insert([]byte("OPTIONS"), []byte("/options0"), 71)
	require.NoError(t, err)

	r, err = r.Insert([]byte("OPTIONS"), []byte("/options1"), 72)
	require.NoError(t, err)

	r, err = r.Insert([]byte("OPTIONS"), []byte("/options1/{param}"), 73)
	require.NoError(t, err)

	r, err = r.Insert([]byte("TRACE"), []byte("/trace0"), 81)
	require.NoError(t, err)

	r, err = r.Insert([]byte("TRACE"), []byte("/trace1"), 82)
	require.NoError(t, err)

	r, err = r.Insert([]byte("TRACE"), []byte("/trace1/{param}"), 83)
	require.NoError(t, err)

	type test struct {
		result      error
		foundParams map[string]interface{}
		method      string
		path        string
		handlerID   uint64
	}

	tests := []test{
		{
			method:      "GET",
			path:        "/get0",
			result:      nil,
			handlerID:   1,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "GET",
			path:        "/get1",
			result:      nil,
			handlerID:   2,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "GET",
			path:      "/get1/foo",
			result:    nil,
			handlerID: 3,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "GET",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "HEAD",
			path:        "/head0",
			result:      nil,
			handlerID:   11,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "HEAD",
			path:        "/head1",
			result:      nil,
			handlerID:   12,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "HEAD",
			path:      "/head1/foo",
			result:    nil,
			handlerID: 13,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "HEAD",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "POST",
			path:        "/post0",
			result:      nil,
			handlerID:   21,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "POST",
			path:        "/post1",
			result:      nil,
			handlerID:   22,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "POST",
			path:      "/post1/foo",
			result:    nil,
			handlerID: 23,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "POST",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "PUT",
			path:        "/put0",
			result:      nil,
			handlerID:   31,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "PUT",
			path:        "/put1",
			result:      nil,
			handlerID:   32,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "PUT",
			path:      "/put1/foo",
			result:    nil,
			handlerID: 33,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "PUT",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "PATCH",
			path:        "/patch0",
			result:      nil,
			handlerID:   41,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "PATCH",
			path:        "/patch1",
			result:      nil,
			handlerID:   42,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "PATCH",
			path:      "/patch1/foo",
			result:    nil,
			handlerID: 43,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "PATCH",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "DELETE",
			path:        "/delete0",
			result:      nil,
			handlerID:   51,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "DELETE",
			path:        "/delete1",
			result:      nil,
			handlerID:   52,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "DELETE",
			path:      "/delete1/foo",
			result:    nil,
			handlerID: 53,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "DELETE",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "CONNECT",
			path:        "/connect0",
			result:      nil,
			handlerID:   61,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "CONNECT",
			path:        "/connect1",
			result:      nil,
			handlerID:   62,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "CONNECT",
			path:      "/connect1/foo",
			result:    nil,
			handlerID: 63,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "CONNECT",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "OPTIONS",
			path:        "/options0",
			result:      nil,
			handlerID:   71,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "OPTIONS",
			path:        "/options1",
			result:      nil,
			handlerID:   72,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "OPTIONS",
			path:      "/options1/foo",
			result:    nil,
			handlerID: 73,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "OPTIONS",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "TRACE",
			path:        "/trace0",
			result:      nil,
			handlerID:   81,
			foundParams: map[string]interface{}{},
		},
		{
			method:      "TRACE",
			path:        "/trace1",
			result:      nil,
			handlerID:   82,
			foundParams: map[string]interface{}{},
		},
		{
			method:    "TRACE",
			path:      "/trace1/foo",
			result:    nil,
			handlerID: 83,
			foundParams: map[string]interface{}{
				"param": []byte("foo"),
			},
		},
		{
			method:      "TRACE",
			path:        "/not/found",
			result:      httprouter.ErrPathNotFound,
			foundParams: map[string]interface{}{},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d#", i), func(t *testing.T) {
			params := make(map[string]interface{})

			handlerID, found := r.Route([]byte(tt.method), []byte(tt.path), func(n string, v interface{}) {
				params[n] = v
			})

			assert.Equal(t, tt.result, found)
			assert.Equal(t, tt.foundParams, params)
			assert.Equal(t, int64(tt.handlerID), int64(handlerID))
		})
	}
}

func TestDelete(t *testing.T) {
	var err error

	r := httprouter.Router{}

	r, err = r.Insert([]byte("GET"), []byte("/foo"), 1)
	require.NoError(t, err)

	r, err = r.Insert([]byte("POST"), []byte("/foo/{bar}"), 2)
	require.NoError(t, err)

	// guard
	_, err = r.Route([]byte("GET"), []byte("/foo"), nil)
	require.NoError(t, err)

	// guard
	_, err = r.Route([]byte("POST"), []byte("/foo/{bar}"), nil)
	require.NoError(t, err)

	r, err = r.Delete([]byte("GET"), []byte("/foo"))
	require.NoError(t, err)

	r, err = r.Delete([]byte("POST"), []byte("/foo/{bar}"))
	require.NoError(t, err)

	_, err = r.Route([]byte("GET"), []byte("/foo"), nil)
	require.EqualError(t, err, httprouter.ErrPathNotFound.Error())

	_, err = r.Route([]byte("POST"), []byte("/foo/{bar}"), nil)
	require.EqualError(t, err, httprouter.ErrPathNotFound.Error())
}
