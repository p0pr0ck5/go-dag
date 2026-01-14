package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type widget struct {
	id   int
	name string
}

func newWidget(id int, name string) *widget {
	return &widget{
		id:   id,
		name: name,
	}
}

func (w *widget) String() string {
	return w.name
}

func TestAddNode(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")

	n1 := dag.AddNode(w1)
	n2 := dag.AddNode(w2)

	assert.Equal(t, w1, n1.Data())
	assert.Equal(t, w2, n2.Data())

	// Adding the same node again should return the existing node
	n1Again := dag.AddNode(w1)
	assert.Same(t, n1, n1Again, "Expected the same node instance when adding duplicate node")
}

func TestAddEdge(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")

	dag.AddNode(w1)
	dag.AddNode(w2)

	err := dag.AddEdge(w1, w2)
	assert.NoError(t, err)

	n1 := dag.Node(w1)
	n2 := dag.Node(w2)

	assert.Len(t, n1.Children(), 1, "Widget1 should have one child")
	assert.Equal(t, n2, n1.Children()[0], "Widget1 should have Widget2 as child")
	assert.Len(t, n2.Parents(), 1, "Widget2 should have one parent")
	assert.Equal(t, n1, n2.Parents()[0], "Widget2 should have Widget1 as parent")

	// Test cycle detection
	err = dag.AddEdge(w2, w1)
	assert.Error(t, err, "Expected error when adding edge that creates a cycle")
}

func TestRemoveNode(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddEdge(w1, w2)

	dag.RemoveNode(w1)

	n2 := dag.Node(w2)
	assert.Len(t, n2.Parents(), 0, "Widget2 should have no parents after Widget1 is removed")

	n1 := dag.Node(w1)
	assert.Nil(t, n1, "Widget1 should be nil after removal")

	// Removing non-existent node should not cause error
	dag.RemoveNode(w1)

	// re-add w1 and edge to test removal of w2
	dag.AddNode(w1)
	dag.AddEdge(w1, w2)

	// remove w2 and check w1's children
	dag.RemoveNode(w2)
	n1 = dag.Node(w1)
	assert.Len(t, n1.Children(), 0, "Widget1 should have no children after Widget2 is removed")
}

func TestRemoveEdge(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddEdge(w1, w2)

	dag.RemoveEdge(w1, w2)

	n1 := dag.Node(w1)
	n2 := dag.Node(w2)

	assert.Len(t, n1.Children(), 0, "Widget1 should have no children after edge removal")
	assert.Len(t, n2.Parents(), 0, "Widget2 should have no parents after edge removal")

	// Removing non-existent edge should not cause error
	dag.RemoveEdge(w1, w2)

	// Removing edge where one node does not exist should not cause error
	dag.RemoveEdge(w1, newWidget(3, "Widget3"))
}

func TestClear(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddEdge(w1, w2)

	dag.Clear()

	assert.Len(t, dag.Nodes(), 0, "DAG should have no nodes after Clear()")
}

func TestNodes(t *testing.T) {
	dag := NewDAG[*widget]()

	widgets := []*widget{
		newWidget(1, "Widget1"),
		newWidget(2, "Widget2"),
		newWidget(3, "Widget3"),
	}

	for _, w := range widgets {
		dag.AddNode(w)
	}

	nodes := dag.Nodes()
	assert.Len(t, nodes, len(widgets), "Expected correct number of nodes")

	widgetSet := make(map[*widget]struct{})
	for _, w := range widgets {
		widgetSet[w] = struct{}{}
	}

	for _, node := range nodes {
		_, exists := widgetSet[node.Data()]
		assert.True(t, exists, "Unexpected node data: %v", node.Data())
	}
}

func TestEdges(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w2, w3)

	edges := dag.Edges()
	assert.Len(t, edges, 2, "Expected 2 edges in the DAG")

	expectedEdges := map[string]struct{}{
		"Widget1->Widget2": {},
		"Widget2->Widget3": {},
	}

	for _, edge := range edges {
		key := edge[0].Data().name + "->" + edge[1].Data().name
		_, exists := expectedEdges[key]
		assert.True(t, exists, "Unexpected edge: %s", key)
	}
}

func TestTraverse(t *testing.T) {
	dag := NewDAG[*widget]()

	// Create nodes
	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	// Create edges: Widget1 -> Widget2 -> Widget3
	dag.AddEdge(w1, w2)
	dag.AddEdge(w2, w3)

	sortedNodes, err := dag.Traverse()
	assert.NoError(t, err)

	expectedOrder := []*widget{w1, w2, w3}
	for i, node := range sortedNodes {
		assert.Equal(t, expectedOrder[i], node.Data(), "Expected node %v at position %d", expectedOrder[i], i)
	}
}

func TestWalk(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")
	w4 := newWidget(4, "Widget4")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddNode(w4)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w2, w3)
	dag.AddEdge(w3, w4)
	dag.AddEdge(w2, w4)

	visited := make(map[*widget]struct{})
	dag.Walk(w1, func(v *widget) {
		visited[v] = struct{}{}
	})

	assert.Len(t, visited, 4, "Expected to visit 4 nodes")
	assert.Contains(t, visited, w1, "Expected to visit Widget1")
	assert.Contains(t, visited, w2, "Expected to visit Widget2")
	assert.Contains(t, visited, w3, "Expected to visit Widget3")
	assert.Contains(t, visited, w4, "Expected to visit Widget4")
}

func TestReverseWalk(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")
	w4 := newWidget(4, "Widget4")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddNode(w4)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w2, w3)
	dag.AddEdge(w3, w4)
	dag.AddEdge(w1, w4)

	visited := make(map[*widget]struct{})
	dag.ReverseWalk(w4, func(v *widget) {
		visited[v] = struct{}{}
	})

	assert.Len(t, visited, 4, "Expected to visit 4 nodes")
	assert.Contains(t, visited, w1, "Expected to visit Widget1")
	assert.Contains(t, visited, w2, "Expected to visit Widget2")
	assert.Contains(t, visited, w3, "Expected to visit Widget3")
	assert.Contains(t, visited, w4, "Expected to visit Widget4")
}

func TestBreadthFirstWalk(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")
	w4 := newWidget(4, "Widget4")
	w5 := newWidget(5, "Widget5")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddNode(w4)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w1, w3)
	dag.AddEdge(w2, w4)
	dag.AddEdge(w3, w4)

	var visitOrder []string
	dag.BreadthFirstWalk(w1, func(v *widget) {
		visitOrder = append(visitOrder, v.name)
	})

	expectedOrder := []string{"Widget1", "Widget2", "Widget3", "Widget4"}
	assert.Equal(t, expectedOrder, visitOrder, "Expected breadth-first visit order")

	// double-check validation - do a depth-first walk to ensure
	// the proper order, which would be different than breadth-first
	var dfsVisitOrder []string
	dag.Walk(w1, func(v *widget) {
		dfsVisitOrder = append(dfsVisitOrder, v.name)
	})

	dfsExpectedOrder := []string{"Widget1", "Widget2", "Widget4", "Widget3"}
	assert.Equal(t, dfsExpectedOrder, dfsVisitOrder, "Expected depth-first visit order")
	assert.NotEqual(t, expectedOrder, dfsVisitOrder, "Expected depth-first visit order to differ from breadth-first")

	// nil check
	var nilVisitOrder []string
	dag.BreadthFirstWalk(w5, func(v *widget) {
		nilVisitOrder = append(nilVisitOrder, v.name)
	})
	assert.Empty(t, nilVisitOrder, "Expected no visits when starting from non-existent node")
}

func TestReverseBreadthFirstWalk(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")
	w4 := newWidget(4, "Widget4")
	w5 := newWidget(5, "Widget5")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddNode(w4)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w1, w3)
	dag.AddEdge(w2, w4)
	dag.AddEdge(w3, w4)

	var visitOrder []string
	dag.ReverseBreadthFirstWalk(w4, func(v *widget) {
		visitOrder = append(visitOrder, v.name)
	})

	expectedOrder := []string{"Widget4", "Widget2", "Widget3", "Widget1"}
	assert.Equal(t, expectedOrder, visitOrder, "Expected reverse breadth-first visit order")

	// double-check validation - do a depth-first walk to ensure
	// the proper order, which would be different than breadth-first
	var dfsVisitOrder []string
	dag.ReverseWalk(w4, func(v *widget) {
		dfsVisitOrder = append(dfsVisitOrder, v.name)
	})

	dfsExpectedOrder := []string{"Widget4", "Widget2", "Widget1", "Widget3"}
	assert.Equal(t, dfsExpectedOrder, dfsVisitOrder, "Expected reverse depth-first visit order")
	assert.NotEqual(t, expectedOrder, dfsVisitOrder, "Expected reverse depth-first visit order to differ from reverse breadth-first")

	// nil check
	var nilVisitOrder []string
	dag.ReverseBreadthFirstWalk(w5, func(v *widget) {
		nilVisitOrder = append(nilVisitOrder, v.name)
	})
	assert.Empty(t, nilVisitOrder, "Expected no visits when starting from non-existent node")
}

func TestLevelOrder(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")
	w4 := newWidget(4, "Widget4")
	w5 := newWidget(5, "Widget5")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)
	dag.AddNode(w4)
	dag.AddNode(w5)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w1, w3)
	dag.AddEdge(w2, w4)
	dag.AddEdge(w3, w4)
	dag.AddEdge(w4, w5)

	levels := dag.LevelOrder()
	assert.Len(t, levels, 4, "Expected 4 levels in the DAG")

	expectedLevels := [][]*widget{
		{w1},
		{w2, w3},
		{w4},
		{w5},
	}

	for i, level := range levels {
		assert.Len(t, level, len(expectedLevels[i]), "Expected correct number of nodes at level %d", i)
		for j, node := range level {
			assert.Equal(t, expectedLevels[i][j], node.Data(), "Expected correct node at level %d position %d", i, j)
		}
	}
}

func TestNode(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	dag.AddNode(w1)

	n := dag.Node(w1)
	assert.NotNil(t, n, "Expected to find node for Widget1")
	assert.Equal(t, w1, n.Data())

	w2 := newWidget(2, "Widget2")
	n2 := dag.Node(w2)
	assert.Nil(t, n2, "Did not expect to find node for Widget2")
}

func TestSearch(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	result := dag.Search(func(n *Node[*widget]) bool {
		return n.Data().name == "Widget2"
	})

	assert.NotNil(t, result, "Expected to find node for Widget2")
	assert.Equal(t, w2, result.Data())

	// Search for a non-existent widget
	result = dag.Search(func(n *Node[*widget]) bool {
		return n.Data().name == "Widget4"
	})

	assert.Nil(t, result, "Did not expect to find node for Widget4")
}

func TestRoots(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w1, w3)

	roots := dag.Roots()
	assert.Len(t, roots, 1, "Expected 1 root node")
	assert.Equal(t, w1, roots[0].Data(), "Expected root node to be Widget1")
}

func TestLeaves(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w1, w3)

	leaves := dag.Leaves()
	assert.Len(t, leaves, 2, "Expected 2 leaf nodes")
	leafSet := make(map[*widget]struct{})
	for _, leaf := range leaves {
		leafSet[leaf.Data()] = struct{}{}
	}
	assert.Contains(t, leafSet, w2, "Expected Widget2 to be a leaf node")
	assert.Contains(t, leafSet, w3, "Expected Widget3 to be a leaf node")
}

func TestVisualize(t *testing.T) {
	dag := NewDAG[*widget]()

	w1 := newWidget(1, "Widget1")
	w2 := newWidget(2, "Widget2")
	w3 := newWidget(3, "Widget3")

	dag.AddNode(w1)
	dag.AddNode(w2)
	dag.AddNode(w3)

	dag.AddEdge(w1, w2)
	dag.AddEdge(w2, w3)

	dot := dag.Visualize()
	expectedSubstrings := []string{
		`"Widget1" -> "Widget2";`,
		`"Widget2" -> "Widget3";`,
	}

	for _, substr := range expectedSubstrings {
		assert.Contains(t, dot, substr, "Expected DOT output to contain: %s", substr)
	}
}

func TestHasEdge(t *testing.T) {
	dag := NewDAG[int]()

	// Create nodes
	for i := 1; i <= 3; i++ {
		dag.AddNode(i)
	}

	// Create edges: 1 -> 2 and 2 -> 3
	dag.AddEdge(1, 2)
	dag.AddEdge(2, 3)

	// Test existing edges
	assert.True(t, dag.HasEdge(1, 2), "Expected edge from 1 to 2")
	assert.True(t, dag.HasEdge(2, 3), "Expected edge from 2 to 3")

	// Test non-existing edges
	assert.False(t, dag.HasEdge(1, 3), "Did not expect edge from 1 to 3")
	assert.False(t, dag.HasEdge(3, 1), "Did not expect edge from 3 to 1")

	// Test non-existing nodes
	assert.False(t, dag.HasEdge(4, 1), "Did not expect edge from non-existent node 4 to 1")
	assert.False(t, dag.HasEdge(1, 4), "Did not expect edge from 1 to non-existent node 4")
}

func TestHasPath(t *testing.T) {
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

	// Test paths
	assert.True(t, dag.HasPath(1, 3), "Expected path from 1 to 3")
	assert.True(t, dag.HasPath(1, 5), "Expected path from 1 to 5")
	assert.False(t, dag.HasPath(3, 1), "Did not expect path from 3 to 1")
	assert.False(t, dag.HasPath(2, 5), "Did not expect path from 2 to 5")
}

func TestHasPathNilNodes(t *testing.T) {
	dag := NewDAG[int]()

	// Create one node
	dag.AddNode(1)

	// Test with non-existent from node
	assert.False(t, dag.HasPath(999, 1), "Expected false when from node doesn't exist")

	// Test with non-existent to node
	assert.False(t, dag.HasPath(1, 999), "Expected false when to node doesn't exist")

	// Test with both nodes non-existent
	assert.False(t, dag.HasPath(888, 999), "Expected false when both nodes don't exist")
}

func TestShortestPath(t *testing.T) {
	dag := NewDAG[int]()

	// Create nodes
	for i := 1; i <= 7; i++ {
		dag.AddNode(i)
	}

	// Create edges: 1 -> 2 -> 3, 1 -> 4 -> 3, and 2 -> 5
	dag.AddEdge(1, 2)
	dag.AddEdge(2, 3)
	dag.AddEdge(1, 4)
	dag.AddEdge(4, 3)
	dag.AddEdge(2, 5)
	dag.AddEdge(6, 7) // disconnected component

	// Test shortest path from 1 to 5
	path := dag.ShortestPath(1, 5)
	expectedPath := []int{1, 2, 5}

	assert.Len(t, path, len(expectedPath), "Expected path length %d", len(expectedPath))
	for i, node := range path {
		assert.Equal(t, expectedPath[i], node.Data(), "Expected node %d at position %d", expectedPath[i], i)
	}

	path = dag.ShortestPath(1, 7)
	assert.Empty(t, path, "Expected no path from 1 to 7")
}

func TestShortestPathNilNodes(t *testing.T) {
	dag := NewDAG[int]()

	// Create one node
	dag.AddNode(1)

	// Test with non-existent from node
	path := dag.ShortestPath(999, 1)
	assert.Nil(t, path, "Expected nil when from node doesn't exist")

	// Test with non-existent to node
	path = dag.ShortestPath(1, 999)
	assert.Nil(t, path, "Expected nil when to node doesn't exist")

	// Test with both nodes non-existent
	path = dag.ShortestPath(888, 999)
	assert.Nil(t, path, "Expected nil when both nodes don't exist")
}

func TestAncestors(t *testing.T) {
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

	// Test ancestor relationships
	a := dag.Ancestors(5)
	assert.Len(t, a, 2, "Expected 2 ancestors for Node 5")
	assert.Contains(t, a, dag.nodes[1], "Expected Node 1 to be an ancestor of Node 5")
	assert.Contains(t, a, dag.nodes[4], "Expected Node 4 to be an ancestor of Node 5")

	a = dag.Ancestors(3)
	assert.Len(t, a, 2, "Expected 2 ancestors for Node 3")
	assert.Contains(t, a, dag.nodes[1], "Expected Node 1 to be an ancestor of Node 3")
	assert.Contains(t, a, dag.nodes[2], "Expected Node 2 to be an ancestor of Node 3")
}

func TestAncestorsNilNode(t *testing.T) {
	dag := NewDAG[int]()

	// Test ancestors of non-existent node
	ancestors := dag.Ancestors(999)
	assert.Nil(t, ancestors, "Expected nil when node doesn't exist")
}

func TestDescendants(t *testing.T) {
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

	// Test descendant relationships
	d := dag.Descendants(1)
	assert.Len(t, d, 4, "Expected 4 descendants for Node 1")
	assert.Contains(t, d, dag.nodes[2], "Expected Node 2 to be a descendant of Node 1")
	assert.Contains(t, d, dag.nodes[3], "Expected Node 3 to be a descendant of Node 1")
	assert.Contains(t, d, dag.nodes[4], "Expected Node 4 to be a descendant of Node 1")
	assert.Contains(t, d, dag.nodes[5], "Expected Node 5 to be a descendant of Node 1")

	d = dag.Descendants(2)
	assert.Len(t, d, 1, "Expected 1 descendant for Node 2")
	assert.Contains(t, d, dag.nodes[3], "Expected Node 3 to be a descendant of Node 2")
}

func TestDescendantsNilNode(t *testing.T) {
	dag := NewDAG[int]()

	// Test descendants of non-existent node
	descendants := dag.Descendants(999)
	assert.Nil(t, descendants, "Expected nil when node doesn't exist")
}

func TestWalkNilNode(t *testing.T) {
	dag := NewDAG[int]()

	// Test walking from non-existent node
	visitCount := 0
	dag.Walk(999, func(v int) {
		visitCount++
	})

	assert.Equal(t, 0, visitCount, "Expected no visits when starting from non-existent node")
}

func TestReverseWalkNilNode(t *testing.T) {
	dag := NewDAG[int]()

	// Test reverse walking from non-existent node
	visitCount := 0
	dag.ReverseWalk(999, func(v int) {
		visitCount++
	})

	assert.Equal(t, 0, visitCount, "Expected no visits when starting from non-existent node")
}
