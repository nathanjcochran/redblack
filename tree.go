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

func (n *node) find(val int) *node {
	cur := n
	for cur != nil {
		if val < cur.val {
			cur = cur.left
		} else if val > cur.val {
			cur = cur.right
		} else {
			return cur
		}
	}
	return nil
}

func (t *Tree) Remove(val int) {
	n := t.root.find(val)
	if n == nil {
		return
	}

	x := nodeToRemove(n)
	t.remove(x)
}

func nodeToRemove(n *node) *node {
	if n.left == nil || n.right == nil {
		return n
	}

	pre := n.left
	for pre.right != nil {
		pre = pre.right
	}
	n.val = pre.val

	return pre
}

// Precondition: node has at most one non-leaf child
func (t *Tree) remove(n *node) {
	// Find non-leaf child, or nil
	var c *node
	if n.left == nil {
		c = n.right
	} else {
		c = n.left
	}

	// Black node - if it has a child, it's red
	if n.black {
		if c != nil {
			c.black = true
			c.parent = n.parent
		} else {
			// Black node has no children
			removeRepair(n)
		}
	}

	// Remove node
	p := n.parent
	if p == nil {
		t.root = c
	} else if p.left == n {
		p.left = c
		t.root = p.root()
	} else {
		p.right = c
		t.root = p.root()
	}
}

func removeRepair(n *node) {
	for {
		// Case 1 - node is root
		if n.parent == nil {
			return
		}

		// Case 2 - sibling is red (parent is black)
		if s := n.sibling(); !s.black {
			p := n.parent
			p.black = false
			s.black = true
			if n == p.left {
				p.rotateLeft()
			} else {
				p.rotateRight()
			}
		}

		p := n.parent
		s := n.sibling() // Sibling is black
		if (s.left == nil || s.left.black) &&
			(s.right == nil || s.right.black) {

			if p.black { // Case 3 - all black
				s.black = false
				n = p
				continue
			} else { // Case 4 - red parent, all black children
				p.black = true
				s.black = false
				return
			}
		} else if n == p.left && // Case 5 - right sibling's left child is red
			s.left != nil && !s.left.black &&
			(s.right == nil || s.right.black) {
			s.black = false
			s.left.black = true
			s.rotateRight()
		} else if n == p.right && // Case 5 - left sibling's right child is red
			s.right != nil && !s.right.black &&
			(s.left == nil || s.left.black) {
			s.black = false
			s.right.black = true
			s.rotateLeft()
		}

		p = n.parent
		s = n.sibling()
		if n == p.left { // Case 6 - right sibling's right child is red
			s.right.black = true
			p.rotateLeft()
		} else { // Case 6 = left sibling's left child is red
			s.left.black = true
			p.rotateRight()
		}
		return
	}
}

func (t *Tree) String() string {
	buf := &bytes.Buffer{}
	if t.root != nil {
		t.root.string(buf)
	} else {
		fmt.Fprintf(buf, "()")
	}
	return buf.String()
}

func (n *node) String() string {
	buf := &bytes.Buffer{}
	if n != nil {
		n.string(buf)
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
