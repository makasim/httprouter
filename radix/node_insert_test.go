package radix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_InsertStatic(main *testing.T) {
	type test struct {
		node       Node
		insertPath string
		insertKey  uint64
		expected   Node
	}

	tests := map[string]test{
		"ToEmpty": {
			node:       Node{},
			insertPath: "/foo",
			insertKey:  1,
			expected: Node{
				path: "/foo",
				key:  1,
			},
		},
		"SplitCurrent": {
			node: Node{
				path: "/foo",
				key:  1,
			},
			insertPath: "/faa",
			insertKey:  2,
			expected: Node{
				path: "/f",
				key:  0,
				children: []Node{
					{path: "oo", key: 1},
					{path: "aa", key: 2},
				},
			},
		},
		"ShorterAsCurrent": {
			node:       Node{path: "/foo", key: 1},
			insertPath: "/fo",
			insertKey:  2,
			expected: Node{
				path: "/fo",
				key:  2,
				children: []Node{
					{path: "o", key: 1},
				},
			},
		},
		"SameAsCurrent": {
			node:       Node{path: "/foo", key: 1},
			insertPath: "/foo",
			insertKey:  1,
			expected: Node{
				path: "/foo",
				key:  1,
			},
		},
		"NewChild": {
			node: Node{
				path: "/",
				children: []Node{
					{path: "foo", key: 1},
					{path: "bar", key: 2},
				},
			},
			insertPath: "/ololo",
			insertKey:  3,
			expected: Node{
				path: "/",
				children: []Node{
					{path: "foo", key: 1},
					{path: "bar", key: 2},
					{path: "ololo", key: 3},
				},
			},
		},
		"SplitChild": {
			node: Node{
				path: "/",
				key:  0,
				children: []Node{
					{
						path: "foo",
						key:  1,
					},
				},
			},
			insertPath: "/faa",
			insertKey:  2,
			expected: Node{
				path: "/",
				key:  0,
				children: []Node{
					{
						path: "f",
						key:  0,
						children: []Node{
							{
								path: "oo",
								key:  1,
							},
							{
								path: "aa",
								key:  2,
							},
						},
					},
				},
			},
		},
		"ShorterChild": {
			node: Node{
				path: "/",
				children: []Node{
					{path: "foo", key: 1},
				},
			},
			insertPath: "/fo",
			insertKey:  3,
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "fo",
						key:  3,
						children: []Node{
							{
								path: "o",
								key:  1,
							},
						},
					},
				},
			},
		},
		"SameChild": {
			node: Node{
				path: "/",
				children: []Node{
					{path: "foo", key: 1},
				},
			},
			insertPath: "/foo",
			insertKey:  1,
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "foo",
						key:  1,
					},
				},
			},
		},
		"Bug1": {
			node: Node{
				path: "insert/10",
				children: []Node{
					{
						path: "10",
						key:  1010,
					},
					{
						path: "0",
						children: []Node{
							{path: "9", key: 1009},
						},
					},
				},
			},
			insertPath: "insert/1011",
			insertKey:  1011,
			expected: Node{
				path: "insert/10",
				children: []Node{
					{
						path: "1",
						children: []Node{
							{path: "0", key: 1010},
							{path: "1", key: 1011},
						},
					},
					{
						path: "0",
						children: []Node{
							{path: "9", key: 1009},
						},
					},
				},
			},
		},
		"Bug2": {
			node: Node{
				path: "insert/10",
				children: []Node{
					{
						path: "1",
						children: []Node{
							{path: "0", key: 1010},
							{path: "2", key: 1012},
						},
					},
					{
						path: "00",
						key:  1000,
					},
				},
			},
			insertPath: "insert/1011",
			insertKey:  1011,
			expected: Node{
				path: "insert/10",
				children: []Node{
					{
						path: "1",
						children: []Node{
							{path: "0", key: 1010},
							{path: "2", key: 1012},
							{path: "1", key: 1011},
						},
					},
					{
						path: "00",
						key:  1000,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		main.Run(name, func(t *testing.T) {
			actual := tt.node.Insert(tt.insertPath, tt.insertKey)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestNode_InsertDynamic(main *testing.T) {
	type test struct {
		node       Node
		insertPath string
		insertKey  uint64
		expected   Node
	}

	tests := map[string]test{
		"ToEmpty0": {
			node:       Node{},
			insertPath: "/{foo}",
			insertKey:  1,
			expected: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
				},
			},
		},
		"ToEmpty1": {
			node:       Node{},
			insertPath: "/foo/{bar}",
			insertKey:  1,
			expected: Node{
				path: "/foo/",
				children: []Node{
					{
						kind: param,
						path: "{bar}",
						key:  1,
					},
				},
			},
		},
		"ToEmpty2": {
			node:       Node{},
			insertPath: "/foo/{bar}/baz",
			insertKey:  1,
			expected: Node{
				path: "/foo/",
				children: []Node{
					{
						kind: param,
						path: "{bar}",
						children: []Node{
							{
								path: "/baz",
								key:  1,
							},
						},
					},
				},
			},
		},
		"SplitCurrent": {
			node: Node{
				path: "/foo/bar",
				key:  1,
			},
			insertPath: "/foo/{name}",
			insertKey:  2,
			expected: Node{
				path: "/foo/",
				key:  0,
				children: []Node{
					{
						kind: param,
						path: "{name}",
						key:  2,
					},
					{
						path: "bar",
						key:  1,
					},
				},
			},
		},
		"SplitCurrentParam": {
			node: Node{
				path: "/foo/",
				children: []Node{
					{
						kind: param,
						path: "{bar}",
						key:  1,
					},
				},
			},
			insertPath: "/foo/{bar}/name",
			insertKey:  2,
			expected: Node{
				path: "/foo/",
				key:  0,
				children: []Node{
					{
						kind: param,
						path: "{bar}",
						key:  1,
						children: []Node{
							{
								path: "/name",
								key:  2,
							},
						},
					},
				},
			},
		},
		"ShorterAsCurrent": {
			node: Node{
				path: "/",
				key:  0,
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
				},
			},
			insertPath: "/fo",
			insertKey:  2,
			expected: Node{
				path: "/",
				key:  0,
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
					{
						path: "fo",
						key:  2,
					},
				},
			},
		},
		"SameAsCurrent": {
			node: Node{
				path: "/",
				key:  0,
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
				},
			},
			insertPath: "/{foo}",
			insertKey:  1,
			expected: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
				},
			},
		},
		"NewParamChild": {
			node: Node{
				path: "/",
				children: []Node{
					{path: "foo", key: 1},
					{path: "bar", key: 2},
				},
			},
			insertPath: "/{ololo}",
			insertKey:  3,
			expected: Node{
				path: "/",
				children: []Node{
					{kind: param, path: "{ololo}", key: 3},
					{path: "foo", key: 1},
					{path: "bar", key: 2},
				},
			},
		},
		"NewParamChildChild1": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						key:  1,
						children: []Node{
							{path: "bar", key: 2},
						},
					},
				},
			},
			insertPath: "/foo/{ololo}",
			insertKey:  3,
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						key:  1,
						children: []Node{
							{kind: param, path: "{ololo}", key: 3},
							{path: "bar", key: 2},
						},
					},
				},
			},
		},
		"NewParamChildChild2": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						key:  1,
						children: []Node{
							{path: "bar", key: 2},
						},
					},
				},
			},
			insertPath: "/foo/{ololo}/aaa",
			insertKey:  3,
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "foo/",
						key:  1,
						children: []Node{
							{
								kind: param,
								path: "{ololo}",
								children: []Node{
									{
										path: "/aaa",
										key:  3,
									},
								},
							},
							{path: "bar", key: 2},
						},
					},
				},
			},
		},
		"NewParamChildChild3": {
			node: Node{
				path: "/bar/",
				children: []Node{
					{
						path: "1",
						key:  2,
					},
					{
						path: "0",
						key:  1,
					},
				},
			},
			insertPath: "/bar/{param}",
			insertKey:  3,
			expected: Node{
				path: "/bar/",
				children: []Node{
					{
						kind: param,
						path: "{param}",
						key:  3,
					},
					{
						path: "1",
						key:  2,
					},
					{
						path: "0",
						key:  1,
					},
				},
			},
		},
		"Wildcard": {
			node:       Node{},
			insertPath: "/{*foo}",
			insertKey:  1,
			expected: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{*foo}",
						key:  1,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			actual := tt.node.Insert(tt.insertPath, tt.insertKey)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestNode_InsertInvalid(main *testing.T) {
	n := Node{}

	main.Run("EmptyPath", func(t *testing.T) {
		require.PanicsWithValue(t, "insert: path empty", func() {
			n = n.Insert("", 2)
		})
	})

	main.Run("EmptyKey", func(t *testing.T) {
		n := Node{}

		require.PanicsWithValue(t, "insert: key empty", func() {
			n = n.Insert("/path", 0)
		})
	})

	main.Run("NoParamEnd1", func(t *testing.T) {
		n := Node{}

		require.PanicsWithValue(t, "no right bracket: /{param", func() {
			n = n.Insert("/{param", 2)
		})
	})

	main.Run("NoParamEnd2", func(t *testing.T) {
		n := Node{path: "/foo/bar", key: 1}

		require.PanicsWithValue(t, "no right bracket: {param", func() {
			n = n.Insert("/foo/{param", 2)
		})
	})
}

func TestNode_InsertConflict(main *testing.T) {
	type test struct {
		node        Node
		insertPath  string
		insertKey   uint64
		expectedErr error
	}

	tests := map[string]test{
		"EmptyNode": {
			node: Node{
				path:     "",
				key:      1,
				children: nil,
			},
			insertPath:  "/foo",
			insertKey:   2,
			expectedErr: ErrPathAlreadyTaken,
		},
		"SameNode": {
			node: Node{
				path:     "/foo",
				key:      1,
				children: nil,
			},
			insertPath:  "/foo",
			insertKey:   2,
			expectedErr: ErrPathAlreadyTaken,
		},
		"ChildNode": {
			node: Node{
				path: "/foo",
				children: []Node{
					{
						path: "/bar",
						key:  1,
					},
				},
			},
			insertPath:  "/foo/bar",
			insertKey:   2,
			expectedErr: ErrPathAlreadyTaken,
		},
		"ParamNode": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{foo}",
						key:  1,
					},
				},
			},
			insertPath:  "/{foo}",
			insertKey:   2,
			expectedErr: ErrPathAlreadyTaken,
		},
		"ParamConflict": {
			node: Node{
				path: "/foo/",
				children: []Node{
					{
						kind: param,
						path: "{name}",
						key:  1,
					},
				},
			},
			insertPath:  "/foo/{another}",
			insertKey:   2,
			expectedErr: ErrParamNameConflict,
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			require.PanicsWithError(t, tt.expectedErr.Error(), func() {
				tt.node.Insert(tt.insertPath, tt.insertKey)
			})
		})
	}
}
