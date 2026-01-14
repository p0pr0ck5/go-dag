package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := NewNode("test")
	assert.Equal(t, "test", node.data, "Node data should be 'test'")
	assert.Empty(t, node.parents, "Node should have no parents initially")
	assert.Empty(t, node.children, "Node should have no children initially")
}

func TestData(t *testing.T) {
	node := NewNode(42)
	assert.Equal(t, 42, node.Data(), "Node data should be 42")
}

func TestParents(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")
	nodeC := NewNode("C")

	// Set up parents
	nodeC.addParent(nodeA)
	nodeC.addParent(nodeB)

	parents := nodeC.Parents()
	assert.Len(t, parents, 2, "Node C should have 2 parents")
	assert.Equal(t, "A", parents[0].Data(), "First parent should be A")
	assert.Equal(t, "B", parents[1].Data(), "Second parent should be B")
}

func TestChildren(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")
	nodeC := NewNode("C")

	// Set up children
	nodeA.addChild(nodeB)
	nodeA.addChild(nodeC)

	children := nodeA.Children()
	assert.Len(t, children, 2, "Node A should have 2 children")
	assert.Equal(t, "B", children[0].Data(), "First child should be B")
	assert.Equal(t, "C", children[1].Data(), "Second child should be C")
}

func TestParentsChildrenDeterministicOrder(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")
	nodeC := NewNode("C")

	// Set up parents and children
	nodeC.addParent(nodeA)
	nodeC.addParent(nodeB)
	nodeA.addChild(nodeC)
	nodeB.addChild(nodeC)

	// Test Parents order
	parents := nodeC.Parents()
	assert.Len(t, parents, 2, "Node C should have 2 parents")
	assert.Equal(t, "A", parents[0].Data(), "First parent should be A")
	assert.Equal(t, "B", parents[1].Data(), "Second parent should be B")

	// Test Children order
	childrenA := nodeA.Children()
	childrenB := nodeB.Children()

	assert.Len(t, childrenA, 1, "Node A should have 1 child")
	assert.Equal(t, "C", childrenA[0].Data(), "Child of A should be C")

	assert.Len(t, childrenB, 1, "Node B should have 1 child")
	assert.Equal(t, "C", childrenB[0].Data(), "Child of B should be C")
}

func Test_addParent(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")

	nodeB.addParent(nodeA)

	parents := nodeB.Parents()
	assert.Len(t, parents, 1, "Node B should have 1 parent")
	assert.Equal(t, "A", parents[0].Data(), "Parent of B should be A")
}

func Test_addChild(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")

	nodeA.addChild(nodeB)

	children := nodeA.Children()
	assert.Len(t, children, 1, "Node A should have 1 child")
	assert.Equal(t, "B", children[0].Data(), "Child of A should be B")
}

func Test_hasAncestor(t *testing.T) {
	dag := NewDAG[int]()

	// Create nodes
	for i := 1; i <= 5; i++ {
		dag.AddNode(i)
	}

	// Create edges: 1 -> 2 -> 3 and 1 -> 4 -> 5
	dag.AddEdge(1, 2)
	dag.AddEdge(2, 3)
	dag.AddEdge(1, 4)
	dag.AddEdge(4, 5)

	node1 := dag.nodes[1]
	node2 := dag.nodes[2]
	node3 := dag.nodes[3]
	node4 := dag.nodes[4]
	node5 := dag.nodes[5]

	// Test ancestor relationships
	assert.True(t, node3.hasAncestor(node1), "Node 1 should be an ancestor of Node 3")
	assert.True(t, node5.hasAncestor(node1), "Node 1 should be an ancestor of Node 5")
	assert.True(t, node3.hasAncestor(node2), "Node 2 should be an ancestor of Node 3")
	assert.True(t, node5.hasAncestor(node4), "Node 4 should be an ancestor of Node 5")
	assert.False(t, node2.hasAncestor(node3), "Node 3 should not be an ancestor of Node 2")
	assert.False(t, node4.hasAncestor(node2), "Node 2 should not be an ancestor of Node 4")
}
