package stdrouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/makasim/httprouter/stdrouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter_AddHandler(main *testing.T) {
	main.Run("OK", func(t *testing.T) {
		r := stdrouter.New()

		h1ID := r.AddHandler(stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		}))

		h2ID := r.AddHandler(stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		}))

		h3ID := r.AddHandler(stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		}))

		require.Equal(t, stdrouter.HandlerID(1), h1ID)
		require.Equal(t, stdrouter.HandlerID(2), h2ID)
		require.Equal(t, stdrouter.HandlerID(3), h3ID)

		r.RemoveHandler(h2ID)

		h4ID := r.AddHandler(stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		}))

		require.Equal(t, stdrouter.HandlerID(2), h4ID)

		h5ID := r.AddHandler(stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		}))

		require.Equal(t, stdrouter.HandlerID(4), h5ID)
	})

	main.Run("Nil", func(t *testing.T) {
		r := stdrouter.New()

		require.Panics(t, func() {
			r.AddHandler(nil)
		})
	})
}

func TestRouter_Insert(main *testing.T) {
	main.Run("MethodEmpty", func(t *testing.T) {
		r := stdrouter.New()

		err := r.Add("", "apath", 123)
		require.EqualError(t, err, "method not allowed")
	})

	main.Run("MethodUnsupported", func(t *testing.T) {
		r := stdrouter.New()

		err := r.Add("unsupported", "apath", 123)
		require.EqualError(t, err, "method not allowed")
	})

	main.Run("PathEmpty", func(t *testing.T) {
		r := stdrouter.New()

		err := r.Add("POST", "", 123)
		require.EqualError(t, err, "path empty")
	})

	main.Run("OK", func(t *testing.T) {
		r := stdrouter.New()

		require.NoError(t, r.Add("GET", "/get", 10))
		require.NoError(t, r.Add("GET", "/get/{param}/foo", 11))

		require.NoError(t, r.Add("HEAD", "/head", 20))
		require.NoError(t, r.Add("HEAD", "/head/{param}/foo", 21))

		require.NoError(t, r.Add("POST", "/post", 30))
		require.NoError(t, r.Add("POST", "/post/{param}/foo", 31))

		require.NoError(t, r.Add("PUT", "/put", 40))
		require.NoError(t, r.Add("PUT", "/put/{param}/foo", 41))

		require.NoError(t, r.Add("PATCH", "/patch", 50))
		require.NoError(t, r.Add("PATCH", "/patch/{param}/foo", 51))

		require.NoError(t, r.Add("DELETE", "/delete", 60))
		require.NoError(t, r.Add("DELETE", "/delete/{param}/foo", 61))

		require.NoError(t, r.Add("CONNECT", "/connect", 70))
		require.NoError(t, r.Add("CONNECT", "/connect/{param}/foo", 71))

		require.NoError(t, r.Add("OPTIONS", "/options", 80))
		require.NoError(t, r.Add("OPTIONS", "/options/{param}/foo", 81))

		require.NoError(t, r.Add("TRACE", "/trace", 90))
		require.NoError(t, r.Add("TRACE", "/TRACE/{param}/foo", 91))

		require.NoError(t, r.Add(stdrouter.MethodAny, "/trace", 100))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/TRACE/{param}/foo", 101))

		// wildcard param test
		require.NoError(t, r.Add(stdrouter.MethodAny, "/{*path}", 110))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/ANY/{*param}", 111))
	})
}

func TestRouter_HandleComplexParametrizedRouting(main *testing.T) {
	main.Run("WhenParametrizedRouteAfterSimpleRoutes_OK", func(t *testing.T) {
		r := stdrouter.New()

		// wildcard param test
		require.NoError(t, r.Add(stdrouter.MethodAny, "/bar/0", 1))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/bar/1", 2))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/bar/{param}", 3))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/bar/rpc.{*param}", 4))

		type test struct {
			status    int
			params    map[string]interface{}
			method    string
			path      string
			handlerID uint64
		}

		tests := []test{
			{
				method:    "POST", // can be any method
				path:      "/bar/0",
				status:    http.StatusOK,
				handlerID: 1,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(1),
				},
			},
			{
				method:    "POST", // can be any method
				path:      "/bar/1",
				status:    http.StatusOK,
				handlerID: 2,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(2),
					"path":                        "bar/1",
				},
			},
			{
				method:    "POST", // can be any method
				path:      "/bar/foobar",
				status:    http.StatusOK,
				handlerID: 3,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(3),
					"param":                       "foobar",
				},
			},

			{
				method:    "POST", // can be any method
				path:      "/bar/rpc.v4",
				status:    http.StatusOK,
				handlerID: 4,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(4),
					"param":                       "rpc.v4",
				},
			},
		}

		r.GlobalHandler = stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		})

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%s-%d", tt.method, tt.handlerID), func(t *testing.T) {
				req := &http.Request{}
				req.Method = tt.method
				req.URL = &url.URL{}
				req.URL.Path = tt.path
				rw := &httptest.ResponseRecorder{}

				r.ServeHTTP(rw, req)

				assert.Equal(t, tt.status, rw.Result().StatusCode)
			})
		}
	})
	main.Run("June2024_Bug", func(t *testing.T) {
		r := stdrouter.New()

		// wildcard param test
		require.NoError(t, r.Add(stdrouter.MethodAny, "/api/{*path}", 1))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/api/v1/foo/bar", 2))

		type test struct {
			status    int
			params    map[string]interface{}
			method    string
			path      string
			handlerID uint64
		}

		tests := []test{
			{
				method:    "POST", // can be any method
				path:      "/api/v1/foo/bar",
				status:    http.StatusOK,
				handlerID: 1,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(1),
				},
			},
			{
				method:    "POST", // can be any method
				path:      "/api/something/else",
				status:    http.StatusOK,
				handlerID: 2,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(2),
					"path":                        "something/else",
				},
			},
		}

		r.GlobalHandler = stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		})

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%s-%d", tt.method, tt.handlerID), func(t *testing.T) {
				req := &http.Request{}
				req.Method = tt.method
				req.URL = &url.URL{}
				req.URL.Path = tt.path
				rw := &httptest.ResponseRecorder{}

				r.ServeHTTP(rw, req)

				assert.Equal(t, tt.status, rw.Result().StatusCode)
			})
		}
	})
	main.Run("June2024_Bug_2", func(t *testing.T) {
		r := stdrouter.New()

		// wildcard param test
		require.NoError(t, r.Add(stdrouter.MethodAny, "/api/v1/foo/bar", 1))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/api/{*path}", 2))

		type test struct {
			status    int
			params    map[string]interface{}
			method    string
			path      string
			handlerID uint64
		}

		tests := []test{
			{
				method:    "POST", // can be any method
				path:      "/api/v1/foo/bar",
				status:    http.StatusOK,
				handlerID: 1,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(1),
				},
			},
			{
				method:    "POST", // can be any method
				path:      "/api/something/else",
				status:    http.StatusOK,
				handlerID: 2,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(2),
					"path":                        "something/else",
				},
			},
		}

		r.GlobalHandler = stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		})

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%s-%d", tt.method, tt.handlerID), func(t *testing.T) {
				req := &http.Request{}
				req.Method = tt.method
				req.URL = &url.URL{}
				req.URL.Path = tt.path
				rw := &httptest.ResponseRecorder{}

				r.ServeHTTP(rw, req)

				assert.Equal(t, tt.status, rw.Result().StatusCode)
			})
		}
	})
}

func TestRouter_Handle(main *testing.T) {
	main.Run("Route", func(t *testing.T) {
		r := stdrouter.New()

		require.NoError(t, r.Add("GET", "/get0", 1))
		require.NoError(t, r.Add("GET", "/get1", 2))
		require.NoError(t, r.Add("GET", "/get1/{param}", 3))
		require.NoError(t, r.Add("HEAD", "/head0", 11))
		require.NoError(t, r.Add("HEAD", "/head1", 12))
		require.NoError(t, r.Add("HEAD", "/head1/{param}", 13))
		require.NoError(t, r.Add("POST", "/post0", 21))
		require.NoError(t, r.Add("POST", "/post1", 22))
		require.NoError(t, r.Add("POST", "/post1/{param}", 23))
		require.NoError(t, r.Add("PUT", "/put0", 31))
		require.NoError(t, r.Add("PUT", "/put1", 32))
		require.NoError(t, r.Add("PUT", "/put1/{param}", 33))
		require.NoError(t, r.Add("PATCH", "/patch0", 41))
		require.NoError(t, r.Add("PATCH", "/patch1", 42))
		require.NoError(t, r.Add("PATCH", "/patch1/{param}", 43))
		require.NoError(t, r.Add("DELETE", "/delete0", 51))
		require.NoError(t, r.Add("DELETE", "/delete1", 52))
		require.NoError(t, r.Add("DELETE", "/delete1/{param}", 53))
		require.NoError(t, r.Add("CONNECT", "/connect0", 61))
		require.NoError(t, r.Add("CONNECT", "/connect1", 62))
		require.NoError(t, r.Add("CONNECT", "/connect1/{param}", 63))
		require.NoError(t, r.Add("OPTIONS", "/options0", 71))
		require.NoError(t, r.Add("OPTIONS", "/options1", 72))
		require.NoError(t, r.Add("OPTIONS", "/options1/{param}", 73))
		require.NoError(t, r.Add("TRACE", "/trace0", 81))
		require.NoError(t, r.Add("TRACE", "/trace1", 82))
		require.NoError(t, r.Add("TRACE", "/trace1/{param}", 83))

		// wildcard param test
		require.NoError(t, r.Add(stdrouter.MethodAny, "/ANY/foo/car", 112))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/ANY/foo/bar/{*barparam}", 111))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/ANY/foo/{*fooparam}", 110))
		require.NoError(t, r.Add(stdrouter.MethodAny, "/ANY/{*rootparam}", 113))

		type test struct {
			status    int
			params    map[string]interface{}
			method    string
			path      string
			handlerID uint64
		}

		tests := []test{
			{
				method:    "GET",
				path:      "/get0",
				status:    http.StatusOK,
				handlerID: 1,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(1),
				},
			},
			{
				method:    "GET",
				path:      "/get1",
				status:    http.StatusOK,
				handlerID: 2,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(2),
				},
			},
			{
				method:    "GET",
				path:      "/get1/foo",
				status:    http.StatusOK,
				handlerID: 3,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(3),
					"param":                       "foo",
				},
			},
			{
				method: "GET",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "HEAD",
				path:      "/head0",
				status:    http.StatusOK,
				handlerID: 11,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(11),
				},
			},
			{
				method:    "HEAD",
				path:      "/head1",
				status:    http.StatusOK,
				handlerID: 12,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(12),
				},
			},
			{
				method:    "HEAD",
				path:      "/head1/foo",
				status:    http.StatusOK,
				handlerID: 13,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(13),
					"param":                       "foo",
				},
			},
			{
				method: "HEAD",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "POST",
				path:      "/post0",
				status:    http.StatusOK,
				handlerID: 21,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(21),
				},
			},
			{
				method:    "POST",
				path:      "/post1",
				status:    http.StatusOK,
				handlerID: 22,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(22),
				},
			},
			{
				method:    "POST",
				path:      "/post1/foo",
				status:    http.StatusOK,
				handlerID: 23,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(23),
					"param":                       "foo",
				},
			},
			{
				method: "POST",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "PUT",
				path:      "/put0",
				status:    http.StatusOK,
				handlerID: 31,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(31),
				},
			},
			{
				method:    "PUT",
				path:      "/put1",
				status:    http.StatusOK,
				handlerID: 32,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(32),
				},
			},
			{
				method:    "PUT",
				path:      "/put1/foo",
				status:    http.StatusOK,
				handlerID: 33,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(33),
					"param":                       "foo",
				},
			},
			{
				method: "PUT",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "PATCH",
				path:      "/patch0",
				status:    http.StatusOK,
				handlerID: 41,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(41),
				},
			},
			{
				method:    "PATCH",
				path:      "/patch1",
				status:    http.StatusOK,
				handlerID: 42,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(42),
				},
			},
			{
				method:    "PATCH",
				path:      "/patch1/foo",
				status:    http.StatusOK,
				handlerID: 43,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(43),
					"param":                       "foo",
				},
			},
			{
				method: "PATCH",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "DELETE",
				path:      "/delete0",
				status:    http.StatusOK,
				handlerID: 51,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(51),
				},
			},
			{
				method:    "DELETE",
				path:      "/delete1",
				status:    http.StatusOK,
				handlerID: 52,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(52),
				},
			},
			{
				method:    "DELETE",
				path:      "/delete1/foo",
				status:    http.StatusOK,
				handlerID: 53,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(53),
					"param":                       "foo",
				},
			},
			{
				method: "DELETE",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "CONNECT",
				path:      "/connect0",
				status:    http.StatusOK,
				handlerID: 61,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(61),
				},
			},
			{
				method:    "CONNECT",
				path:      "/connect1",
				status:    http.StatusOK,
				handlerID: 62,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(62),
				},
			},
			{
				method:    "CONNECT",
				path:      "/connect1/foo",
				status:    http.StatusOK,
				handlerID: 63,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(63),
					"param":                       "foo",
				},
			},
			{
				method: "CONNECT",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "OPTIONS",
				path:      "/options0",
				status:    http.StatusOK,
				handlerID: 71,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(71),
				},
			},
			{
				method:    "OPTIONS",
				path:      "/options1",
				status:    http.StatusOK,
				handlerID: 72,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(72),
				},
			},
			{
				method:    "OPTIONS",
				path:      "/options1/foo",
				status:    http.StatusOK,
				handlerID: 73,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(73),
					"param":                       "foo",
				},
			},
			{
				method: "OPTIONS",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "TRACE",
				path:      "/trace0",
				status:    http.StatusOK,
				handlerID: 81,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(81),
				},
			},
			{
				method:    "TRACE",
				path:      "/trace1",
				status:    http.StatusOK,
				handlerID: 82,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(82),
				},
			},
			{
				method:    "TRACE",
				path:      "/trace1/foo",
				status:    http.StatusOK,
				handlerID: 83,
				params: map[string]interface{}{
					stdrouter.HandlerKeyUserValue: uint64(83),
					"param":                       "foo",
				},
			},
			{
				method: "TRACE",
				path:   "/not/found",
				status: http.StatusNotFound,
				params: map[string]interface{}{},
			},
			{
				method:    "POST", // can be any method
				path:      "/ANY/foo/match/by/wildcard",
				status:    http.StatusOK,
				handlerID: 110,
			},
			{
				method:    "POST", // can be any method
				path:      "/ANY/foo/bar/match/by/wildcard",
				status:    http.StatusOK,
				handlerID: 111,
			},
			{
				method:    "POST", // can be any method
				path:      "/ANY/foo/bar/car",
				status:    http.StatusOK,
				handlerID: 112,
			},
			{
				method:    "POST", // can be any method
				path:      "/ANY/root",
				status:    http.StatusOK,
				handlerID: 113,
			},
			{
				method: "UNSUPPORTED",
				path:   "/method/unsupported",
				status: http.StatusMethodNotAllowed,
				params: map[string]interface{}{},
			},
		}

		r.GlobalHandler = stdrouter.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, params stdrouter.Params) {
			rw.WriteHeader(http.StatusOK)
		})

		for _, tt := range tests {
			tt := tt
			t.Run(fmt.Sprintf("%s-%d", tt.method, tt.handlerID), func(t *testing.T) {
				req := &http.Request{}
				req.Method = tt.method
				req.URL = &url.URL{}
				req.URL.Path = tt.path
				rw := &httptest.ResponseRecorder{}

				r.ServeHTTP(rw, req)

				//params := make(map[string]interface{})
				//ctx.VisitUserValues(func(bytes []byte, i interface{}) {
				//	params[string(bytes)] = i
				//})

				assert.Equal(t, tt.status, rw.Result().StatusCode)
				//assert.Equal(t, tt.params, params)
			})
		}
	})

	main.Run("CustomHandlers", func(t *testing.T) {
		r := stdrouter.New()
		r.PageNotFoundHandler = func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte(`custom_not_found_handler`))
			require.NoError(t, err)
		}
		r.MethodNotAllowedHandler = func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, err := writer.Write([]byte(`custom_method_not_allowed_handler`))
			require.NoError(t, err)
		}
		h1ID := r.AddHandler(stdrouter.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, params stdrouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`custom_handler`))
			require.NoError(t, err)
		}))
		r.GlobalHandler = stdrouter.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, params stdrouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`custom_global_handler`))
			require.NoError(t, err)
		})

		require.NoError(t, r.Add("GET", "/get/123", h1ID))
		require.NoError(t, r.Add("GET", "/get/321", 321))

		req := httptest.NewRequest(`GET`, `/get/123`, http.NoBody)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
		require.Equal(t, `custom_handler`, rw.Body.String())

		req = httptest.NewRequest(`GET`, `/get/321`, http.NoBody)
		rw = httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
		require.Equal(t, `custom_global_handler`, rw.Body.String())

		req = httptest.NewRequest(`GET`, `/not/found`, http.NoBody)
		rw = httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusNotFound, rw.Result().StatusCode)
		require.Equal(t, `custom_not_found_handler`, rw.Body.String())

		req = httptest.NewRequest(`UNSUPPORTED`, `/method/not/allowed`, http.NoBody)
		rw = httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusMethodNotAllowed, rw.Result().StatusCode)
		require.Equal(t, `custom_method_not_allowed_handler`, rw.Body.String())
	})

	main.Run("AnyMethod", func(t *testing.T) {
		r := stdrouter.New()
		h1ID := r.AddHandler(stdrouter.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, params stdrouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`get_handler`))
			require.NoError(t, err)
		}))
		h2ID := r.AddHandler(stdrouter.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, params stdrouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`any_handler`))
			require.NoError(t, err)
		}))
		require.NoError(t, r.Add("GET", "/apath", h1ID))
		require.NoError(t, r.Add(`ANY`, "/apath", h2ID))

		req := httptest.NewRequest(`GET`, `/apath`, http.NoBody)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
		require.Equal(t, `get_handler`, rw.Body.String())

		req = httptest.NewRequest(`POST`, `/apath`, http.NoBody)
		rw = httptest.NewRecorder()
		r.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
		require.Equal(t, `any_handler`, rw.Body.String())
	})
}

func TestRouter_Delete(t *testing.T) {
	r := stdrouter.New()
	r.GlobalHandler = stdrouter.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, params stdrouter.Params) {
		writer.WriteHeader(http.StatusOK)
	})

	require.NoError(t, r.Add("GET", "/foo", 1))
	require.NoError(t, r.Add("POST", "/foo/{bar}", 2))

	//guard
	req := httptest.NewRequest("GET", "/foo", http.NoBody)
	rw := httptest.NewRecorder()

	r.ServeHTTP(rw, req)
	require.Equal(t, http.StatusOK, rw.Result().StatusCode)

	require.NoError(t, r.Remove("GET", "/foo"))
	require.NoError(t, r.Remove("POST", "/foo/{bar}"))

	//guard
	req = httptest.NewRequest("GET", "/foo", http.NoBody)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	require.Equal(t, http.StatusNotFound, rw.Result().StatusCode)
}
