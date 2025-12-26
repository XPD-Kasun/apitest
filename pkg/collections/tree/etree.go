package tree

type Operator string

type ETree struct {
	Tree
}

func NewETree(node *Node) *ETree {
	return &ETree{
		Tree: Tree{
			Root: node,
		},
	}
}

func (etree *ETree) InsertNode(node *Node, op Operator) {
	root := etree.Root
	etree.Root = NewNode(op, "op", nil)
	root.Parent = etree.Root
	etree.Root.Left = root
	etree.Root.Right = node
	node.Parent = etree.Root
}
