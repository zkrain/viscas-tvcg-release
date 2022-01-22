package api

// import (
// 	"congest/graph"
// 	// "congest/api"
// 	"fmt"
// 	"math"

// 	geo "github.com/kellydunn/golang-geo"
// )

// // InferEdge ..
// type InferEdge struct {
// 	StartTimeI int
// 	EndTimeI   int
// 	StartTimeJ int
// 	EndTimeJ   int
// 	P          float64
// 	Strength   int
// }

// type InferEdges = []InferEdge

// type RevertEdgeIndexInfo struct {
// 	CausalLocationID int
// 	EdgeIndex        int
// }

// // CalculateDistanceBetweenLocations return in km
// func CalculateDistanceBetweenLocations(locationI Location, locationJ Location) float64 {
// 	var pointI *geo.Point
// 	var pointJ *geo.Point
// 	pointI = geo.NewPoint(float64(locationI.Lat), float64(locationI.Lng))
// 	pointJ = geo.NewPoint(float64(locationJ.Lat), float64(locationJ.Lng))
// 	d := pointJ.GreatCircleDistance(pointI)
// 	return d
// }

// // GetDistanceMatrix return a distance dictionary like map[int]map[int]float64
// func GetDistanceMatrix(locationIDs []int, locations []Location) map[int]map[int]float64 {
// 	var distanceMatrix map[int]map[int]float64
// 	distanceMatrix = make(map[int]map[int]float64)

// 	maxValidDistance := 0.0
// 	minValidDistance := 10000.0

// 	// 问题在于，输入进来的地点之间的距离能否确保不会太大？
// 	for i, locationIDi := range locationIDs {
// 		var row map[int]float64
// 		row = make(map[int]float64)

// 		for j, locationIDj := range locationIDs {
// 			if i != j {
// 				distance := CalculateDistanceBetweenLocations(locations[i], locations[j])

// 				if distance > maxValidDistance {
// 					maxValidDistance = distance
// 				}
// 				if distance < minValidDistance {
// 					minValidDistance = distance
// 				}

// 				row[locationIDj] = distance
// 			}
// 		}

// 		distanceMatrix[locationIDi] = row
// 	}

// 	for locationIDi, row := range distanceMatrix {
// 		for locationIDj, d := range row {
// 			distanceMatrix[locationIDi][locationIDj] = (d - minValidDistance) / (maxValidDistance - minValidDistance)
// 		}
// 	}

// 	return distanceMatrix
// }

// // GetEdges ..
// func GetEdges(locationIDs []int, timeRangesList [][]int, distanceMatrix map[int]map[int]float64, filename string) map[int]map[int]InferEdges {
// 	var edges map[int]map[int]InferEdges
// 	edges = make(map[int]map[int]InferEdges)

// 	for i, locationIDi := range locationIDs {
// 		timeRangesOfLocationI := timeRangesList[i]
// 		for j, locationIDj := range locationIDs {
// 			if i == j {
// 				continue
// 			}

// 			timeRangesOfLocationJ := timeRangesList[j]
// 			distance := distanceMatrix[locationIDj][locationIDi]

// 			edgeList := GetCausalLinksBetweenLocations(timeRangesOfLocationI, timeRangesOfLocationJ, distance, filename)

// 			var subEdges map[int]InferEdges
// 			subEdges = make(map[int]InferEdges)
// 			if _, ok := edges[locationIDi]; ok {
// 				edges[locationIDi][locationIDj] = edgeList
// 			} else {
// 				subEdges[locationIDj] = edgeList
// 				edges[locationIDi] = subEdges
// 			}
// 		}
// 	}

// 	return edges
// }

// // GetCausalLinksBetweenLocations ...
// func GetCausalLinksBetweenLocations(timeRangesOfLocationI []int, timeRangesOfLocationJ []int, distance float64, filename string) []InferEdge {
// 	var length int
// 	if filename == "congestion" || filename == "congestion2" {
// 		length = 8785
// 	} else if filename == "flow" {
// 		length = 7221
// 	} else {
// 		length = 10000
// 	}
// 	// inputs
// 	// timeRangesOfLocationI: casual location
// 	// timeRangesOfLocationJ: location that is affected

// 	// dissipation

// 	// // store with map/dict
// 	// // map[congestTimeI]map[startTimeJ]{}
// 	// var edgeDict map[int]map[int]InferEdge
// 	// edgeDict = make(map[int]map[int]InferEdge)

// 	// store with slice/list
// 	edgeList := []InferEdge{}

// 	a := 2.0
// 	// lambda := 25.0
// 	beta := 18.0
// 	gamma := 1.0

// 	for startTimeI, endTimeI := range timeRangesOfLocationI {
// 		if endTimeI == 0 { // means no congestion
// 			continue
// 		}

// 		// var subEdgeDict map[int]InferEdge
// 		// subEdgeDict = make(map[int]InferEdge)

// 		// for tj := startTimeI + 1; tj < endTimeI+4 && tj < length; tj++ { // The road section B is congested before the congestion of the road section A dissipates
// 		for tj := startTimeI + 1; tj < startTimeI+7 && tj < length; tj++ {
// 			if timeRangesOfLocationJ[tj] > 0 {
// 				startTimeJ := tj
// 				endTimeJ := timeRangesOfLocationJ[tj]
// 				// Transmission Likelihood: (d, startTimeJ - startTimeI), (endTimeJ - endTimeI))
// 				// Degree of Impact: (endTimeJ -startTimeJ)

// 				A1 := float64(startTimeJ-startTimeI) / 6.0                             // normalized
// 				A2 := distance                                                         // have been normalized
// 				A3 := math.Abs(float64(endTimeJ - startTimeJ - endTimeI + startTimeI)) // abs(L1 - L2)
// 				A4 := endTimeJ - startTimeJ + endTimeI - startTimeI                    // L1 + L2
// 				weight := float64(A4) - A3*gamma

// 				p := math.Exp(-a * (A1 + beta*A2) / weight)
// 				strength := endTimeJ - startTimeJ

// 				edge := InferEdge{startTimeI, endTimeI, startTimeJ, endTimeJ, p, strength}
// 				edgeList = append(edgeList, edge)
// 				// subEdgeDictstartTimeJ] = edge
// 			}
// 		}

// 		// edgeDict startTimeI] = subEdgeDict
// 	}

// 	// return edgeDict
// 	return edgeList
// }

// // GetEventTimeRanges
// func GetEventTimeRanges(eventsList [][]int, filename string) [][]int {
// 	timeRangesList := [][]int{}

// 	var length int
// 	if filename == "congestion" || filename == "congestion2" {
// 		length = 8785
// 	} else if filename == "flow" {
// 		length = 7221
// 	} else {
// 		length = 10000
// 	}

// 	for _, events := range eventsList {
// 		timeRanges := make([]int, length)

// 		movingIndex := 0
// 		for index, t := range events {
// 			// 1,2,3,6,7,8  ,11,12,13,14,30
// 			// 0,1,2,3,4,5  ,6, 7, 8, 9, 10
// 			if index < movingIndex {
// 				continue
// 			}

// 			movingIndex = index
// 			for (events[movingIndex] - t) == (movingIndex - index) {
// 				movingIndex++
// 				if movingIndex >= len(events) {
// 					break
// 				}
// 			}

// 			timeRanges[t] = events[movingIndex-1] + 1
// 		}

// 		timeRangesList = append(timeRangesList, timeRanges)
// 	}

// 	return timeRangesList
// }

// // GetEdgeRevertIndex ..
// func GetEdgeRevertIndex(edges map[int]map[int]InferEdges, locationIDs []int) map[int]map[int][]RevertEdgeIndexInfo {
// 	var edgeRevertIndex map[int]map[int][]RevertEdgeIndexInfo
// 	edgeRevertIndex = make(map[int]map[int][]RevertEdgeIndexInfo)

// 	for _, locationIDi := range locationIDs {
// 		for _, locationIDj := range locationIDs {
// 			inferEdgesIJ := edges[locationIDi][locationIDj]

// 			for index, inferEdge := range inferEdgesIJ {
// 				startTimeJ := inferEdge.StartTimeJ

// 				_, ok := edgeRevertIndex[startTimeJ]
// 				if ok {
// 					_, okok := edgeRevertIndex[startTimeJ][locationIDj]
// 					if !okok {
// 						edgeRevertIndex[startTimeJ][locationIDj] = []RevertEdgeIndexInfo{}
// 					}
// 					edgeRevertIndex[startTimeJ][locationIDj] = append(edgeRevertIndex[startTimeJ][locationIDj], RevertEdgeIndexInfo{locationIDi, index})
// 				} else {
// 					var subMap map[int][]RevertEdgeIndexInfo
// 					subMap = make(map[int][]RevertEdgeIndexInfo)
// 					subMap[locationIDj] = []RevertEdgeIndexInfo{}
// 					subMap[locationIDj] = append(subMap[locationIDj], RevertEdgeIndexInfo{locationIDi, index})
// 					edgeRevertIndex[startTimeJ] = subMap
// 				}
// 			}
// 		}
// 	}

// 	return edgeRevertIndex
// }

// // NetworkInfer ..
// func NetworkInfer(locationIDs []int, locations []Location, timeRangesList [][]int, filename string) (graph.Graph, map[int]map[int]InferEdges) {
// 	// initialize the distance matrix
// 	distanceMatrix := GetDistanceMatrix(locationIDs, locations)

// 	// initialize edges between any pair of roads
// 	// egdes: []InferEdge
// 	edges := GetEdges(locationIDs, timeRangesList, distanceMatrix, filename)

// 	// map[startTimeJ][locationIDj] => {locationIDi, index}
// 	edgeRevertIndex := GetEdgeRevertIndex(edges, locationIDs)

// 	// k is the number of edges in cascading network
// 	k := math.Min(float64(len(locationIDs)), 10)

// 	var resultG graph.Graph

// 	for resultG.CountEdge() < int(k) {
// 		maxDeltaIJ := 0.0
// 		maxIID := -1
// 		maxJID := -1

// 		for _, locationIDi := range locationIDs {
// 			for _, locationIDj := range locationIDs {
// 				// each pair of i_j here represents the pair of P
// 				ni := graph.Node{locationIDi}
// 				nj := graph.Node{locationIDj}
// 				if resultG.IsEdgeExist(&ni, &nj) || resultG.IsEdgeExist(&nj, &ni) {
// 					continue
// 				}
// 				// each pair of i_j here represents the pair of P\G

// 				deltaIJ := 0.0

// 				for _, edge := range edges[locationIDi][locationIDj] {
// 					wcij := edge.P
// 					startTimeJ := edge.StartTimeJ

// 					deltaCIJ := 0.0

// 					for _, info := range edgeRevertIndex[startTimeJ][locationIDj] {
// 						// info 	CausalLocationID int
// 						//        EdgeIndex    int
// 						roadMID := info.CausalLocationID
// 						nm := graph.Node{roadMID}

// 						if roadMID == locationIDi || !resultG.IsNodeExist(&nm) {
// 							continue
// 						}

// 						edgeIndex := info.EdgeIndex

// 						wcmj := edges[roadMID][locationIDj][edgeIndex].P

// 						deltaCIJ += wcmj
// 					}

// 					if deltaCIJ == 0 {
// 						deltaIJ += math.Log(1 + (wcij / 0.01))
// 					} else {
// 						deltaIJ += math.Log(1 + (wcij / deltaCIJ))
// 					}
// 				}

// 				// debug

// 				if deltaIJ > maxDeltaIJ {
// 					maxDeltaIJ = deltaIJ
// 					maxIID = locationIDi
// 					maxJID = locationIDj
// 				}
// 			}
// 		}

// 		fmt.Println(maxIID, maxJID, maxDeltaIJ)

// 		ni := graph.Node{maxIID}
// 		nj := graph.Node{maxJID}
// 		resultG.AddNode(&ni)
// 		resultG.AddNode(&nj)
// 		resultG.AddEdge(&ni, &nj)
// 	}

// 	return resultG, edges
// }
