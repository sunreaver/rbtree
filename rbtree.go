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
	r.root, _ = r.insert(r.root, k, v)
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
func (r *RBTree) ShowTree(graphName string) string {
	return fmt.Sprintf("digraph %s {\n%s\n}", graphName, r.root.showDotFormat())
}

// insert 插入
func (r *RBTree) insert(insertAt *node, k, v interface{}) (n *node, newNode bool) {
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

func (r *RBTree) equal(k1, k2 interface{}) bool {
	return !r.less(k1, k2) && !r.less(k2, k1)
}

// Remove 删除k键
func (r *RBTree) Remove(k interface{}) (ok bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.root, ok = r.delete(r.root, k)
	r.root.red = false
	return
}

func (r *RBTree) delete(delNode *node, k interface{}) (new *node, ok bool) {
	// 核心，保证每个被检查的节点都为红
	if delNode == nil {
		return nil, false
	} else if r.less(k, delNode.k) {
		// 左删除
		delNode = delNode.moveRed2Left()
		delNode.left, ok = r.delete(delNode.left, k)
	} else {
		if r.equal(delNode.k, k) {
			// 删当前
			if delNode.right == nil {
				// 直接删
				return delNode.left, true
			}
			// 走右删除逻辑
			// 根据情况，把当前delNode变红
			if delNode.left.isRed() {
				delNode = delNode.rightRotate()
			} else {
				delNode = delNode.moveRed2Right()
			}
			// 删除delNode.right
			smallest := delNode.right.min()
			delNode.k = smallest.k
			delNode.v = smallest.v
			// 问题改为了：删除右节点的最小节点
			delNode.right = delNode.right.deleteMin()
			ok = true
		} else {
			// 大于、右删除
			// 根据情况，把当前delNode变红
			if delNode.left.isRed() {
				delNode = delNode.rightRotate()
			} else {
				delNode = delNode.moveRed2Right()
			}
			delNode.right, ok = r.delete(delNode.right, k)
		}
	}
	if ok {
		// 有删除，才重新平衡
		return delNode.fixUp(), true
	}
	return delNode, false
}

// Keys 获取rb树所有key
func (r *RBTree) Keys() []interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.root.keys()
}
