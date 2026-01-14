package dag

import (
	"fmt"
	"sort"
)

// Node represents a node in the DAG.
type Node[T comparable] struct {
	data T

	parents  map[*Node[T]]struct{}
	children map[*Node[T]]struct{}
}

// NewNode creates and returns a new Node with the given data.
func NewNode[T comparable](data T) *Node[T] {
	return &Node[T]{
		data:     data,
		parents:  make(map[*Node[T]]struct{}),
		children: make(map[*Node[T]]struct{}),
	}
}

// Data returns the data stored in the node.
func (n *Node[T]) Data() T {
	return n.data
}

// Parents returns the parent nodes of the current node in deterministic order.
func (n *Node[T]) Parents() []*Node[T] {
	parents := make([]*Node[T], 0, len(n.parents))
	for parent := range n.parents {
		parents = append(parents, parent)
	}
	// Sort by string representation of node data for deterministic ordering
	sort.Slice(parents, func(i, j int) bool {
		return fmt.Sprintf("%v", parents[i].data) < fmt.Sprintf("%v", parents[j].data)
	})
	return parents
}

// Children returns the child nodes of the current node in deterministic order.
func (n *Node[T]) Children() []*Node[T] {
	children := make([]*Node[T], 0, len(n.children))
	for child := range n.children {
		children = append(children, child)
	}
	// Sort by string representation of node data for deterministic ordering
	sort.Slice(children, func(i, j int) bool {
		return fmt.Sprintf("%v", children[i].data) < fmt.Sprintf("%v", children[j].data)
	})
	return children
}

func (n *Node[T]) addParent(parent *Node[T]) {
	n.parents[parent] = struct{}{}
}

func (n *Node[T]) addChild(child *Node[T]) {
	n.children[child] = struct{}{}
}

// hasAncestor checks if the current node has the specified ancestor node.
// It performs a depth-first search to determine the relationship.
func (n *Node[T]) hasAncestor(ancestor *Node[T]) bool {
	visited := make(map[*Node[T]]struct{})
	var visit func(node *Node[T]) bool
	visit = func(node *Node[T]) bool {
		if node == ancestor {
			return true
		}
		visited[node] = struct{}{}
		// Use deterministic iteration order
		for _, parent := range node.Parents() {
			if _, seen := visited[parent]; !seen {
				if visit(parent) {
					return true
				}
			}
		}
		return false
	}
	return visit(n)
}
