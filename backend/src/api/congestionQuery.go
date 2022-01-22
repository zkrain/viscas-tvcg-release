package api

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type CongestedRoad struct {
	Count         int
	ID            int
	Len           float32
	CongestEvents []int
}

type Roads []CongestedRoad
type RoadDict map[int]CongestedRoad

// type RoadsFilter struct {
// 	Roads Roads
// }
type RoadsFilter struct {
	Roads RoadDict
}

func (R Roads) Len() int { return len(R) }
func (R Roads) Less(i, j int) bool {
	return R[i].Count > R[j].Count
}
func (R Roads) Swap(i, j int) { R[i], R[j] = R[j], R[i] }

// func (roadsFilter *RoadsFilter) GetTopKCongestedRoads(k int) []int {
// 	// initially filter
// 	var validRoads Roads

// 	// debug
// 	blindCountThreshold := 0
// 	blindCountThreshold2 := 230000
// 	for _, road := range roadsFilter.Roads {
// 		if road.Count > blindCountThreshold && road.Count < blindCountThreshold2 {
// 			validRoads = append(validRoads, road)
// 		}
// 	}

// 	fmt.Println(len(validRoads), " roads left")
// 	sort.Sort(validRoads)

// 	roadIds := []int{}
// 	for _, road := range validRoads {
// 		roadIds = append(roadIds, road.ID)
// 	}
// 	return roadIds
// }

func (roadsFilter *RoadsFilter) FilterRoadsByCongestCount(vLow int, vHigh int) []int {
	// initially filter
	roadIds := []int{}
	for _, road := range roadsFilter.Roads {
		if road.Count > vLow && road.Count < vHigh {
			roadIds = append(roadIds, road.ID)
		}
	}

	fmt.Println(len(roadIds), " roads left")
	return roadIds
}

// GetAllRoadsID 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
func (roadsFilter *RoadsFilter) GetAllRoadsID() []int {
	keys := make([]int, 0, len(roadsFilter.Roads))
	for k := range roadsFilter.Roads {
		keys = append(keys, k)
	}
	return keys
}

// GetCongestionDistribution 获取拥堵时间数量的分布
func (roadsFilter *RoadsFilter) GetCongestionDistribution() ([]int, int, int) {
	nBin := 30
	congestedEventNumbers := []int{}
	// 这里的Count指的是拥堵的时间段的个数，不是出租车的记录
	maxCount := 0
	minCount := 8000
	for _, congestedRoad := range roadsFilter.Roads {
		congestedEventNumbers = append(congestedEventNumbers, congestedRoad.Count)
		if congestedRoad.Count > maxCount {
			maxCount = congestedRoad.Count
		}
		if congestedRoad.Count < minCount {
			minCount = congestedRoad.Count
		}
	}

	var distribution = make([]int, nBin)
	interval := math.Ceil(float64(maxCount-minCount) / float64(nBin))

	fmt.Println(minCount, maxCount, interval)

	for _, n := range congestedEventNumbers {
		bin := math.Floor(float64(n-minCount) / interval)
		if int(bin) >= nBin {
			bin--
		}
		distribution[int(bin)]++
	}

	return distribution, maxCount, minCount
}

func (roadsFilter *RoadsFilter) GetTopKCongestedRoadsWithRegion(k int, validRoadDict map[int]bool) []int {
	// initially filter
	var validRoads Roads
	for _, road := range roadsFilter.Roads {
		if _, ok := validRoadDict[road.ID]; ok {
			validRoads = append(validRoads, road)
		}
	}

	roadIds := []int{}

	fmt.Println(len(validRoads), " roads left")
	sort.Sort(validRoads)
	if len(validRoads) == 0 {
		return roadIds
	}

	for _, road := range validRoads[:k] {
		roadIds = append(roadIds, road.ID)
		// fmt.Println(road.Count, road.ID)
	}
	return roadIds
}

// func

// func (roadsFilter *RoadsFilter) CorrelationSearch(targetRoadId int, roads []RoadShape) { // candidate roads
// 	targetRoadCongestEvents := roadsFilter.Roads[targetRoadId].CongestEvents
// 	for _, road := range roads {
// 		rid := road.ID
// 		congestEvents := roadsFilter.Roads[rid].CongestEvents

// 	}
// }

func InitialRoadsFilter() RoadsFilter {
	var roads map[int]CongestedRoad
	roads = make(map[int]CongestedRoad)

	file, err := os.Open("../output/congestionData/congestionEvents.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), "#")

		attributes := strings.Split(params[0], ",")
		count, _ := strconv.Atoi(attributes[1])
		rid, _ := strconv.Atoi(attributes[0])
		l, _ := strconv.ParseFloat(attributes[2], 32)

		CongestEventsStr := strings.Split(params[1], ",")
		CongestEvents := []int{}
		for _, t := range CongestEventsStr {
			it, _ := strconv.Atoi(t)
			CongestEvents = append(CongestEvents, it)
		}

		roads[rid] = CongestedRoad{count, rid, float32(l), CongestEvents}
		// roads = append(roads, CongestedRoad{count, rid, float32(l), CongestEvents})
	}

	roadFilter := RoadsFilter{roads}
	return roadFilter
}
