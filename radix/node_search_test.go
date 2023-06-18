package radix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_Search(main *testing.T) {
	type test struct {
		node            Node
		searchPath      string
		expectedKey     uint64
		expecctedParams map[string]interface{}
	}

	tests := map[string]test{
		"NoMatchEmptyNode": {
			node:            Node{},
			searchPath:      "/foo",
			expectedKey:     0,
			expecctedParams: map[string]interface{}{},
		},
		"NoMatch": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "bar",
					},
				},
			},
			searchPath:      "/foo",
			expectedKey:     0,
			expecctedParams: map[string]interface{}{},
		},
		"NoKey": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
					},
				},
			},
			searchPath:      "/aa",
			expectedKey:     0,
			expecctedParams: map[string]interface{}{},
		},
		"Match": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  1,
					},
				},
			},
			searchPath:      "/aa",
			expectedKey:     1,
			expecctedParams: map[string]interface{}{},
		},
		"NoMatchChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						children: []Node{
							{
								path: "bar",
							},
						},
					},
				},
			},
			searchPath:      "/foo/foo",
			expectedKey:     0,
			expecctedParams: map[string]interface{}{},
		},
		"NoKeyChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						children: []Node{
							{
								path: "bar",
							},
						},
					},
				},
			},
			searchPath:      "/foo/bar",
			expectedKey:     0,
			expecctedParams: map[string]interface{}{},
		},
		"MatchChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						children: []Node{
							{
								path: "bar",
								key:  1,
							},
						},
					},
				},
			},
			searchPath:      "/foo/bar",
			expectedKey:     1,
			expecctedParams: map[string]interface{}{},
		},
		"MatchParam0": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{foo}",
					},
				},
			},
			searchPath:      "/ololo",
			expectedKey:     1,
			expecctedParams: map[string]interface{}{"foo": []byte("ololo")},
		},
		"MatchParam1": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{foo}",
						children: []Node{
							{
								path: "/",
								children: []Node{
									{
										key:  2,
										path: "bar",
									},
								},
							},
						},
					},
				},
			},
			searchPath:      "/john/bar",
			expectedKey:     2,
			expecctedParams: map[string]interface{}{"foo": []byte("john")},
		},
		"MatchParam2": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{foo}",
						children: []Node{
							{
								path: "/",
								children: []Node{
									{
										kind: param,
										key:  2,
										path: "{bar}",
									},
								},
							},
						},
					},
				},
			},
			searchPath:  "/john/doe",
			expectedKey: 2,
			expecctedParams: map[string]interface{}{
				"bar": []byte("doe"),
				"foo": []byte("john"),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			params := make(map[string]interface{})

			key := tt.node.Search(tt.searchPath, func(n string, v interface{}) {
				params[n] = v
			})

			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, tt.expecctedParams, params)
		})
	}
}

func TestNode_SearchWildcard(main *testing.T) {
	type test struct {
		node            Node
		searchPath      string
		expectedKey     uint64
		expecctedParams map[string]interface{}
	}

	tests := map[string]test{
		"RootParam": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{*foo}",
					},
				},
			},
			searchPath:  "/john/doe/bar",
			expectedKey: 1,
			expecctedParams: map[string]interface{}{
				"*foo": []byte("john/doe/bar"),
			},
		},
		"RootParamWithChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{*foo}",
						children: []Node{
							{
								kind: static,
								key:  3,
								path: "/doe",
							},
						},
					},
				},
			},
			searchPath:  "/john/doe/bar",
			expectedKey: 1,
			expecctedParams: map[string]interface{}{
				"*foo": []byte("john/doe/bar"),
			},
		},

		"SubParam": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						key:  1,
						path: "{*foo}",
						children: []Node{
							{
								kind: static,
								path: "/",
								children: []Node{
									{
										kind: param,
										key:  4,
										path: "{*bar}",
									},
									{
										kind:     static,
										key:      3,
										path:     "doe",
										children: []Node{},
									},
								},
							},
						},
					},
				},
			},
			searchPath:  "/john/doe/bar/baz",
			expectedKey: 4,
			expecctedParams: map[string]interface{}{
				"*foo": []byte("john"),
				"*bar": []byte("doe/bar/baz"),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			params := make(map[string]interface{})

			key := tt.node.Search(tt.searchPath, func(n string, v interface{}) {
				params[n] = v
			})

			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, tt.expecctedParams, params)
		})
	}
}
