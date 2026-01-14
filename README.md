# DAG

[![CI](https://github.com/p0pr0ck5/go-dag/actions/workflows/ci.yml/badge.svg)](https://github.com/p0pr0ck5/go-dag/actions/workflows/ci.yml)

A tiny, generic Directed Acyclic Graph (DAG) library for Go. It provides cycle-safe edge creation, deterministic traversal, and handy helpers like topological sort, roots/leaves, ancestors/descendants, shortest paths, and DOT output for visualization.

## Install

```bash
go get github.com/p0pr0ck5/go-dag@latest
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/p0pr0ck5/go-dag"
)

func main() {
    g := dag.NewDAG[int]()

    // Nodes are created implicitly when adding edges
    _ = g.AddEdge(1, 2)
    _ = g.AddEdge(2, 3)

    // Topological order
    order, _ := g.Traverse()
    for _, n := range order {
        fmt.Println(n.Data())
    }

    // Queries
    fmt.Println("HasPath 1->3:", g.HasPath(1, 3))

    // DOT visualization
    fmt.Println(g.Visualize())
}
```

## Features

- Generic `DAG[T]` and `Node[T]` with `T` comparable
- Cycle detection on `AddEdge`
- Deterministic `Walk`, `ReverseWalk`, BFS variants, and `LevelOrder`
- Topological sort via `Traverse`
- `Roots()`, `Leaves()`, `Ancestors()`, `Descendants()`
- `HasEdge()`, `HasPath()`, `ShortestPath()`
- `Visualize()` → DOT format (works with Graphviz)

## Test

```bash
go test ./...
```
