## LLRB红树左旋红黑树实现

### 示例

```golang
package main

import (
	"fmt"
	"math/rand"

	"github.com/sunreaver/rbtree"
)

func main() {
	tree := rbtree.NewRBTree(compare)
	for i := 0; i < 10; i++ {
		n := rand.Intn(1000)
		tree.Insert(n, i+100)
	}

	fmt.Println(tree.ShowTree("a"))
}

func compare(k1, k2 interface{}) bool {
	if ki1, ok := k1.(int); ok {
		if ki2, ok := k2.(int); ok {
			return ki1 < ki2
		}
	}
	return false
}
```

运行以上示例 `go run main.go > out.dot`

完成后即可生成 `out.dot` 文件

```
digraph a {
425 -> 81 [color = "black", label = "L"]
425 -> 847 [color = "black", label = "R"]
81 -> 59 [color = "black", label = "L"]
81 -> 318 [color = "black", label = "R"]
318 -> 300 [color = "red", label = "L"]
847 -> 540 [color = "black", label = "L"]
847 -> 887 [color = "black", label = "R"]
540 -> 456 [color = "red", label = "L"]
}
```

使用 `dot` 工具生成可视图片 `dot -Tpng ./out.dot > ./out.png`


