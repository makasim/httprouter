package radix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_String(main *testing.T) {
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

func TestNode_Count(t *testing.T) {
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
