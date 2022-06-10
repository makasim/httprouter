package radix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertStatic(main *testing.T) {
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
					{path: "aa", key: 2},
					{path: "oo", key: 1},
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
					{path: "ololo", key: 3},
					{path: "foo", key: 1},
					{path: "bar", key: 2},
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
								path: "aa",
								key:  2,
							},
							{
								path: "oo",
								key:  1,
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
							{path: "1", key: 1011},
							{path: "0", key: 1010},
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
							{path: "1", key: 1011},
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

func TestInsertDynamic(main *testing.T) {
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
						path: "bar",
						key:  1,
					},
					{
						kind: param,
						path: "{name}",
						key:  2,
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
						path: "fo",
						key:  2,
					},
					{
						kind: param,
						path: "{foo}",
						key:  1,
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
					{path: "foo", key: 1},
					{path: "bar", key: 2},
					{kind: param, path: "{ololo}", key: 3},
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
							{path: "bar", key: 2},
							{kind: param, path: "{ololo}", key: 3},
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
							{path: "bar", key: 2},
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
						},
					},
				},
			},
		},
		"NewParamChildChild3": {
			node: Node{
				path: "/bar",
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
				path: "/bar",
				children: []Node{
					{
						path: "/",
						children: []Node{
							{
								kind: param,
								path: "{param}",
								key:  3,
							},
						},
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
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			actual := tt.node.Insert(tt.insertPath, tt.insertKey)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestInsertEmptyPath(t *testing.T) {
	n := Node{}

	require.PanicsWithValue(t, "insert: path empty", func() {
		n = n.Insert("", 2)
	})
}

func TestInsertEmptyKey(t *testing.T) {
	n := Node{}

	require.PanicsWithValue(t, "insert: key empty", func() {
		n = n.Insert("/path", 0)
	})
}

func TestInsertNoParamEnd1(t *testing.T) {
	n := Node{}

	require.PanicsWithValue(t, "no right bracket: /{param", func() {
		n = n.Insert("/{param", 2)
	})
}

func TestInsertNoParamEnd2(t *testing.T) {
	n := Node{path: "/foo/bar", key: 1}

	require.PanicsWithValue(t, "no right bracket: {param", func() {
		n = n.Insert("/foo/{param", 2)
	})
}

func TestInsertConflict(main *testing.T) {
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

func TestDeleteStatic(main *testing.T) {
	type test struct {
		node       Node
		deletePath string
		expected   Node
	}

	tests := map[string]test{
		"NoMatch": {
			node: Node{
				path: "/",
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
			deletePath: "/cc",
			expected: Node{
				path: "/",
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
		"LongerNoMatch": {
			node: Node{
				path: "/",
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
			deletePath: "oooo",
			expected: Node{
				path: "/",
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
		"ShorterNoMatch": {
			node: Node{
				path: "/",
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
			deletePath: "o",
			expected: Node{
				path: "/",
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
		"DeleteFirstChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						path: "aa",
						key:  2,
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
			deletePath: "oo",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  2,
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
		},
		"DeleteSecondChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						path: "aa",
						key:  2,
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
			deletePath: "aa",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
		},
		"DeleteThirdChild": {
			node: Node{
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						path: "aa",
						key:  2,
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
			deletePath: "bb",
			expected: Node{
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
		"DeleteWithChildren": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  1,
						children: []Node{
							{path: "/aa", key: 2},
						},
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
			deletePath: "aa",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						children: []Node{
							{path: "/aa", key: 2},
						},
					},
					{
						path: "bb",
						key:  3,
					},
				},
			},
		},
		"DeleteNextToLast0": {
			node: Node{
				path: "/",
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
			deletePath: "aa",
			expected: Node{
				path: "/oo",
				key:  1,
			},
		},
		"DeleteNextToLast1": {
			node: Node{
				path: "/",
				key:  3,
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
			deletePath: "aa",
			expected: Node{
				path: "/",
				key:  3,
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
				},
			},
		},
		"DeleteFirstChildChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "aa",
								key:  1,
							},
							{
								path: "bb",
								key:  2,
							},
							{
								path: "cc",
								key:  3,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
			deletePath: "aa/aa",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "bb",
								key:  2,
							},
							{
								path: "cc",
								key:  3,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
		},
		"DeleteSecondChildChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "aa",
								key:  1,
							},
							{
								path: "bb",
								key:  2,
							},
							{
								path: "cc",
								key:  3,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
			deletePath: "aa/bb",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "aa",
								key:  1,
							},
							{
								path: "cc",
								key:  3,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
		},
		"DeleteThirdChildChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "aa",
								key:  1,
							},
							{
								path: "bb",
								key:  2,
							},
							{
								path: "cc",
								key:  3,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
			deletePath: "aa/cc",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa/",
						children: []Node{
							{
								path: "aa",
								key:  1,
							},
							{
								path: "bb",
								key:  2,
							},
						},
					},
					{
						path: "bb",
						key:  4,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			n := tt.node.Delete(tt.deletePath)
			require.Equal(t, tt.expected, n)
		})
	}
}

func TestDeleteDynamic(main *testing.T) {
	type test struct {
		node       Node
		deletePath string
		expected   Node
	}

	tests := map[string]test{
		"NoMatch": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "aa",
						key:  2,
					},
				},
			},
			deletePath: "/{cc}",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "aa",
						key:  2,
					},
				},
			},
		},
		"LongerNoMatch": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "{aa}",
						key:  2,
					},
				},
			},
			deletePath: "{aaaaa}",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "{aa}",
						key:  2,
					},
				},
			},
		},
		"ShorterNoMatch": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "{aa}",
						key:  2,
					},
				},
			},
			deletePath: "{a}",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "oo",
						key:  1,
					},
					{
						kind: param,
						path: "{aa}",
						key:  2,
					},
				},
			},
		},
		"DeleteParamNoOtherChildren": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{aa}",
						key:  3,
					},
				},
			},
			deletePath: "{aa}",
			expected: Node{
				path:     "/",
				children: []Node{},
			},
		},
		"DeleteParamWithOneOtherChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  1,
					},
					{
						kind: param,
						path: "{bb}",
						key:  3,
					},
				},
			},
			deletePath: "{bb}",
			expected: Node{
				path: "/aa",
				key:  1,
			},
		},
		"DeleteParamWithOtherChildren": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  1,
					},
					{
						path: "bb",
						key:  2,
					},
					{
						kind: param,
						path: "{cc}",
						key:  3,
					},
				},
			},
			deletePath: "{cc}",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  1,
					},
					{
						path: "bb",
						key:  2,
					},
				},
			},
		},
		"DeleteParamWithSubChild": {
			node: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{aa}",
						key:  1,
						children: []Node{
							{
								path: "/aa",
								key:  2,
							},
						},
					},
				},
			},
			deletePath: "{aa}",
			expected: Node{
				path: "/",
				children: []Node{
					{
						kind: param,
						path: "{aa}",
						children: []Node{
							{
								path: "/aa",
								key:  2,
							},
						},
					},
				},
			},
		},
		"DeleteNextToLast2": {
			node: Node{
				path: "/",
				children: []Node{
					{
						path: "aa",
						key:  2,
					},
					{
						path: "{oo}",
						kind: param,
						key:  1,
					},
				},
			},
			deletePath: "aa",
			expected: Node{
				path: "/",
				children: []Node{
					{
						path: "{oo}",
						kind: param,
						key:  1,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			n := tt.node.Delete(tt.deletePath)
			require.Equal(t, tt.expected, n)
		})
	}
}

func TestSearchStatic(main *testing.T) {
	type test struct {
		node            Node
		searchPath      string
		expectedKey     uint64
		expecctedParams map[string]interface{}
	}

	tests := map[string]test{
		"NoMatchEmptyNode": {
			node:            Node{},
			searchPath:      "foo",
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
			searchPath:      "foo",
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
			searchPath:      "aa",
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
			searchPath:      "aa",
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
			searchPath:      "foo/foo",
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
			searchPath:      "foo/bar",
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
			searchPath:      "foo/bar",
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
			searchPath:      "ololo",
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
			searchPath:      "john/bar",
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
			searchPath:  "john/doe",
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

func TestString(main *testing.T) {
	type test struct {
		n   Node
		exp string
	}

	tests := []test{
		{
			n: Node{},
			exp: `
/`,
		},
		{
			n: Node{path: "/foo", key: 0},
			exp: `
/foo`,
		},
		{
			n: Node{path: "/foo", key: 5},
			exp: `
/foo=5`,
		},
		{
			n: Node{
				path: "/f",
				children: []Node{
					{path: "oo", key: 1},
					{path: "aa", key: 2},
				}},
			exp: `
/f
 ├oo=1
 └aa=2
`,
		},
		{
			n: Node{
				path: "/fo",
				key:  1,
				children: []Node{
					{path: "o", key: 2},
				}},
			exp: `
/fo=1
  └o=2
`,
		},
		{
			n: Node{
				path: "/f",
				children: []Node{
					{
						path: "o",
						children: []Node{
							{path: "a", key: 1},
							{path: "b", key: 2},
							{path: "c", key: 3},
						},
					},
					{
						path: "a",
						children: []Node{
							{path: "a", key: 4},
							{path: "b", key: 5},
							{path: "c", key: 6},
						},
					},
				}},
			exp: `
/f
 ├o
 |├a=1
 |├b=2
 |└c=3
 └a
  ├a=4
  ├b=5
  └c=6
`,
		},
	}

	for i, tt := range tests {
		tt := tt

		main.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			require.Equal(t, tt.exp, "\n"+tt.n.String(), tt.n.String())
		})
	}
}

func TestNodeCount(t *testing.T) {
	n0 := Node{}
	assert.Equal(t, 0, n0.Count())

	var n1 Node
	assert.Equal(t, 0, n1.Count())

	n2 := Node{key: 123}
	assert.Equal(t, 1, n2.Count())

	n3 := Node{
		children: []Node{
			{},
		},
	}
	assert.Equal(t, 0, n3.Count())

	n4 := Node{
		children: []Node{
			{},
			{key: 1},
			{key: 2},
		},
	}
	assert.Equal(t, 2, n4.Count())

	n5 := Node{
		key: 3,
		children: []Node{
			{},
			{key: 1},
			{key: 2},
		},
	}
	assert.Equal(t, 3, n5.Count())

	n6 := Node{
		key: 3,
		children: []Node{
			{},
			{
				key: 1,
				children: []Node{
					{key: 4},
				},
			},
			{key: 2},
		},
	}
	assert.Equal(t, 4, n6.Count())
}
