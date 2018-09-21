package rbtree

import (
	"fmt"
	"sync"
)

// RBTree 红黑树
type RBTree struct {
	lock   sync.RWMutex
	root   *node
	length int
	less   func(k1, k2 interface{}) bool
}

// NewRBTree NewRBTree
func NewRBTree(less func(k1, k2 interface{}) bool) *RBTree {
	if less == nil {
		return nil
	}
	return &RBTree{
		less: less,
		lock: sync.RWMutex{},
	}
}

// Insert Insert
func (r *RBTree) Insert(k, v interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.insert(r.root, k, v)
	r.root.red = false
	r.length++
}

// Len 包含的元素个数
func (r *RBTree) Len() int {
	return r.length
}

// Get 获取
func (r *RBTree) Get(k interface{}) (interface{}, bool) {
	r.lock.RLock()
	r.lock.RUnlock()
	return r.get(r.root, k)
}

func (r *RBTree) get(n *node, k interface{}) (interface{}, bool) {
	if n == nil {
		// not find
		return nil, false
	}
	if r.less(n.k, k) {
		// right
		return r.get(n.right, k)
	} else if r.less(k, n.k) {
		// left
		return r.get(n.left, k)
	}
	return n.v, true
}

// ShowTree 打印key结构
func (r *RBTree) ShowTree() string {
	return fmt.Sprintf("digraph edge_settings {\n%s\n}", r.root.showDotFormat())
}

// insert 插入
func (r *RBTree) insert(insertAt *node, k, v interface{}) (n *node, newNode bool) {
	defer func() {
		if insertAt == r.root && r.root != n {
			r.root = n
		}
	}()

	if insertAt == nil {
		n = &node{
			k:   k,
			v:   v,
			red: true,
		}
		return n, true
	}

	if r.less(k, insertAt.k) {
		insertAt.left, newNode = r.insert(insertAt.left, k, v)
	} else if r.less(insertAt.k, k) {
		insertAt.right, newNode = r.insert(insertAt.right, k, v)
	} else {
		insertAt.v = v
		newNode = false
	}
	if newNode {
		return insertAt.fixUp(), true
	}
	return insertAt, false
}
