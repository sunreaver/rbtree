package rbtree

import "fmt"

type node struct {
	k, v        interface{}
	red         bool
	left, right *node
}

func (n *node) say() string {
	if n == nil {
		return ""
	}
	if n.left == nil && n.right == nil {
		return ""
	} else if n.left == nil {
		return fmt.Sprintf("%v -> %v [color = \"%s\", label = \"R\"]\n", n.key(), n.right.key(), n.right.color())
	} else if n.right == nil {
		return fmt.Sprintf("%v -> %v [color = \"%s\", label = \"L\"]\n", n.key(), n.left.key(), n.left.color())
	}
	return fmt.Sprintf("%v -> %v [color = \"%s\", label = \"L\"]\n%v -> %v [color = \"%s\", label = \"R\"]\n", n.key(), n.left.key(), n.left.color(), n.key(), n.right.key(), n.right.color())
}

func (n *node) showDotFormat() string {
	if n == nil {
		return ""
	}
	leftOut := n.left.showDotFormat()
	rightOut := n.right.showDotFormat()
	return fmt.Sprintf("%v%v%v", n.say(), leftOut, rightOut)
}

func (n *node) isRed() bool {
	return n != nil && n.red
}

func (n *node) color() string {
	if n.isRed() {
		return "red"
	}
	return "black"
}

func (n *node) key() interface{} {
	if n == nil {
		return nil
	}
	return n.k
}

func newNode(k, v interface{}) *node {
	return &node{
		k: k,
		v: v,
	}
}

// leftRotate 左旋
func (n *node) leftRotate() *node {
	if n == nil || n.right == nil {
		return nil
	}
	rightChild := n.right
	n.right = rightChild.left
	rightChild.left = n
	rightChild.red = n.red
	n.red = true
	return rightChild
}

// rightRotate 右旋
func (n *node) rightRotate() *node {
	if n == nil || n.left == nil {
		return nil
	}
	leftChild := n.left
	n.left = leftChild.right
	leftChild.right = n
	leftChild.red = n.red
	n.red = true
	return leftChild
}

func (n *node) colorFlip() *node {
	if n == nil {
		return n
	}
	n.red = !n.red
	if n.left != nil {
		n.left.red = !n.left.red
	}
	if n.right != nil {
		n.right.red = !n.right.red
	}
	return n
}

func (n *node) fixUp() *node {
	if n == nil {
		return nil
	}
	tmp := n

	// 左右是红色，变色
	if tmp.left.isRed() && tmp.right.isRed() {
		tmp = tmp.colorFlip()
	}
	// 右边是红色，左旋
	if tmp.right.isRed() {
		tmp = n.leftRotate()
	}
	// 左边是连续红色，右旋
	if tmp.left.isRed() && tmp.left.left.isRed() {
		tmp = tmp.rightRotate()
	}
	return tmp
}
