package radix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_DeleteStatic(main *testing.T) {
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

func TestNode_DeleteDynamic(main *testing.T) {
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
