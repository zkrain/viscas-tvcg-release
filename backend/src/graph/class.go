package graph

import (
	"fmt"
	"sync"
)

// Node ..
type Node struct {
	ID int // roadID
}

// NodeString get the string format of node n
func (n *Node) NodeString() string {
	return fmt.Sprintf("%v", n.ID)
}

// Graph ..
type Graph struct {
	Nodes     []*Node
	Edges     map[int][]int
	edgeCount int
	lock      sync.RWMutex
}

// CountEdge ..
func (g *Graph) CountEdge() int {
	return g.edgeCount
}

// IsNodeExist .. if exist, return true, else return false
func (g *Graph) IsNodeExist(n *Node) bool {
	for _, node := range g.Nodes {
		if node.ID == n.ID {
			return true
		}
	}

	return false
}

// IsEdgeExist ..
func (g *Graph) IsEdgeExist(ni, nj *Node) bool {
	if g.Edges == nil {
		return false
	}

	if len(g.Edges[ni.ID]) == 0 {
		return false
	}

	for _, njID := range g.Edges[ni.ID] {
		if njID == nj.ID {
			return true
		}
	}

	return false
}

// AddNode ..
func (g *Graph) AddNode(n *Node) {
	if g.IsNodeExist(n) {
		return
	}

	g.lock.Lock()
	g.Nodes = append(g.Nodes, n)
	g.lock.Unlock()
}

// AddEdge ..
func (g *Graph) AddEdge(ni, nj *Node) { // each edge is directed
	g.lock.Lock()
	if g.Edges == nil {
		g.Edges = make(map[int][]int)
	}
	g.Edges[ni.ID] = append(g.Edges[ni.ID], nj.ID)
	g.edgeCount++
	// g.Edges[ni.ID] = append(g.Edges[ni.ID], nj.ID)
	g.lock.Unlock()
}

// GetLinkList ..
func (g *Graph) GetLinkList() [][2]int {
	var linkList [][2]int
	for roadIID, roadJIDs := range g.Edges {
		for _, roadJID := range roadJIDs {
			linkList = append(linkList, [2]int{roadIID, roadJID})
		}
	}
	return linkList
}
