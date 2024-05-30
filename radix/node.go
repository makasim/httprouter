package radix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/savsgio/gotils"
)

var ErrPathAlreadyTaken = fmt.Errorf("path already taken")
var ErrParamNameConflict = fmt.Errorf("param name conflict")

type kind uint8

const (
	static kind = iota
	param
)

type Node struct {
	path     string
	children []Node
	key      uint64
	kind     kind
}

func (n Node) Insert(path string, key uint64) Node {
	if path == "" {
		panic("insert: path empty")
	}
	if key < 1 {
		panic("insert: key empty")
	}

	if n.path == "" {
		if paramStart := findParamStart(path); paramStart != -1 {
			paramEnd := findParamEnd(path[paramStart:])
			if paramEnd == -1 {
				panic(fmt.Sprintf("no right bracket: %v", path))
			}

			n.path = path[:paramStart]
			path = path[paramStart:]

			child := Node{
				kind: param,
				path: path[:paramEnd],
			}

			path = path[paramEnd:]
			if path != "" {
				child = child.Insert(path, key)
			} else {
				child.key = key
			}

			n = n.setParamNode(child)
			return n
		}

		if n.key > 0 && n.key != key {
			panic(ErrPathAlreadyTaken)
		}

		n.path = path
		n.key = key
		return n
	}

	if n.path == path {
		if n.key > 0 && n.key != key {
			panic(ErrPathAlreadyTaken)
		}

		n.key = key
		return n
	}

	i := longestCommonPrefix(path, n.path)
	if i > 0 {
		if len(n.path) > i {
			n = n.split(i)
		}

		path = path[i:]
		if path == "" {
			n.key = key
			return n
		}

		if findParamStart(path) == 0 {
			pn, _ := n.paramNode()
			n = n.setParamNode(updateParamNode(pn, path, key))
			return n
		}
	}

	for i, child := range n.children {
		prefix := longestCommonPrefix(path, child.path)
		if prefix == 0 {
			continue
		}

		if len(child.path) > prefix {
			child = child.split(prefix)
			child = child.Insert(path, key)
			n.children[i] = child
			return n
		}

		if len(path) > prefix {
			n.children[i] = child.Insert(path, key)
			return n
		}

		if child.key > 0 && child.key != key {
			panic(ErrPathAlreadyTaken)
		}

		child.key = key
		n.children[i] = child
		return n
	}

	if start := findParamStart(path); start >= 0 {
		if start == 0 {
			n = n.setParamNode(updateParamNode(Node{kind: param}, path, key))
			return n
		}

		if len(n.children) > 0 {
			n.children = append(n.children, Node{
				path: path[:start],
				children: []Node{
					updateParamNode(Node{kind: param}, path[start:], key),
				},
			})
		} else {
			n.children = []Node{{
				path: path[:start],
				children: []Node{
					updateParamNode(Node{kind: param}, path[start:], key),
				},
			}}
		}
		return n
	}

	n.children = append([]Node{
		{
			path:     path,
			key:      key,
			children: nil,
		},
	}, n.children...)
	return n
}

func (n Node) Delete(path string) Node {
	removeChild := -1

loop:
	for i, child := range n.children {
		switch {
		case len(path) < len(child.path):
			continue
		case path == child.path && len(child.children) > 0:
			n.children[i].key = 0
			break loop
		case path == child.path && len(child.children) == 0:
			removeChild = i
			break loop
		case path[:len(child.path)] == child.path:
			n.children[i] = child.Delete(path[len(child.path):])
			break loop
		}
	}

	if removeChild >= 0 {
		copy(n.children[removeChild:], n.children[removeChild+1:])
		n.children[len(n.children)-1] = Node{}
		n.children = n.children[:len(n.children)-1]
	}

	if n.key == 0 && len(n.children) == 1 && n.children[0].kind != param {
		n.path += n.children[0].path
		n.key = n.children[0].key

		children := n.children[0].children
		n.children[0] = Node{}
		n.children = n.children[:0]
		n.children = children
	}

	return n
}

func (n *Node) Search(path string, kv func(n string, v interface{})) uint64 {
	var n1 *Node

	switch n.kind {
	case static:
		if len(path) > len(n.path) {
			if len(n.children) == 0 || n.path != path[:len(n.path)] {
				return 0
			}

			i := 0
			l := len(n.children)
			hasChildParam := false
			if n.children[0].kind == param {
				i = 1
				hasChildParam = true
			}

			path = path[len(n.path):]
			for ; i < l; i++ {
				n1 = &n.children[i]
				if path[0] == n1.path[0] {
					if key := n1.Search(path, kv); key > 0 {
						return key
					}
					break
				}
			}

			if hasChildParam {
				return n.children[0].Search(path, kv)
			}

			return 0
		} else if n.path == path {
			return n.key
		}

		return 0
	case param:
		i := findSlashOrEnd(path)
		switch {
		case i > 0:
			pn := n.paramName()
			originI := i
			originPath := path
			path = path[i:]

			if len(path) == 0 {
				kv(pn, gotils.S2B(originPath[:originI]))
				return n.key
			}

			if len(n.children) == 0 {
				// wildcard
				if pn[0] == '*' {
					kv(pn, gotils.S2B(originPath))
					return n.key
				}

				kv(pn, gotils.S2B(originPath[:originI]))
				return 0
			}

			i := 0
			l := len(n.children)
			hasChildParam := false
			if n.children[0].kind == param {
				i = 1
				hasChildParam = true
			}

			for ; i < l; i++ {
				n1 := &n.children[i]
				if path[0] == n1.path[0] {
					if key := n1.Search(path, kv); key > 0 {
						kv(pn, gotils.S2B(originPath[:originI]))
						return key
					}
					break
				}
			}

			if hasChildParam {
				kv(pn, gotils.S2B(originPath[:originI]))
				return n.children[0].Search(path, kv)
			}

			// wildcard
			if pn[0] == '*' {
				kv(pn, gotils.S2B(originPath))
				return n.key
			}

			return 0
		default:
			return 0
		}
	default:
		return 0
	}
}

func (n Node) paramName() string {
	if n.kind != param {
		return ""
	}

	return n.path[1 : len(n.path)-1]
}

func (n Node) paramNode() (Node, bool) {
	if len(n.children) == 0 {
		return Node{kind: param}, false
	}

	pn := n.children[0]
	if pn.kind != param {
		return Node{kind: param}, false
	}

	return pn, true
}

func (n Node) setParamNode(pn Node) Node {
	if len(n.children) == 0 {
		n.children = append(n.children, pn)
		return n
	}

	if n.children[0].kind != param {
		n.children = append([]Node{pn}, n.children...)
		return n
	}

	n.children[0] = pn
	return n
}

func (n Node) split(i int) Node {
	rightPath := n.path[:i]
	leftPath := n.path[i:]

	n.path = leftPath

	return Node{
		path:     rightPath,
		children: []Node{n},
	}
}

func (n Node) Count() int {
	cnt := 0
	if n.key != 0 {
		cnt++
	}

	for _, c := range n.children {
		cnt += c.Count()
	}

	return cnt
}

func (n Node) Clone() Node {
	cloneNode := n
	cloneNode.path = n.path
	cloneNode.key = n.key
	cloneNode.kind = n.kind

	if len(n.children) > 0 {
		cloneNode.children = make([]Node, len(n.children))
		for i, child := range n.children {
			cloneNode.children[i] = child.Clone()
		}
	}

	return cloneNode
}

func (n Node) String() string {
	res := "/"
	if n.path != "" {
		res = n.path
	}

	if n.key != 0 {
		res += "=" + strconv.FormatUint(n.key, 10)
	}

	if len(n.children) > 0 {
		res += "\n"
	}

	last := len(n.children) - 1

	prefix := ""

	if len(n.path) > 1 {
		prefix = strings.Repeat(" ", len(n.path)-1)
	}

	for i, child := range n.children {
		if i == last {
			for j, line := range strings.Split(child.String(), "\n") {
				if line == "" {
					continue
				}

				if j == 0 {
					res += prefix + `└` + line + "\n"
				} else {
					res += prefix + " " + line + "\n"
				}
			}
		} else {
			for j, line := range strings.Split(child.String(), "\n") {
				if line == "" {
					continue
				}

				if j == 0 {
					res += prefix + `├` + line + "\n"
				} else {
					res += prefix + `|` + line + "\n"
				}
			}
		}
	}

	return res
}
