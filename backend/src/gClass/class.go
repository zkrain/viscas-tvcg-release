package gClass

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Node struct {
	ID string
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.ID)
}

type TemporalEdge struct {
	O         *Node
	D         *Node
	Timestamp int
}

type TemporalEdgeSequence []*TemporalEdge

func (S TemporalEdgeSequence) Len() int { return len(S) }
func (S TemporalEdgeSequence) Less(i, j int) bool {
	return S[i].Timestamp < S[j].Timestamp
}
func (S TemporalEdgeSequence) Swap(i, j int) { S[i], S[j] = S[j], S[i] }

type IGraph struct {
	Nodes []*Node
	Edges []*TemporalEdge
	lock  sync.RWMutex
}

func (g *IGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.Nodes = append(g.Nodes, n)
	g.lock.Unlock()
}

func (g *IGraph) AddEdge(t int, n1, n2 *Node) {
	g.lock.Lock()
	// if g.Edges == nil {
	// 	g.Edges = make(map[Node]map[int]*Node)
	// }
	// et := make(map[int]*Node)
	// et[t] = n2
	// g.Edges[*n1] = et
	// g.Edges[*n1] = append(g.Edges[*n1], n2)
	// g.Edges[*n2] = append(g.Edges[*n2], n1)
	g.Edges = append(g.Edges, &TemporalEdge{n1, n2, t})
	g.lock.Unlock()
}

type SimpleEdge struct {
	O *Node
	D *Node
}

type TMotif struct {
	Nodes        []*Node
	OrderedEdges []*SimpleEdge
}

func (tm *TMotif) AddNode(n *Node) {
	tm.Nodes = append(tm.Nodes, n)
}

func (tm *TMotif) AddOrderedEdge(e *SimpleEdge) {
	tm.OrderedEdges = append(tm.OrderedEdges, e)
}

type byLength []string

func (s byLength) Len() int           { return len(s) }
func (s byLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byLength) Less(i, j int) bool { return strings.Count(s[i], ";") < strings.Count(s[j], ";") }

type byLengthReverse []string

func (s byLengthReverse) Len() int      { return len(s) }
func (s byLengthReverse) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byLengthReverse) Less(i, j int) bool {
	return strings.Count(s[i], ";") > strings.Count(s[j], ";")
}

func getKeySlice(counts *map[string]int, isReverse bool, l int) []string {
	var keySlice []string
	for suffix := range *counts {
		keyLength := strings.Count(suffix, ";")
		if keyLength < l {
			keySlice = append(keySlice, suffix)
		}
	}

	if isReverse {
		sort.Sort(byLengthReverse(keySlice))
	} else {
		sort.Sort(byLength(keySlice))
	}

	return keySlice
}

func DecrementCounts(e *TemporalEdge, counts *map[string]int, l int) {
	ekey := e.O.ID + e.D.ID + ";"
	(*counts)[ekey]--

	keySlice := getKeySlice(counts, false, l-1)
	for _, suffix := range keySlice {
		(*counts)[ekey+suffix] -= (*counts)[suffix]
	}
}

func IncrementCounts(e *TemporalEdge, counts *map[string]int, l int) {
	ekey := e.O.ID + e.D.ID + ";"

	keySlice := getKeySlice(counts, true, l)
	for _, prefix := range keySlice {
		_, ok := (*counts)[prefix+ekey]
		if ok {
			(*counts)[prefix+ekey] += (*counts)[prefix]
		} else {
			(*counts)[prefix+ekey] = (*counts)[prefix]
		}
	}

	_, ok := (*counts)[ekey]
	if ok {
		(*counts)[ekey]++
	} else {
		(*counts)[ekey] = 1
	}
}

func (g *IGraph) MotifCounting(tm *TMotif, tDelta int) {
	l := len(tm.OrderedEdges)

	// retrieve edges
	var S TemporalEdgeSequence
	for _, simpleEdge := range tm.OrderedEdges {
		O := simpleEdge.O
		D := simpleEdge.D
		for _, edge := range g.Edges {
			if edge.O.ID == O.ID && edge.D.ID == D.ID {
				S = append(S, edge)
			}
		}
	}

	// sort edges
	sort.Sort(S)
	counts := make(map[string]int)

	start := 0
	for end, temporalEdge := range S {
		fmt.Println(temporalEdge.O.ID, temporalEdge.D.ID, temporalEdge.Timestamp)

		tstart := S[start].Timestamp
		tend := S[end].Timestamp
		for tstart+tDelta < tend {
			DecrementCounts(S[start], &counts, l)
			start++
			tstart = S[start].Timestamp
		}

		IncrementCounts(S[end], &counts, l)

		for key, count := range counts {
			fmt.Println(key, count)
		}
		fmt.Println("......divider......")
	}
}
