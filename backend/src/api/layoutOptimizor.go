package api

import (
	"congest/graph"
)

// SortRoads ..
func sortRoads(roadsID []int, nodes []*graph.Node) {
	indexOfNodeID := 0

	var nodeDict map[int]int
	nodeDict = make(map[int]int)
	for i, node := range nodes {
		nodeDict[node.ID] = i
	}

	for i, roadID := range roadsID {
		if _, ok := nodeDict[roadID]; ok {
			if i == indexOfNodeID {
				continue
			}
			tmp := roadsID[i]
			roadsID[i] = roadsID[indexOfNodeID]
			roadsID[indexOfNodeID] = tmp

			indexOfNodeID++
		}
	}
}

// OptimizeLayout ..
func optimizeBySwap(roadsID []int, g *graph.Graph) {
	k := 10
	curLoss := loss(roadsID, g)
	for k > 0 {
		breakFlag := false
		for i1 := range roadsID {
			for i2 := range roadsID {
				// swap attempt
				tmp := roadsID[i1]
				roadsID[i1] = roadsID[i2]
				roadsID[i2] = tmp

				tmpLoss := loss(roadsID, g)

				if tmpLoss < curLoss {
					curLoss = tmpLoss
					breakFlag = true
				} else {
					// undo swap
					tmp := roadsID[i1]
					roadsID[i1] = roadsID[i2]
					roadsID[i2] = tmp
				}

				if breakFlag {
					break
				}
			}

			if breakFlag {
				break
			}
		}

		k--
	}
}

func loss(roadsID []int, g *graph.Graph) float64 {
	sumLength := 0
	sumCross := 0
	links := g.GetLinkList()

	var nodeDict map[int]int
	nodeDict = make(map[int]int)
	for i, roadID := range roadsID {
		nodeDict[roadID] = i
	}

	for li, link := range links {
		o, d := link[0], link[1]
		var id1, id2 int // idWithSmallerIndex, idWithLargerIndex

		if nodeDict[o] > nodeDict[d] {
			id1, id2 = nodeDict[d], nodeDict[o]
		} else {
			id2, id1 = nodeDict[d], nodeDict[o]
		}

		len := id2 - id1 // length of the link

		nCross := 0
		for lj, linkj := range links {
			if li == lj {
				continue
			}

			oj, dj := linkj[0], linkj[1]
			var id1j, id2j int // idWithSmallerIndex, idWithLargerIndex

			if nodeDict[oj] > nodeDict[dj] {
				id1j, id2j = nodeDict[dj], nodeDict[oj]
			} else {
				id2j, id1j = nodeDict[dj], nodeDict[oj]
			}

			if ((id1 < id1j) && (id1j < id2) && (id2 < id2j)) || ((id1j < id1) && (id1 < id2j) && (id2j < id2)) {
				nCross++
			}
		}

		sumCross += nCross
		sumLength += len
	}

	return float64(sumCross+sumLength) * 1.0
}

// LayoutOptimize ..
func LayoutOptimize(roadsID []int, g *graph.Graph) map[int]int {
	var originMap map[int]int
	originMap = make(map[int]int)
	for i, roadID := range roadsID {
		originMap[roadID] = i
	}

	sortRoads(roadsID, g.Nodes)
	optimizeBySwap(roadsID, g)

	var orderMap map[int]int
	orderMap = make(map[int]int)
	for i, roadID := range roadsID {
		orderMap[i] = originMap[roadID]
	}
	return orderMap
}
