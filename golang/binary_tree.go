package main

import (
	"fmt"
	"github.com/apache/arrow/go/arrow/memory"
)
import "github.com/gomem/gomem/pkg/dataframe"

type Integer int
type Float float64

type QItem interface {
	Less(QItem) bool
}

func (a Integer) Less(b QItem) bool {
	if a < b.(Integer) {
		return true
	}
	return false
}

func (a Float) Less(b QItem) bool {
	if a < b.(Float) {
		return true
	}
	return false
}

type BinaryNode struct {
	value QItem
	children []*BinaryNode
}

func MakeBinaryNode(val QItem) *BinaryNode {
	return & BinaryNode{
		value:val,
		children:make([]*BinaryNode, 2),
	}
}

type BinaryTree struct {
	root *BinaryNode
	height int
}

func MakeBinaryTree() *BinaryTree {
	return &BinaryTree{
		root: nil,
		height: 0,
	}
}



func recursiveInsert(parent *BinaryNode, value QItem) *BinaryNode {
	if parent == nil {
		return MakeBinaryNode(value)
	} else {
		whichIndex := 0
		if value.Less(parent.value) {
			whichIndex = 0
		} else {
			whichIndex = 1
		}
		child := recursiveInsert(parent.children[whichIndex], value)
		if parent.children[whichIndex] == nil {
			parent.children[whichIndex] = child
		}
		return parent
	}
}

func (self *BinaryTree)Insert(value QItem) {
	if self.root == nil {
		self.root = MakeBinaryNode(value)
	} else {
		recursiveInsert(self.root, value)
	}
}

func (self *BinaryTree) Print() {
	if self.root != nil {
		recursivePrint(self.root, 0)
	}
}

func recursivePrint(parent *BinaryNode, level int) {
	if parent != nil {
		recursivePrint(parent.children[1], level + 1)
		for i := 0; i < level; i++ {
			fmt.Print("....")
		}
		fmt.Println(parent.value)
		recursivePrint(parent.children[0], level + 1)
	}
}


func main() {
	bt := MakeBinaryTree()
	bt.Insert(Float(5.0))
	bt.Insert(Float(3))
	bt.Insert(Float(7))
	bt.Insert(Float(2))
	bt.Insert(Float(4))
	bt.Insert(Float(9))
	bt.Insert(Float(6))
	bt.Print()


	pool := memory.NewGoAllocator()
	df, _ := dataframe.NewDataFrameFromMem(pool, dataframe.Dict{
		"col1": []int32{1, 2, 3, 4, 5},
		"col2": []float64{1.1, 2.2, 3.3, 4.4, 5},
		"col3": []string{"foo", "bar", "ping", "", "pong"},
		"col4": []interface{}{2, 4, 6, nil, 8},
	})
	fmt.Println(df)
	defer df.Release()
}
