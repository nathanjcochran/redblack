package redblack

import (
	"bytes"
	"fmt"
)

type Tree struct {
	root *node
}

type node struct {
	parent *node
	left   *node
	right  *node
	black  bool
	val    int
}

func (n *node) sibling() *node {
	p := n.parent
	if p.left == n {
		return p.right
	}
	return p.left
}

func (n *node) uncle() *node {
	g := n.parent.parent
	if g.left == n.parent {
		return g.right
	}
	return g.left
}

func (n *node) root() *node {
	r := n
	for r.parent != nil {
		r = r.parent
	}
	return r
}

func (n *node) rotateLeft() {
	x := n.right
	n.right = x.left
	x.left = n

	x.parent = n.parent
	n.parent = x

	if n.right != nil {
		n.right.parent = n
	}

	if x.parent != nil {
		if x.parent.left == n {
			x.parent.left = x
		} else {
			x.parent.right = x
		}
	}
}

func (n *node) rotateRight() {
	x := n.left
	n.left = x.right
	x.right = n

	x.parent = n.parent
	n.parent = x

	if n.left != nil {
		n.left.parent = n
	}

	if x.parent != nil {
		if x.parent.left == n {
			x.parent.left = x
		} else {
			x.parent.right = x
		}
	}
}

func (t *Tree) Insert(val int) {
	n := t.insert(val)
	insertRepair(n)
	t.root = n.root()
}

func (t *Tree) insert(val int) *node {
	if t.root == nil {
		t.root = &node{
			val: val,
		}
		return t.root
	}

	cur := t.root
	for {
		if val < cur.val {
			if cur.left == nil {
				cur.left = &node{
					parent: cur,
					val:    val,
				}
				return cur.left
			}
			cur = cur.left
		} else {
			if cur.right == nil {
				cur.right = &node{
					parent: cur,
					val:    val,
				}
				return cur.right
			}
			cur = cur.right
		}
	}
}

func insertRepair(n *node) {
	cur := n
	for {
		if cur.parent == nil {
			// Make the root black
			cur.black = true
			return
		} else if cur.parent.black {
			// Parent is black, do nothing
			return
		} else if u := cur.uncle(); u != nil && !u.black {
			// Make parent and uncle black
			cur.parent.black = true
			u.black = true

			// Make grandparent red
			g := cur.parent.parent
			g.black = false

			// Iteratively repair grandparent
			cur = g
		} else {
			// Ensure cur is on outside of tree
			p := cur.parent
			g := cur.parent.parent
			if p == g.left && p.right == cur {
				p.rotateLeft()
				cur = p
				p = cur.parent
			} else if p == g.right && p.left == cur {
				p.rotateRight()
				cur = p
				p = cur.parent
			}

			// Rotate parent into grandparent position
			if p == g.left {
				g.rotateRight()
			} else {
				g.rotateLeft()
			}

			// Make parent black and grandparent red
			p.black = true
			g.black = false
			return
		}
	}
}

func (t *Tree) String() string {
	buf := &bytes.Buffer{}
	if t.root != nil {
		t.root.string(buf)
	}
	return buf.String()
}

func (n *node) string(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "(")
	if n.left != nil {
		n.left.string(buf)
	}
	if n.parent == nil {
		fmt.Fprintf(buf, "*")
	}
	fmt.Fprintf(buf, "%d", n.val)
	if n.black {
		fmt.Fprintf(buf, "b")
	} else {
		fmt.Fprintf(buf, "r")
	}
	if n.right != nil {
		n.right.string(buf)
	}
	fmt.Fprintf(buf, ")")
}
