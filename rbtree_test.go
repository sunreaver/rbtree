package rbtree

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tree *RBTree
	data = []int{0, 2, 26, 33, 15, 59, 78, 47, 89, 94, 90, 81, 107, 137, 133, 156, 159, 147, 106, 189, 199, 194, 205, 211, 241, 237, 266, 258, 202, 162, 287, 318, 300, 336, 353, 356, 376, 355, 324, 408, 420, 413, 387, 274, 429, 445, 433, 451, 447, 463, 485, 466, 456, 503, 511, 510, 538, 528, 495, 541, 552, 561, 546, 577, 598, 623, 605, 563, 643, 703, 694, 718, 721, 705, 631, 540, 425, 737, 783, 746, 828, 843, 831, 790, 878, 887, 891, 940, 888, 953, 996, 957, 947, 847, 728}
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
	tree = NewRBTree(compare)
	rand.Seed(time.Now().UnixNano())
	for _, v := range data {
		tree.Insert(v, rand.Intn(0xffffffff))
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
			So(tree.root.isRed(), ShouldBeFalse)
		})
		Convey("leaf is black", func() {
			// 定律3 叶节点黑
			leafIsBlack(tree.root)
		})
		Convey("red sons is black", func() {
			// 定律4 如果节点是红，它的两个子节点都黑
			redSonsIsBlack(tree.root)
		})
		Convey("black is equal", func() {
			// 定律5 每个节点到其子孙节点的所有路径上黑节点数目相等
			blackIsEqual(tree.root, 1)
		})
		Convey("no right red", func() {
			// LLRB 验证，不存在右红节点
			noRightRed(tree.root)
		})
	})
}
