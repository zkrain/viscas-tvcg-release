package api

import (
	"congest/graph"
	// "congest/api"
	"fmt"
	"math"
	"sync"
	"strconv"
	"strings"
	"sort"

	geo "github.com/kellydunn/golang-geo"
)

// Influence ..
type Influence struct {
	StartTimeI int
	EndTimeI   int
	StartTimeJ int
	EndTimeJ   int
	P          float64
	W          float64
	// Strength   int
}

// nodeIKey, nodeJKey, probabilityStr, durationI, durationJ
type FullInfluenceInfo struct {
	LocationIDi int
	StartTimeI int
	DurationI int
	LocationIDj int
	StartTimeJ int
	DurationJ int
	P float64
	OtherPotentialCauses []Cause
}

type Cause struct {
	LocationID int
	StartTime int
	Duration int
	P float64
}

type DeletedInfluenceDict map[string]bool // string = locationIDi/j, startTime

type Influences = []Influence

type RevertEdgeIndexInfo struct {
	ParentLocationID int //
	EdgeIndex        int
}

type KeyPair = [5]string
type PropagationTree struct {
	TreeStartTime int
	TreeEndTime int
	FullInfluenceInfoList []FullInfluenceInfo
}
func (t *PropagationTree) SetStartTimeIfValid(time int) {
	if t.TreeStartTime > time || t.TreeStartTime == 0 {
		t.TreeStartTime = time
	}
}
func (t *PropagationTree) SetEndTimeIfValid(time int) {
	if t.TreeEndTime < time || t.TreeEndTime == 0 {
		t.TreeEndTime = time
	}
}
func (t *PropagationTree) AddInfluenceInTree(fullInfluenceInfo FullInfluenceInfo) {
	t.FullInfluenceInfoList = append(t.FullInfluenceInfoList, fullInfluenceInfo)
}
type PropagationTrees []PropagationTree
func (ts PropagationTrees) Len() int { return len(ts) } // 获取此 slice 的长度
func (ts PropagationTrees) Less(i, j int) bool { // 根据元素的TreeStartTime升序排序
    return ts[i].TreeStartTime < ts[j].TreeStartTime
}
func (ts PropagationTrees) Swap(i, j int) { ts[i], ts[j] = ts[j], ts[i] } // 交换数据

type GraphRelation struct {
	Parent []string // "locationID,startTime,duration"
	Children []string // "locationID,startTime,duration"
	ParentP []float64
	ChildrenP []float64
	lock      sync.RWMutex
}

func (myGraphRelation *GraphRelation) AddParent(parentKey string, p float64) {
	myGraphRelation.lock.Lock()
	myGraphRelation.Parent = append(myGraphRelation.Parent, parentKey)
	myGraphRelation.ParentP = append(myGraphRelation.ParentP, p)
	myGraphRelation.lock.Unlock()
}

func (myGraphRelation *GraphRelation) AddChildren(childKey string, p float64) {
	myGraphRelation.lock.Lock()
	myGraphRelation.Children = append(myGraphRelation.Children, childKey)
	myGraphRelation.ChildrenP = append(myGraphRelation.ChildrenP, p)
	myGraphRelation.lock.Unlock()
}


// CalculateDistanceBetweenLocations return in km
func CalculateDistanceBetweenLocations(locationI Location, locationJ Location) float64 {
	var pointI *geo.Point
	var pointJ *geo.Point
	pointI = geo.NewPoint(float64(locationI.Lat), float64(locationI.Lng))
	pointJ = geo.NewPoint(float64(locationJ.Lat), float64(locationJ.Lng))
	d := pointJ.GreatCircleDistance(pointI)
	return d
}

// GetDistanceMatrix return a distance dictionary like map[int]map[int]float64
func GetDistanceMatrix(locationIDs []int, locations []Location) map[int]map[int]float64 {
	var distanceMatrix map[int]map[int]float64
	distanceMatrix = make(map[int]map[int]float64)

	maxValidDistance := 0.0
	minValidDistance := 10000.0

	// 问题在于，输入进来的地点之间的距离能否确保不会太大？
	for i, locationIDi := range locationIDs {
		var row map[int]float64
		row = make(map[int]float64)

		for j, locationIDj := range locationIDs {
			if i != j {
				distance := CalculateDistanceBetweenLocations(locations[i], locations[j])

				if distance > maxValidDistance {
					maxValidDistance = distance
				}
				if distance < minValidDistance {
					minValidDistance = distance
				}

				row[locationIDj] = distance
			}
		}

		distanceMatrix[locationIDi] = row
	}

	for locationIDi, row := range distanceMatrix {
		for locationIDj, d := range row {
			distanceMatrix[locationIDi][locationIDj] = (d - minValidDistance) / (maxValidDistance - minValidDistance)
		}
	}

	return distanceMatrix
}

// GetInfluences ..
func GetInfluences(locationIDs []int, tw int, infectEventsList [][]InfectEvent, distanceMatrix map[int]map[int]float64, basicEpss []float64, timeSpanLength int) map[int]map[int]Influences {
	var edges map[int]map[int]Influences
	edges = make(map[int]map[int]Influences)
	// type Influence struct {
	// 	StartTimeI int
	// 	EndTimeI   int
	// 	StartTimeJ int
	// 	EndTimeJ   int
	// 	P          float64
	// 	Strength   int
	// }

	for i, locationIDi := range locationIDs {
		infectEventsOfLi := infectEventsList[i]
		for j, locationIDj := range locationIDs {
			if i == j {
				continue
			}

			basicEpsI := basicEpss[i]
			basicEpsJ := basicEpss[j]

			infectEventsOfLj := infectEventsList[j]
			distance := distanceMatrix[locationIDj][locationIDi] // km

			var influences []Influence = ModelInfluences(infectEventsOfLi, infectEventsOfLj, tw, distance, basicEpsI, basicEpsJ, timeSpanLength)

			if _, ok := edges[locationIDi]; ok {
				edges[locationIDi][locationIDj] = influences
			} else {
				var subInfluencesDict map[int]Influences
				subInfluencesDict = make(map[int]Influences)
				subInfluencesDict[locationIDj] = influences
				edges[locationIDi] = subInfluencesDict
			}
		}
	}

	return edges
}

// ModelInfluences ...
func ModelInfluences(infectEventsOfLi []InfectEvent, infectEventsOfLj []InfectEvent, tw int, distance float64, basicEpsI float64, basicEpsJ float64, timeSpanLength int) []Influence {
	// inputs
		// infectEventsOfLi: casual location
		// infectEventsOfLj: location that is affected

	// store with slice/list
	influences := []Influence{}

	// constants
	MAX_INFLUENCE_TIME := float64(tw)
	ALPHA := 1.0/3.0
	LAMBDA := 3.0
	// beta := 18.0
	// gamma := 1.0


	// // 以i作为出发点的
	// flagJ := 0 // a flag indicating the current valid infectEvent.
	// for _, infectEventI := range infectEventsOfLi {
	// 	startTimeI := infectEventI.StartTime
	// 	durationI := infectEventI.DurationTime
	// 	endTimeI := infectEventI.EndTime

	// 	for j := flagJ; j < len(infectEventsOfLj); j++ {
	// 		infectEventJ := infectEventsOfLj[j]

	// 		if infectEventJ.startTime - startTimeI > MAX_INFLUENCE_TIME {
	// 			// 如果j太晚发生，当前i的查找可以结束了
	// 			break;
	// 		}

	// 		startTimeJ = infectEventJ.StartTime
	// 		endTimeJ = infectEventJ.EndTime

	// 		if startTimeJ < startTimeI {
	// 			// 如果这个j已经早于当前i了，那自然这个j也会早于下一个i
	// 			flagJ++
	// 		}

	// 		A1 := float64(startTimeJ-startTimeI) / MAX_INFLUENCE_TIME 	// normalized
	// 		A2 := distance  																						// have been normalized
	// 		p := math.Exp(-a * (A1 + LAMBDA*A2))

	// 		influence := Influence{startTimeI, endTimeI, startTimeJ, endTimeJ, p, strength}
	// 		influences = append(influences, influence)
	// 		// 如果已经找到一个了，说明这个是
	// 		break
	// 	}
	// }

	// 以j作为出发点的
	flagI := 0 // a flag indicating the current valid infectEvent.
	for _, infectEventJ := range infectEventsOfLj {
		startTimeJ := infectEventJ.StartTime
		durationJ := infectEventJ.DurationTime
		endTimeJ := infectEventJ.EndTime

		for i := flagI; i < len(infectEventsOfLi); i++ {
			infectEventI := infectEventsOfLi[i]

			if infectEventI.StartTime >= startTimeJ {
				// 如果i比j还要晚发生，当前i的查找可以结束了
				break;
			}

			startTimeI := infectEventI.StartTime
			durationI := infectEventI.DurationTime
			endTimeI := infectEventI.EndTime

			if startTimeJ - startTimeI >= int(MAX_INFLUENCE_TIME) {
				// 如果这个i远远早于当前j，那肯定也会远远早于下一个j
				flagI++
				continue
			}

			// 判断当前这个i是不是最后一个晚于j的i
			if (i + 1) < len(infectEventsOfLi) {
				nextInfectEventI := infectEventsOfLi[i+1]
				if nextInfectEventI.StartTime < startTimeJ { // 如果下一个i仍早于j,继续
					continue
				} else { // 如果当前i是最后一个晚于j的i，break
					A1 := float64(startTimeJ-startTimeI) / MAX_INFLUENCE_TIME 	// normalized
					A2 := distance  																						// have been normalized
					p := math.Exp(-ALPHA * (A1 + LAMBDA*A2))
					eps := math.Pow(basicEpsI, float64(durationI)) * math.Pow(basicEpsJ, float64(durationJ))
					w := p/eps

					influence := Influence{startTimeI, endTimeI, startTimeJ, endTimeJ, p, w}
					influences = append(influences, influence)
					break
				}
			} else { // 已经找到最后了
				A1 := float64(startTimeJ-startTimeI) / MAX_INFLUENCE_TIME 	// normalized
				A2 := distance  																						// have been normalized
				p := math.Exp(-ALPHA * (A1 + LAMBDA*A2))
				eps := math.Pow(basicEpsI, float64(durationI)) * math.Pow(basicEpsJ, float64(durationJ))
				w := p/eps

				influence := Influence{startTimeI, endTimeI, startTimeJ, endTimeJ, p, w}
				influences = append(influences, influence)
			}
		}
	}

	return influences
}

// GetEventTimeRanges
func GetEventTimeRanges(eventsList [][]int, filename string) [][]int {
	infectEventsList := [][]int{}

	var length int
	if filename == "congestion" || filename == "congestion2" {
		length = 8785
	} else if filename == "flow" {
		length = 7221
	} else {
		length = 10000
	}

	for _, events := range eventsList {
		timeRanges := make([]int, length)

		movingIndex := 0
		for index, t := range events {
			// 1,2,3,6,7,8  ,11,12,13,14,30
			// 0,1,2,3,4,5  ,6, 7, 8, 9, 10
			if index < movingIndex {
				continue
			}

			movingIndex = index
			for (events[movingIndex] - t) == (movingIndex - index) {
				movingIndex++
				if movingIndex >= len(events) {
					break
				}
			}

			timeRanges[t] = events[movingIndex-1] + 1
		}

		infectEventsList = append(infectEventsList, timeRanges)
	}

	return infectEventsList
}

// GetEdgeRevertIndex ..
func GetEdgeRevertIndex(influences map[int]map[int]Influences, locationIDs []int) map[int]map[int][]RevertEdgeIndexInfo {
	var edgeRevertIndex map[int]map[int][]RevertEdgeIndexInfo
	edgeRevertIndex = make(map[int]map[int][]RevertEdgeIndexInfo)

	for _, locationIDi := range locationIDs {
		for _, locationIDj := range locationIDs {
			influencesIJ := influences[locationIDi][locationIDj]

			for index, influence := range influencesIJ {
				startTimeJ := influence.StartTimeJ

				_, ok := edgeRevertIndex[startTimeJ]
				if ok {
					_, okok := edgeRevertIndex[startTimeJ][locationIDj]
					if !okok {
						edgeRevertIndex[startTimeJ][locationIDj] = []RevertEdgeIndexInfo{}
					}
					edgeRevertIndex[startTimeJ][locationIDj] = append(edgeRevertIndex[startTimeJ][locationIDj], RevertEdgeIndexInfo{locationIDi, index})
				} else {
					var subMap map[int][]RevertEdgeIndexInfo
					subMap = make(map[int][]RevertEdgeIndexInfo)
					subMap[locationIDj] = []RevertEdgeIndexInfo{}
					subMap[locationIDj] = append(subMap[locationIDj], RevertEdgeIndexInfo{locationIDi, index})
					edgeRevertIndex[startTimeJ] = subMap
				}
			}
		}
	}

	return edgeRevertIndex
}

func visitInfluence(
	locationIDi int,
	startTimeI int,
	durationI string,
	locationIDj int,
	startTimeJ int,
	durationJ string,
	probability float64,
	propagationTree *PropagationTree,
	relationDict map[string]*GraphRelation,
	visitedInfluenceDict *map[string]int,
	edgeRevertIndex *map[int]map[int][]RevertEdgeIndexInfo,
	influencesDict *map[int]map[int]Influences) {

	nodeIKey := fmt.Sprintf("%d,%d", locationIDi, startTimeI)
	nodeJKey := fmt.Sprintf("%d,%d", locationIDj, startTimeJ)
	influenceKey := nodeIKey + "_" + nodeJKey
	_, exist := (*visitedInfluenceDict)[influenceKey]
	if exist {
		return
	} else {
		(*visitedInfluenceDict)[influenceKey] = 1
	}

	// get other potential influences
	// edgeRevertIndex[startTimeJ][locationIDj] // to do
	otherPotentialCauses := []Cause{}
	for _, info := range (*edgeRevertIndex)[startTimeJ][locationIDj] {
		// info
		// ParentLocationID int
		// EdgeIndex    		int
		locationIDm := info.ParentLocationID
		if locationIDi == locationIDm {
			continue
		}
		edgeIndex := info.EdgeIndex

		potentialInfluence := (*influencesDict)[locationIDm][locationIDj][edgeIndex]
		startTimeM := potentialInfluence.StartTimeI
		endTimeM := potentialInfluence.EndTimeI
		durationM := endTimeM - startTimeM
		p := potentialInfluence.P
		otherPotentialCauses = append(otherPotentialCauses, Cause{locationIDm, startTimeM, durationM, p})
	}

	// *resultKeyPairs = append(*resultKeyPairs, [2]string{nodeIKey, nodeJKey})
	(*propagationTree).SetStartTimeIfValid(startTimeI)
	durationIInt, _ := strconv.ParseInt(durationI, 10, 64)
	durationJInt, _ := strconv.ParseInt(durationJ, 10, 64)
	(*propagationTree).SetEndTimeIfValid(startTimeJ+int(durationJInt))
	// (*propagationTree).AddKeyPair([5]string{nodeIKey, nodeJKey, probabilityStr, durationI, durationJ})
	(*propagationTree).AddInfluenceInTree(FullInfluenceInfo{locationIDi, startTimeI, int(durationIInt), locationIDj, startTimeJ, int(durationJInt), probability, otherPotentialCauses})

	relationKeyI := nodeIKey + "," + durationI
	relationKeyJ := nodeJKey + "," + durationJ

	graphRelationI, oki := relationDict[relationKeyI]
	graphRelationJ, okj := relationDict[relationKeyJ]

	if oki {
		for piIndex, pikey := range graphRelationI.Parent { // pi: "localLocationIDi,localStartTimeI"
			strs := strings.Split(pikey, ",")
			localLocationIDi, _ := strconv.ParseInt(strs[0], 10, 64)
			localStartTimeI, _ := strconv.ParseInt(strs[1], 10, 64)
			localDurationI := strs[2]
			p := graphRelationI.ParentP[piIndex]
			visitInfluence(int(localLocationIDi), int(localStartTimeI), localDurationI, locationIDi, startTimeI, durationI, p, propagationTree, relationDict, visitedInfluenceDict, edgeRevertIndex, influencesDict)
		}
	}
	if okj {
		for pjIndex, pjkey := range graphRelationJ.Parent { // pi: "localLocationIDi,localStartTimeI"
			strs := strings.Split(pjkey, ",")
			localLocationIDi, _ := strconv.ParseInt(strs[0], 10, 64)
			localStartTimeI, _ := strconv.ParseInt(strs[1], 10, 64)
			localDurationI := strs[2]
			p := graphRelationJ.ParentP[pjIndex]
			visitInfluence(int(localLocationIDi), int(localStartTimeI), localDurationI, locationIDj, startTimeJ, durationJ, p, propagationTree, relationDict, visitedInfluenceDict, edgeRevertIndex, influencesDict)
		}
	}
	if oki {
		for ciIndex, cikey := range graphRelationI.Children { // pi: "localLocationIDj,localStartTimeJ"
			strs := strings.Split(cikey, ",")
			localLocationIDj, _ := strconv.ParseInt(strs[0], 10, 64)
			localStartTimeJ, _ := strconv.ParseInt(strs[1], 10, 64)
			localDurationJ := strs[2]
			p := graphRelationI.ChildrenP[ciIndex]
			visitInfluence(locationIDi, startTimeI, durationI, int(localLocationIDj), int(localStartTimeJ), localDurationJ, p, propagationTree, relationDict, visitedInfluenceDict, edgeRevertIndex, influencesDict)
		}
	}
	if okj {
		for cjIndex, cjkey := range graphRelationJ.Children { // pi: "localLocationIDj,localStartTimeJ"
			strs := strings.Split(cjkey, ",")
			localLocationIDj, _ := strconv.ParseInt(strs[0], 10, 64)
			localStartTimeJ, _ := strconv.ParseInt(strs[1], 10, 64)
			localDurationJ := strs[2]
			p := graphRelationJ.ChildrenP[cjIndex]
			visitInfluence(locationIDj, startTimeJ, durationJ, int(localLocationIDj), int(localStartTimeJ), localDurationJ, p, propagationTree, relationDict, visitedInfluenceDict, edgeRevertIndex, influencesDict)
		}
	}
	return
}

func GetPropagationTree(influencesDict *map[int]map[int]Influences, links [][2]int, edgeRevertIndex *map[int]map[int][]RevertEdgeIndexInfo) PropagationTrees {
	// initialize visition dict
	visitedInfluenceDict := make(map[string]int) // visitedInfluenceDict["2,32_3,35"][1] 表示 2,32->3,35的influence是否被访问
	// initialized relation dict
	relationDict := make(map[string]*GraphRelation)

	for _, link := range links {
		locationIDi := link[0]
		locationIDj := link[1]

		// // construct visition dict
		// linkStr := fmt.Sprintf("%d,%d", locationIDi, locationIDj)
		// subDict := make(map[int]int)
		// visitedInfluenceDict[linkStr] = subDict

		// construct relation dict
		influences := (*influencesDict)[locationIDi][locationIDj]
		for _, influence := range influences {
			startTimeI := influence.StartTimeI
			durationI := influence.EndTimeI - startTimeI
			startTimeJ := influence.StartTimeJ
			durationJ := influence.EndTimeJ - startTimeJ
			p := influence.P

			// 如果有任何一个节点的duration太低，则不考虑
			if durationI < 1 || durationJ < 1 {
				continue
			}

			keyStrI := fmt.Sprintf("%d,%d,%d", locationIDi, startTimeI, durationI)
			keyStrJ := fmt.Sprintf("%d,%d,%d", locationIDj, startTimeJ, durationJ)

			_, okI := relationDict[keyStrI]
			if okI {
				relationDict[keyStrI].AddChildren(keyStrJ, p)
			} else {
				var graphRelation GraphRelation
				graphRelation.AddChildren(keyStrJ, p)
				relationDict[keyStrI] = &graphRelation
			}

			_, okJ := relationDict[keyStrJ]
			if okJ {
				relationDict[keyStrJ].AddParent(keyStrI, p)
			} else {
				var graphRelation GraphRelation
				graphRelation.AddParent(keyStrI, p)
				relationDict[keyStrJ] = &graphRelation
			}
		}
	}

	// get trees
	var trees PropagationTrees
	for _, link := range links {
		locationIDi := link[0]
		locationIDj := link[1]

		influences := (*influencesDict)[locationIDi][locationIDj]
		for _, influence := range influences {
			var propagationTree PropagationTree
			p := influence.P
			durationI := influence.EndTimeI - influence.StartTimeI
			durationJ := influence.EndTimeJ - influence.StartTimeJ
			if durationI < 1 || durationJ < 1 {
				continue
			}
			durationIStr := fmt.Sprintf("%d", durationI)
			durationJStr := fmt.Sprintf("%d", durationJ)
			visitInfluence(locationIDi, influence.StartTimeI, durationIStr, locationIDj, influence.StartTimeJ, durationJStr, p, &propagationTree, relationDict, &visitedInfluenceDict, edgeRevertIndex, influencesDict)
			// DFS

			if len(propagationTree.FullInfluenceInfoList) > 0 {
				trees = append(trees, propagationTree)
			}
		}
	}

	return trees
}


func NetworkInfer(
	locationIDs []int,
	k int,
	tw int,
	locations []Location,
	infectEventsList [][]InfectEvent,
	basicEpss []float64,
	timeSpanLength int) ([][2]int, map[int]map[int]Influences, []PropagationTree, map[string]Influences, map[string]Influences) {
	// the order of infectEventsList is consistent with that in locations

	// initialize the distance matrix
	distanceMatrix := GetDistanceMatrix(locationIDs, locations) // map[int]map[int]float64

	// initialize edges between any pair of roads
	influences := GetInfluences(locationIDs, tw, infectEventsList, distanceMatrix, basicEpss, timeSpanLength) // map[int]map[int]Influences

	// map[startTimeJ][locationIDj] => []{locationIDi, index}
	edgeRevertIndex := GetEdgeRevertIndex(influences, locationIDs)

	// k is the number of edges in cascading network

	// start inference process
	var cascadingPattern graph.Graph
	var linkList [][2]int
	for cascadingPattern.CountEdge() < int(k) {
		maxDeltaIJ := 0.0
		maxIID := -1
		maxJID := -1

		for _, locationIDi := range locationIDs {
			for _, locationIDj := range locationIDs {
				// each pair of i_j here represents the pair of P
				ni := graph.Node{locationIDi}
				nj := graph.Node{locationIDj}
				if cascadingPattern.IsEdgeExist(&ni, &nj) || cascadingPattern.IsEdgeExist(&nj, &ni) {
					continue
				}
				// each pair of i_j here represents the pair of P\G

				deltaIJ := 0.0
				for _, influence := range influences[locationIDi][locationIDj] {
					wcij := influence.W
					if (wcij > float64(10E+30)) {
						wcij = float64(10E+30)
					}
					startTimeJ := influence.StartTimeJ

					deltaijm := 0.0

					for _, info := range edgeRevertIndex[startTimeJ][locationIDj] {
						// info
						// ParentLocationID int
						// EdgeIndex    		int
						locationIDm := info.ParentLocationID
						edgeIndex := info.EdgeIndex

						nm := graph.Node{locationIDm}
						if locationIDm == locationIDi || !cascadingPattern.IsNodeExist(&nm) { // 如果m等于i，或者nm不在g里面，continue
							continue
						}

						wcmj := influences[locationIDm][locationIDj][edgeIndex].W
						if (wcmj > float64(10E+30)) {
							wcmj = float64(10E+30)
						}

						deltaijm += wcmj
					}

					if deltaijm == 0 {
						deltaIJ += math.Log(wcij)
					} else {
						deltaIJ += math.Log(1 + (wcij / deltaijm))
					}
				}

				// pick up the edge with maximum margin gain
				if deltaIJ > maxDeltaIJ {
					maxDeltaIJ = deltaIJ
					maxIID = locationIDi
					maxJID = locationIDj
				}
			}
		}

		fmt.Println(maxIID, maxJID, maxDeltaIJ)

		ni := graph.Node{maxIID}
		nj := graph.Node{maxJID}
		cascadingPattern.AddNode(&ni)
		cascadingPattern.AddNode(&nj)
		cascadingPattern.AddEdge(&ni, &nj)
		linkList = append(linkList, [2]int{maxIID, maxJID})
	}

	// get the edges/links in the cascading pattern
	// links := cascadingPattern.GetLinkList()

	// get the influences contained in the pattern
	influencesInG := make(map[int]map[int]Influences)

	for _, link := range linkList {
		locationIDi := link[0]
		locationIDj := link[1]

		var subInfluencesDict map[int]Influences
		subInfluencesDict = make(map[int]Influences)

		if _, ok := influencesInG[locationIDi]; ok {
			influencesInG[locationIDi][locationIDj] = influences[locationIDi][locationIDj]
		} else {
			subInfluencesDict[locationIDj] = influences[locationIDi][locationIDj]
			influencesInG[locationIDi] = subInfluencesDict
		}
	}

	// get which infection events have been considered in G
	infectionEventsDictInG := make(map[string]map[string]bool)
	for _, link := range linkList {
		locationIDi := link[0]
		locationIDj := link[1]

		deletedInfluences := make(map[string]bool)
		influencesIJ := influencesInG[locationIDi][locationIDj]
		for _, influence := range influencesIJ {
			startTimeI := influence.StartTimeI
			startTimeJ := influence.StartTimeJ
			nodeIKey := fmt.Sprintf("%d,%d", locationIDi, startTimeI)
			nodeJKey := fmt.Sprintf("%d,%d", locationIDj, startTimeJ)
			deletedInfluences[nodeIKey] = true
			deletedInfluences[nodeJKey] = true
		}

		IJStr := fmt.Sprintf("%d,%d", locationIDi, locationIDj)
		infectionEventsDictInG[IJStr] = deletedInfluences
	}

	// get influences that are contrary to the result and in G
	influencesDictContraryInG := make(map[string]Influences)
	// get influences that are contrary to the result and not in G
	influencesDictContraryNotInG := make(map[string]Influences)
	for _, link := range linkList {
		locationIDi := link[0]
		locationIDj := link[1]

		IJStr := fmt.Sprintf("%d,%d", locationIDi, locationIDj)
		JIStr := fmt.Sprintf("%d,%d", locationIDj, locationIDi)

		influencesContraryInG := []Influence{}
		influencesContraryNotInG := []Influence{}
		for _, influence := range influences[locationIDj][locationIDi] {
			startTimeI := influence.StartTimeI
			nodeIKey := fmt.Sprintf("%d,%d", locationIDi, startTimeI)

			startTimeJ := influence.StartTimeJ
			nodeJKey := fmt.Sprintf("%d,%d", locationIDj, startTimeJ)

			_, oki := infectionEventsDictInG[IJStr][nodeIKey]
			_, okj := infectionEventsDictInG[IJStr][nodeJKey]

			if oki || okj {
				influencesContraryInG = append(influencesContraryInG, influence)
			} else {
				influencesContraryNotInG = append(influencesContraryNotInG, influence)
			}
		}

		influencesDictContraryInG[JIStr] = influencesContraryInG
		influencesDictContraryNotInG[JIStr] = influencesContraryNotInG
	}

	// get propagation trees
	trees := GetPropagationTree(&influences, linkList, &edgeRevertIndex)
	fmt.Println("number of trees", len(trees))
	sort.Sort(trees)

	fmt.Println(linkList)

	return linkList, influencesInG, trees, influencesDictContraryInG, influencesDictContraryNotInG
}