package main

import (
	"fmt"
)

// BOOKS
//an introduction to programming in by caleb
//inducing go caleb
//the way to bo by ivo
//go in action william kennedy
//the go programming language
//
//youtube
//goo.gl/heBonM
//
//github.com/golang-book/bootcamp-example
//godoc.org
//
//github.com/avelion/awesome-go

//https://www.youtube.com/playlist?list=PLSak_q1UXfPp2VwUQ4ZdUVJdMO6pfi5v_
//https://www.linkedin.com/learning/code-clinic-go/working-with-files

type Node struct {
	order       int
	count       int
	values 		[]int
	children 	[]*Node
}
func MakeNode(order int) *Node {
	return &Node{
		order:order,
		count:0,
		values:make([]int, order + 1),
		children:make([]*Node, order + 2),
	}
}
type Status int

const (
	StatusOk Status = iota
	StatusOverflow
	StatusUnderflow
)

func (self *Node) InsertInto(position int, value int) {
	index := self.count
	for index > position {
		self.values[index] = self.values[index-1]
		self.children[index+1] = self.children[index]
		index--
	}
	//self.children[index+1] = self.children[index]

	self.values[index] = value
	self.count++
}

func (self *Node) Push(value int) {
	self.InsertInto(self.count, value)
}

func (self *Node) RecursiveInsert(value int) Status {
	index := 0
	for index < self.count && self.values[index] < value {
		index++
	}
	if self.children[index] == nil {
		self.InsertInto(index, value)
	} else {
		status := self.children[index].RecursiveInsert(value)
		if status == StatusOverflow {
			SplitNode(self, index)
		}
	}
	if self.count == self.order {
		return StatusOverflow
	}
	return StatusOk
}

func SplitNode(parent *Node, position int) {
	ptr := parent.children[position]
	mid := ptr.count / 2

	midValue := ptr.values[mid]
	left := MakeNode(ptr.order)
	right := MakeNode(ptr.order)

	index := 0
	for ; index < mid; index++ {
		left.children[left.count] = ptr.children[index]
		left.Push(ptr.values[index])
	}
	left.children[left.count] = ptr.children[index]
	index++
	for ; index < ptr.count; index++ {
		right.children[right.count] = ptr.children[index]
		right.Push(ptr.values[index])
	}
	right.children[right.count] = ptr.children[index]

	parent.InsertInto(position, midValue)
	parent.children[position] = left
	parent.children[position+1] = right
}

func SplitRoot(ptr *Node) {

	mid := ptr.count / 2

	midValue := ptr.values[mid]
	left := MakeNode(ptr.order)
	right := MakeNode(ptr.order)

	index := 0
	for ; index < mid; index++ {
		left.children[left.count] = ptr.children[index]
		left.Push(ptr.values[index])
	}
	left.children[left.count] = ptr.children[index]
	index++

	for ; index < ptr.count; index++ {
		right.children[right.count] = ptr.children[index]
		right.Push(ptr.values[index])
	}
	right.children[right.count] = ptr.children[index]

	ptr.values[0] = midValue
	ptr.children[0] = left
	ptr.children[1] = right
	ptr.count = 1
}

type Btree struct {
	root *Node
	height int
}
func MakeBtree(order int) *Btree {
	return &Btree{
		root: MakeNode(order),
		height: 0,
	}
}

func (self *Btree)Insert(value int) {
	status := self.root.RecursiveInsert(value)
	if status == StatusOverflow {
		SplitRoot(self.root)
	}
}
func RecursivePrint(ptr *Node, level int) {
	if ptr != nil {
		index := ptr.count
		for ; index > 0; index-- {
			RecursivePrint(ptr.children[index], level+1)
			for j := 0; j < level * 5 ; j++  {
				fmt.Printf(".")
			}
			fmt.Println(ptr.values[index-1])
		}
		RecursivePrint(ptr.children[index], level+1)
	}
}

func (self *Btree) Print() {
	RecursivePrint(self.root, 0)
}




func main() {

	elements := "qazxswcdevrftgbnyhujmiklop"
	btree := MakeBtree(3)
	for _, item := range elements {
		fmt.Println("inserting: ", int(item))
		btree.Insert(int(item))
		btree.Print()
	}
	//f, err := os.Open("C:/Users/aocsa/go/data/input.csv")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//rdr := csv.NewReader(f)
	//rdr.Comma = '\t'
	//rdr.TrimLeadingSpace = true
	//
	//rows, err := rdr.ReadAll()
	//if err != nil {
	//	panic(err)
	//}
	//
	//for i, row := range rows {
	//	fmt.Println(row)
	//	if i == 1 {
	//		at, _ := strconv.ParseFloat(row[1], 64)
	//		bm, _ := strconv.ParseFloat(row[2], 64)
	//		ws, _ := strconv.ParseFloat(row[3], 64)
	//
	//		fmt.Printf("%T %T %T\n", row[1], row[2], row[7])
	//		fmt.Println(row[1], row[2], row[7])
	//		fmt.Printf("%T %T %T\n", at, bm, ws)
	//
	//		break
	//	}
	//}
}
