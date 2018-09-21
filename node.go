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
	// 一下三个if次序不能乱
	// 因为他们从上到下，依次会造成新状态

	// 右边是红色，左旋
	if tmp.right.isRed() {
		tmp = n.leftRotate()
	}
	// 左边是连续红色，右旋
	if tmp.left.isRed() && tmp.left.left.isRed() {
		tmp = tmp.rightRotate()
	}
	// 左右是红色，变色
	if tmp.left.isRed() && tmp.right.isRed() {
		tmp = tmp.colorFlip()
	}
	return tmp
}

func (n *node) moveRed2Right() *node {
	tmp := n
	// 右存在
	// 右为黑
	// 右左为黑
	// 不判断右右，因为我们是默认LLRB
	if tmp.right != nil && !tmp.right.isRed() && !tmp.right.left.isRed() {
		tmp = tmp.colorFlip()
		if tmp.left != nil && tmp.left.left.isRed() {
			tmp = tmp.rightRotate()
			tmp = tmp.colorFlip()
		}
	}
	return tmp
}

func (n *node) deleteMax() *node {
	tmp := n
	if tmp.right == nil {
		return nil
	}
	if tmp.left.isRed() {
		tmp = tmp.rightRotate()
	}
	tmp.right = tmp.moveRed2Right()
	tmp.right = tmp.right.deleteMax()
	return tmp.fixUp()
}

func (n *node) moveRed2Left() *node {
	tmp := n
	if n.left != nil && !n.left.isRed() && !n.left.left.isRed() {
		tmp := tmp.colorFlip()
		if tmp.right.left.isRed() {
			tmp.right = tmp.right.rightRotate()
			tmp = tmp.leftRotate()
			tmp = tmp.colorFlip()
		}
	}
	return tmp
}

func (n *node) deleteMin() *node {
	tmp := n
	if tmp.left == nil {
		return nil
	}
	// 此处不进行右红判断
	// 因为我们是默认LLRB
	tmp = tmp.moveRed2Left()
	tmp.left = tmp.left.deleteMin()
	return tmp.fixUp()
}

func (n *node) min() *node {
	tmp := n
	for tmp.left != nil {
		tmp = tmp.left
	}
	return tmp
}

// keys 获取节点下所有keys
func (n *node) keys() []interface{} {
	var out []interface{}
	if n.left != nil {
		left := n.left.keys()
		out = append(out, left...)
	}
	if n.right != nil {
		right := n.right.keys()
		out = append(out, right...)
	}
	out = append(out, n.k)
	return out
}
