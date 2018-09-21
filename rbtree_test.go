package rbtree

import (
	"math/rand"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	randTree *RBTree
	verify   *sync.Map
)

func compare(k1, k2 interface{}) bool {
	if ki1, ok := k1.(int); ok {
		if ki2, ok := k2.(int); ok {
			return ki1 < ki2
		}
	}
	return false
}

func TestMain(t *testing.T) {
	verify = &sync.Map{}
	randTree = NewRBTree(compare)
	for index := 0; index < 10000; index++ {
		k := rand.Intn(0xffffffff)
		v := rand.Intn(0xffffffff)
		randTree.Insert(k, v)
		verify.Store(k, v)
	}
}

// leafIsBlack 所有叶节点为黑
func leafIsBlack(n *node) {
	if n == nil {
		// 叶节点
		So(n.isRed(), ShouldBeFalse)
		return
	}
	if n.left != nil {
		leafIsBlack(n.left)
	}
	if n.right != nil {
		leafIsBlack(n.right)
	}
}

// redSonsIsBlack 红节点的儿子都是黑
func redSonsIsBlack(n *node) {
	if n == nil {
		return
	}
	if n.isRed() {
		// 叶节点
		So(n.left.isRed(), ShouldBeFalse)
	}
	redSonsIsBlack(n.left)
	redSonsIsBlack(n.right)
}

// blackIsEqual node到每个节点的路径的黑节点数量相等
var max = -1

func blackIsEqual(n *node, deep int) {
	if !n.isRed() {
		// 黑节点，深度+1
		deep++
		if max < deep {
			max = deep
		}
	}
	if n == nil {
		// 所有叶节点，深度都是max
		So(max, ShouldEqual, deep)
		return
	}
	blackIsEqual(n.left, deep)
	blackIsEqual(n.right, deep)
}

// LLRB 验证
func noRightRed(n *node) {
	if n == nil {
		return
	}
	// 所有叶节点，深度都是max
	So(n.right.isRed(), ShouldBeFalse)
	noRightRed(n.left)
	noRightRed(n.right)
}
func TestInsert(t *testing.T) {
	Convey("Insert", t, func() {
		// 定律1 每个节点红or黑
		Convey("root is black", func() {
			// 定律2 根黑
			So(randTree.root.isRed(), ShouldBeFalse)
		})
		Convey("leaf is black", func() {
			// 定律3 叶节点黑
			leafIsBlack(randTree.root)
		})
		Convey("red sons is black", func() {
			// 定律4 如果节点是红，它的两个子节点都黑
			redSonsIsBlack(randTree.root)
		})
		Convey("black is equal", func() {
			// 定律5 每个节点到其子孙节点的所有路径上黑节点数目相等
			blackIsEqual(randTree.root, 1)
		})
		Convey("no right red", func() {
			// LLRB 验证，不存在右红节点
			noRightRed(randTree.root)
		})
		Convey("value ok", func() {
			// 所有值得存储正确
			verify.Range(func(k, v interface{}) bool {
				v0, ok := randTree.Get(k)
				So(ok, ShouldBeTrue)
				So(v0, ShouldEqual, v)
				return true
			})
		})
	})
}
