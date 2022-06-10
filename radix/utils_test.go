package radix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLongestCommonPrefix(main *testing.T) {
	type test struct {
		a   string
		b   string
		exp int
	}

	tests := map[string]test{
		"NoCommon":    {a: "foo", b: "bar", exp: 0},
		"Common1":     {a: "foo", b: "foa", exp: 2},
		"Common2":     {a: "foooooooo", b: "foooooaaaaaaaaaaaaa", exp: 6},
		"NoCommonUTF": {a: "ααα", b: "ββββ", exp: 0},
		"CommonUTF1":  {a: "ααα", b: "αααββββ", exp: 6},
		"Param1":      {a: "/foo/{bar}", b: "/foo/bar", exp: 5},
		"Param2":      {a: "{bar}", b: "{foo}", exp: 0},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			require.Equal(t, tt.exp, longestCommonPrefix(tt.a, tt.b))
		})
	}
}

func TestFindParamEnd(main *testing.T) {
	type test struct {
		a   string
		exp int
	}

	tests := map[string]test{
		"NoParam":          {a: "foo", exp: -1},
		"NoEnd":            {a: "{foo", exp: -1},
		"NoName":           {a: "{}", exp: -1},
		"NoEndBeforeSlash": {a: "{foo/oo}", exp: -1},
		"Param1":           {a: "{foo}", exp: 5},
		"Param2":           {a: "{foo}/bar", exp: 5},
		"Param3":           {a: "{foo}/{bar}", exp: 5},
		"ParamUTF1":        {a: "{ααα}", exp: 8},
		"ParamUTF2":        {a: "{ααα}/βββ", exp: 8},
		"ParamUTF3":        {a: "{ααα}/{βββ}", exp: 8},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			i := findParamEnd(tt.a)
			require.Equal(t, tt.exp, i)
		})
	}
}

func TestFindParamStart(main *testing.T) {
	type test struct {
		a   string
		exp int
	}

	tests := map[string]test{
		"NoParam":      {a: "foo", exp: -1},
		"NoEnd":        {a: "{foo", exp: 0},
		"NoName":       {a: "{}", exp: 0},
		"AfterStatic":  {a: "foo{bar}", exp: 3},
		"AtPathEnd":    {a: "/f/o/o/{foo}", exp: 7},
		"InPathMiddle": {a: "/f/{foo}/o/o", exp: 3},
		"ParamUTF1":    {a: "{ααα}", exp: 0},
		"ParamUTF2":    {a: "ααα/{ααα}/βββ", exp: 7},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			i := findParamStart(tt.a)
			require.Equal(t, tt.exp, i)
		})
	}
}

func TestFindSlashOrEnd(main *testing.T) {
	type test struct {
		a    string
		rest string
		exp  int
	}

	tests := map[string]test{
		"NoSlash0": {a: "", rest: "", exp: 0},
		"NoSlash1": {a: "foo", rest: "", exp: 3},
		"NoSlash2": {a: "ααα", rest: "", exp: 6},
		"Slash0":   {a: "foo/", rest: "/", exp: 3},
		"Slash1":   {a: "fo/bar", rest: "/bar", exp: 2},
		"Slash2":   {a: "ααα/", rest: "/", exp: 6},
		"Slash3":   {a: "αα/βββ", rest: "/βββ", exp: 4},
	}

	for name, tt := range tests {
		tt := tt

		main.Run(name, func(t *testing.T) {
			i := findSlashOrEnd(tt.a)
			assert.Equal(t, tt.rest, tt.a[i:])
			require.Equal(t, tt.exp, i)
		})
	}
}
