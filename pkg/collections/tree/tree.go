package tree

import (
	"fmt"
)

type Node struct {
	Value               any
	Type                string
	Left, Right, Parent *Node
}

func NewNode(val any, nodeType string, parent *Node) *Node {
	return &Node{val, nodeType, nil, nil, parent}
}

type Tree struct {
	Root *Node
}

type Events struct {
	OnNodeStart func(node *Node)
	OnNodeEnd   func(node *Node)
	Mover       func(node *Node, isLeft bool)
}

func InfixHelper(n *Node, fn func(node *Node, isLeft bool), isLeft bool, events *Events) {
	if n == nil {
		return
	}
	mover := events.Mover
	if n.Left != nil {
		if mover != nil {
			mover(n.Left, true)
		}
		InfixHelper(n.Left, fn, true, events)
	}
	if events.OnNodeStart != nil {
		events.OnNodeStart(n)
	}

	fn(n, isLeft)

	if events.OnNodeEnd != nil {
		events.OnNodeEnd(n)
	}
	if n.Right != nil {
		if mover != nil {
			mover(n.Right, false)
		}
		InfixHelper(n.Right, fn, false, events)
	}
}

func (tree *Tree) Infix(fn func(node *Node, isLeft bool), events *Events) {
	InfixHelper(tree.Root, fn, true, events)
}

func recursiveTreePrint(node *Node, level int, dir string) {

	if node == nil {
		return
	}

	//get repr lines
	leftStr, rightStr := fmt.Sprintf("Level %d nil", level), fmt.Sprintf("Level %d nil", level)
	if node.Left != nil {
		leftStr = fmt.Sprintf("Level %d {value: %s, type: %s}", level, node.Left.Value, node.Left.Type)
	}
	if node.Right != nil {
		rightStr = fmt.Sprintf("Level %d{value: %s, type: %s}", level, node.Right.Value, node.Right.Type)
	}

	//print left and right nodes
	fmt.Println(dir, leftStr, "\t", rightStr)

	recursiveTreePrint(node.Left, level+1, "left")
	recursiveTreePrint(node.Right, level+1, "right")

}

func (tree *Tree) Print() {
	rootStr := fmt.Sprintln("(Root) {value:", tree.Root.Value, ", type:", tree.Root.Type, "}")
	fmt.Println(rootStr)
	recursiveTreePrint(tree.Root, 1, "root")
}
