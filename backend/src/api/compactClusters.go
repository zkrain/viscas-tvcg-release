package api

import (
	"math"
)

// type CompactEventsTimeSegment struct {
// 	NEvent                 int   `json:"nEvent"`
// 	NTimeRange             int   `json:"nTimeRange"`
// 	OrderedTimeRangeLength []int `json:"orderedTimeRangeLength"`
// }

// type CompactEvents = []CompactEventsTimeSegment

// func ContextualizeClustersCompactly(clusters [][]int, eventDataDict EventDataDict, filename string) [][]CompactEvents {
// 	var length int
// 	if filename == "congestion" {
// 		length = 8785
// 	} else if filename == "flow" {
// 		length = 7221
// 	} else {
// 		length = 10000
// 	}

// 	nBin := 60
// 	interval := length/nBin + 1
// 	compactClusters := [][]CompactEvents{}

// 	for _, cluster := range clusters {
// 		compactEventsList := []CompactEvents{}

// 		for _, locationID := range cluster {
// 			compactEvents := make([]CompactEventsTimeSegment, nBin)
// 			events := eventDataDict[locationID].Events
// 			timeRanges := getEventTimeRanges(events, length)

// 			for startTime, endTime := range timeRanges {
// 				if endTime == 0 {
// 					continue
// 				}
// 				bin := startTime / interval
// 				l := endTime - startTime
// 				compactEvents[bin].NEvent++
// 				compactEvents[bin].NTimeRange += l
// 				compactEvents[bin].OrderedTimeRangeLength = append(compactEvents[bin].OrderedTimeRangeLength, l)
// 			}
// 			compactEventsList = append(compactEventsList, compactEvents)
// 		}
// 		compactClusters = append(compactClusters, compactEventsList)
// 	}

// 	return compactClusters
// }

type CompactEventsTimeSegment = [3]float64

type CompactEvents = []CompactEventsTimeSegment

type ListWithVariance struct {
	List     []float64
	Variance []float64
}

type ClusterInformation struct {
	// StartTimeNumberInformation ListWithVariance
	// RawEventNumberInformation  ListWithVariance
	// LengthInformation          ListWithVariance
	StartTimeNumberMatrix [][]float64
	RawEventNumberMatrix  [][]float64
	LengthMatrix          [][]float64
}

type SpatialDistributionOfClusters struct {
	MaxDistance      float64
	MinDistance      float64
	DistributionList [][]int
}

func ContextualizeClustersCompactly(clusters [][]int, eventDataDict EventDataDict, filename string) []ClusterInformation {
	var length int
	if filename == "congestion" || filename == "congestion2" {
		length = 8785
	} else if filename == "flow" {
		length = 7221
	} else {
		length = 10000
	}

	nBin := 120
	interval := length/nBin + 1
	clusterInformationList := []ClusterInformation{}

	for _, cluster := range clusters {
		lengthMatrix := [][]float64{}
		startTimeNumberMatrix := [][]float64{}
		rawEventNumberMatrix := [][]float64{}

		for _, locationID := range cluster {
			lengthVector := []float64{}
			startTimeNUmberVector := []float64{}
			rawEventNumberVector := []float64{}

			compactEvents := make([]CompactEventsTimeSegment, nBin)
			events := eventDataDict[locationID].Events
			timeRanges := getEventTimeRanges(events, length)

			for startTime, endTime := range timeRanges {
				if endTime == 0 {
					continue
				}
				bin := startTime / interval
				l := endTime - startTime
				compactEvents[bin][0] += float64(1)
				compactEvents[bin][1] += float64(l)
				compactEvents[bin][2] += float64(l * l)
			}

			for bin, _ := range compactEvents {
				if compactEvents[bin][1] == 0 {
					compactEvents[bin][2] = 0
				} else {
					compactEvents[bin][2] = compactEvents[bin][2] / compactEvents[bin][1]
				}

				lengthVector = append(lengthVector, compactEvents[bin][2])
				startTimeNUmberVector = append(startTimeNUmberVector, compactEvents[bin][0])
				rawEventNumberVector = append(rawEventNumberVector, compactEvents[bin][1])
			}

			lengthMatrix = append(lengthMatrix, lengthVector)
			startTimeNumberMatrix = append(startTimeNumberMatrix, startTimeNUmberVector)
			rawEventNumberMatrix = append(rawEventNumberMatrix, rawEventNumberVector)
		}

		clusterInformation := ClusterInformation{startTimeNumberMatrix, rawEventNumberMatrix, lengthMatrix}
		clusterInformationList = append(clusterInformationList, clusterInformation)
	}

	return clusterInformationList
}

func average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

func getsd(nums []float64) float64 {
	var sum float64
	var mean float64
	var sd float64
	lengthFloat64 := float64(len(nums))
	for _, num := range nums {
		sum += num
	}
	mean = sum / lengthFloat64

	for _, num := range nums {
		sd += math.Pow(float64(num)-mean, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	sd = math.Sqrt(sd / lengthFloat64)

	return sd
}

// GetEventTimeRanges
func getEventTimeRanges(events []int, length int) []int {
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

	return timeRanges
}

func GetSpatialDistributionOfClusters(clusters [][]int, locationDict map[int]Location) SpatialDistributionOfClusters {
	nBin := 10.0
	maxDistance := 0.0
	minDistance := 10000.0

	distancesList := [][]float64{}

	distributions := [][]int{}

	for _, locationIDs := range clusters {

		distances := []float64{}
		for i, locationIDi := range locationIDs {
			for j, locationIDj := range locationIDs {
				if i >= j {
					continue
				}

				locationI := locationDict[locationIDi]
				locationJ := locationDict[locationIDj]
				d := Get_Geo_Distance(float64(locationI.Lat), float64(locationI.Lng), float64(locationJ.Lat), float64(locationJ.Lng))

				distances = append(distances, d)
				if d > maxDistance {
					maxDistance = d
				}
				if d < minDistance {
					minDistance = d
				}
			}
		}

		distancesList = append(distancesList, distances)
	}

	interval := (maxDistance - minDistance) / nBin

	for i := range clusters {
		distanceDistribution := make([]int, int(nBin))
		distances := distancesList[i] // distances of the current cluster
		for _, d := range distances {
			bin := int(math.Floor((d - minDistance) / interval))
			if bin >= int(nBin) {
				bin = bin - 1
			}
			distanceDistribution[bin]++
		}
		distributions = append(distributions, distanceDistribution)
	}

	return SpatialDistributionOfClusters{minDistance, maxDistance, distributions}
}

func GetClusterAverageEventCount(clusters [][]int, eventDataDict EventDataDict) []float32 {
	var clusterAverageEventCount []float32
	for _, cluster := range clusters{
		var iSum int = 0
		for _, rid := range cluster {
			iSum += eventDataDict[rid].Count
		}
		clusterAverageEventCount = append(clusterAverageEventCount, float32(iSum)/float32(len(cluster)))
	}
	return clusterAverageEventCount
}

func GetClusterEventCount(clusters [][]int, eventDataDict EventDataDict) [][]int {
	var clusterEventCount [][]int
	for _, cluster := range clusters{
		var counts []int
		for _, rid := range cluster {
			counts = append(counts, eventDataDict[rid].Count)
		}
		clusterEventCount = append(clusterEventCount, counts)
	}
	return clusterEventCount
}

func GetClusterInfectEvent(clusters [][]int, infectEventDataDict InfectEventDataDict, deltaTime int) [][][]InfectEvent {
	var clusterInfectEvent [][][]InfectEvent
	for _, cluster := range clusters{
		var infectEvents [][]InfectEvent
		for _, rid := range cluster{
			infectEvent := infectEventDataDict[rid]
			infectEvents = append(infectEvents, infectEvent)
		}
		clusterInfectEvent = append(clusterInfectEvent, infectEvents)
	}
	return clusterInfectEvent
}
