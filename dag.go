package dag

import "fmt"

// ErrCycleDetected is returned when an operation would create a cycle in the DAG.
var ErrCycleDetected = fmt.Errorf("adding this edge would create a cycle")

// DAG represents a directed acyclic graph.
type DAG[T comparable] struct {
	nodes map[T]*Node[T]
}

// NewDAG creates and returns a new empty DAG.
func NewDAG[T comparable]() *DAG[T] {
	return &DAG[T]{
		nodes: make(map[T]*Node[T]),
	}
}

// AddNode adds a node with the given data to the DAG.
// If a node with the same data already exists, it returns the existing node.
func (d *DAG[T]) AddNode(data T) *Node[T] {
	if node, exists := d.nodes[data]; exists {
		return node
	}
	node := NewNode(data)
	d.nodes[data] = node
	return node
}

// AddEdge adds a directed edge from the node with data 'from' to the node with data 'to'.
// It returns an error if adding the edge would create a cycle.
func (d *DAG[T]) AddEdge(from, to T) error {
	fromNode := d.AddNode(from)
	toNode := d.AddNode(to)

	// Check for cycles
	if fromNode.hasAncestor(toNode) {
		return ErrCycleDetected
	}

	fromNode.addChild(toNode)
	toNode.addParent(fromNode)
	return nil
}

// RemoveNode removes the node with the given data from the DAG.
func (d *DAG[T]) RemoveNode(data T) {
	node, exists := d.nodes[data]
	if !exists {
		return
	}

	// Remove this node from its parents' children
	for parent := range node.parents {
		delete(parent.children, node)
	}

	// Remove this node from its children's parents
	for child := range node.children {
		delete(child.parents, node)
	}

	delete(d.nodes, data)
}

// RemoveEdge removes the directed edge from the node with data 'from' to the node with data 'to'.
func (d *DAG[T]) RemoveEdge(from, to T) {
	fromNode, fromExists := d.nodes[from]
	toNode, toExists := d.nodes[to]
	if !fromExists || !toExists {
		return
	}

	delete(fromNode.children, toNode)
	delete(toNode.parents, fromNode)
}

// Clear removes all nodes and edges from the DAG.
func (d *DAG[T]) Clear() {
	d.nodes = make(map[T]*Node[T])
}

// Nodes returns all nodes in the DAG.
func (d *DAG[T]) Nodes() []*Node[T] {
	nodes := make([]*Node[T], 0, len(d.nodes))
	for _, node := range d.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// Edges returns all edges in the DAG as pairs of nodes.
func (d *DAG[T]) Edges() [][2]*Node[T] {
	var edges [][2]*Node[T]
	for _, node := range d.nodes {
		// Use deterministic iteration order
		for _, child := range node.Children() {
			edges = append(edges, [2]*Node[T]{node, child})
		}
	}
	return edges
}

// Traverse performs a topological sort of the DAG and returns the nodes in sorted order.
func (d *DAG[T]) Traverse() ([]*Node[T], error) {
	inDegree := make(map[*Node[T]]int)
	for _, node := range d.nodes {
		inDegree[node] = len(node.parents)
	}

	var queue []*Node[T]
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	var sorted []*Node[T]
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		// Use deterministic iteration order
		for _, child := range current.Children() {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}

	return sorted, nil
}

// Walk performs a depth-first traversal starting from the node with the given data.
func (d *DAG[T]) Walk(start T, visit func(v T)) {
	startNode := d.nodes[start]
	if startNode == nil {
		return
	}

	visited := make(map[*Node[T]]struct{})
	var walk func(node *Node[T])
	walk = func(node *Node[T]) {
		if _, seen := visited[node]; seen {
			return
		}
		visited[node] = struct{}{}
		visit(node.Data())
		// Use deterministic iteration order
		for _, child := range node.Children() {
			walk(child)
		}
	}
	walk(startNode)
}

// ReverseWalk performs a depth-first traversal in reverse
// starting from the node with the given data.
func (d *DAG[T]) ReverseWalk(start T, visit func(v T)) {
	startNode := d.nodes[start]
	if startNode == nil {
		return
	}

	visited := make(map[*Node[T]]struct{})
	var walk func(node *Node[T])
	walk = func(node *Node[T]) {
		if _, seen := visited[node]; seen {
			return
		}
		visited[node] = struct{}{}
		visit(node.Data())
		// Use deterministic iteration order
		for _, parent := range node.Parents() {
			walk(parent)
		}
	}
	walk(startNode)
}

// BreadthFirstWalk performs a breadth-first traversal starting from the node with the given data.
func (d *DAG[T]) BreadthFirstWalk(start T, visit func(v T)) {
	startNode := d.nodes[start]
	if startNode == nil {
		return
	}

	visited := make(map[*Node[T]]struct{})
	queue := []*Node[T]{startNode}
	visited[startNode] = struct{}{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visit(current.Data())

		// Use deterministic iteration order
		for _, child := range current.Children() {
			if _, seen := visited[child]; !seen {
				visited[child] = struct{}{}
				queue = append(queue, child)
			}
		}
	}
}

// ReverseBreadthFirstWalk performs a breadth-first traversal in reverse
// starting from the node with the given data.
func (d *DAG[T]) ReverseBreadthFirstWalk(start T, visit func(v T)) {
	startNode := d.nodes[start]
	if startNode == nil {
		return
	}

	visited := make(map[*Node[T]]struct{})
	queue := []*Node[T]{startNode}
	visited[startNode] = struct{}{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visit(current.Data())

		// Use deterministic iteration order
		for _, parent := range current.Parents() {
			if _, seen := visited[parent]; !seen {
				visited[parent] = struct{}{}
				queue = append(queue, parent)
			}
		}
	}
}

// LevelOrder returns the nodes of the DAG in level order.
func (d *DAG[T]) LevelOrder() [][]*Node[T] {
	inDegree := make(map[*Node[T]]int)
	for _, node := range d.nodes {
		inDegree[node] = len(node.parents)
	}

	var levels [][]*Node[T]
	currentLevel := make([]*Node[T], 0)

	for node, degree := range inDegree {
		if degree == 0 {
			currentLevel = append(currentLevel, node)
		}
	}

	for len(currentLevel) > 0 {
		levels = append(levels, currentLevel)
		nextLevel := make([]*Node[T], 0)

		// Use deterministic iteration order
		for _, node := range currentLevel {
			for _, child := range node.Children() {
				inDegree[child]--
				if inDegree[child] == 0 {
					nextLevel = append(nextLevel, child)
				}
			}
		}

		currentLevel = nextLevel
	}

	return levels
}

// Node finds and returns the node with the given data, or nil if it doesn't exist.
func (d *DAG[T]) Node(data T) *Node[T] {
	return d.nodes[data]
}

// Search finds a node based on a predicate function.
func (d *DAG[T]) Search(predicate func(node *Node[T]) bool) *Node[T] {
	for _, node := range d.nodes {
		if predicate(node) {
			return node
		}
	}
	return nil
}

// Roots returns all root nodes (nodes with no parents) in the DAG.
func (d *DAG[T]) Roots() []*Node[T] {
	var roots []*Node[T]
	for _, node := range d.nodes {
		if len(node.parents) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}

// Leaves returns all leaf nodes (nodes with no children) in the DAG.
func (d *DAG[T]) Leaves() []*Node[T] {
	var leaves []*Node[T]
	for _, node := range d.nodes {
		if len(node.children) == 0 {
			leaves = append(leaves, node)
		}
	}
	return leaves
}

// Visualize generates a DOT format representation of the DAG for visualization.
func (d *DAG[T]) Visualize() string {
	result := "digraph G {\n"
	for _, node := range d.nodes {
		// Use deterministic iteration order
		for _, child := range node.Children() {
			result += fmt.Sprintf("    \"%v\" -> \"%v\";\n", node.Data(), child.Data())
		}
	}
	result += "}\n"
	return result
}

// HasEdge checks if there is a directed edge from the node with data 'from'
// to the node with data 'to'.
func (d *DAG[T]) HasEdge(from, to T) bool {
	fromNode := d.nodes[from]
	toNode := d.nodes[to]
	if fromNode == nil || toNode == nil {
		return false
	}
	for _, child := range fromNode.Children() {
		if child == toNode {
			return true
		}
	}
	return false
}

// HasPath checks if there is a path from the node with data 'from'
// to the node with data 'to'.
func (d *DAG[T]) HasPath(from, to T) bool {
	fromNode := d.nodes[from]
	toNode := d.nodes[to]
	if fromNode == nil || toNode == nil {
		return false
	}

	visited := make(map[*Node[T]]struct{})
	var visit func(node *Node[T]) bool
	visit = func(node *Node[T]) bool {
		if node == toNode {
			return true
		}
		visited[node] = struct{}{}
		// Use deterministic iteration order
		for _, child := range node.Children() {
			if _, seen := visited[child]; !seen {
				if visit(child) {
					return true
				}
			}
		}
		return false
	}

	return visit(fromNode)
}

// ShortestPath finds the shortest path from the node with data 'from'
// to the node with data 'to'.
func (d *DAG[T]) ShortestPath(from, to T) []*Node[T] {
	fromNode := d.nodes[from]
	toNode := d.nodes[to]
	if fromNode == nil || toNode == nil {
		return nil
	}

	type pathNode struct {
		node *Node[T]
		path []*Node[T]
	}

	queue := []pathNode{{node: fromNode, path: []*Node[T]{fromNode}}}
	visited := make(map[*Node[T]]struct{})
	visited[fromNode] = struct{}{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.node == toNode {
			return current.path
		}

		// Use deterministic iteration order
		for _, child := range current.node.Children() {
			if _, seen := visited[child]; !seen {
				visited[child] = struct{}{}
				newPath := append([]*Node[T]{}, current.path...)
				newPath = append(newPath, child)
				queue = append(queue, pathNode{node: child, path: newPath})
			}
		}
	}

	return nil
}

// Ancestors returns all ancestor nodes of the node with the given data.
func (d *DAG[T]) Ancestors(data T) []*Node[T] {
	startNode := d.nodes[data]
	if startNode == nil {
		return nil
	}
	var ancestors []*Node[T]
	visited := make(map[*Node[T]]struct{})

	var visit func(node *Node[T])
	visit = func(node *Node[T]) {
		for _, parent := range node.Parents() {
			if _, seen := visited[parent]; !seen {
				visited[parent] = struct{}{}
				ancestors = append(ancestors, parent)
				visit(parent)
			}
		}
	}

	visit(startNode)
	return ancestors
}

// Descendants returns all descendant nodes of the node with the given data.
func (d *DAG[T]) Descendants(data T) []*Node[T] {
	startNode := d.nodes[data]
	if startNode == nil {
		return nil
	}
	var descendants []*Node[T]
	visited := make(map[*Node[T]]struct{})

	var visit func(node *Node[T])
	visit = func(node *Node[T]) {
		for _, child := range node.Children() {
			if _, seen := visited[child]; !seen {
				visited[child] = struct{}{}
				descendants = append(descendants, child)
				visit(child)
			}
		}
	}

	visit(startNode)
	return descendants
}
