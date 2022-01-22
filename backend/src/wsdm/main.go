package main

import (
	"gClass"
)

func main() {
	var graph gClass.IGraph
	na := gClass.Node{"a"}
	nb := gClass.Node{"b"}
	nc := gClass.Node{"c"}
	nd := gClass.Node{"d"}

	graph.AddNode(&na)
	graph.AddNode(&nb)
	graph.AddNode(&nc)
	graph.AddNode(&nd)

	graph.AddEdge(14, &na, &nd)
	graph.AddEdge(15, &nc, &na)
	graph.AddEdge(17, &na, &nc)
	graph.AddEdge(25, &na, &nb)
	graph.AddEdge(28, &na, &nc)
	graph.AddEdge(30, &na, &nc)
	graph.AddEdge(31, &nc, &nd)
	graph.AddEdge(32, &nc, &na)
	graph.AddEdge(35, &na, &nc)

	var tmotif gClass.TMotif
	e1 := gClass.SimpleEdge{&na, &nb}
	e2 := gClass.SimpleEdge{&na, &nc}
	e3 := gClass.SimpleEdge{&nc, &na}

	tmotif.AddNode(&na)
	tmotif.AddNode(&nb)
	tmotif.AddNode(&nc)

	tmotif.AddOrderedEdge(&e1)
	tmotif.AddOrderedEdge(&e2)
	tmotif.AddOrderedEdge(&e3)

	graph.MotifCounting(&tmotif, 10)
}
