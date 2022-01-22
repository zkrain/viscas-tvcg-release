package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	// "sort"
	"strconv"
	"strings"
)

// overwrite congestionQuery.go

type EventsWithCountID struct {
	Count  int
	ID     int
	Events []int
}

type InfectEvent struct {
	StartTime    int `json:"startTime"`
	DurationTime int `json:"durationTime"`
	EndTime      int `json:"endTime"`
}

type EventDataDict map[int]EventsWithCountID
type InfectEventData []InfectEvent
type InfectEventDataDict map[int]InfectEventData
type AverageEventDataDict map[int][][]float32

type IDsFilter struct {
	EventDataDict EventDataDict
	InfectEventDataDict InfectEventDataDict
	AverageEventDataDict AverageEventDataDict
}

// FilterIDsByCount 根据count来过滤
func (myIDsFilter *IDsFilter) FilterIDsByCount(vLow int, vHigh int) []int {
	// initially filter
	IDs := []int{}
	for _, EventsWithCountID := range myIDsFilter.EventDataDict {
		if EventsWithCountID.Count > vLow && EventsWithCountID.Count < vHigh {
			IDs = append(IDs, EventsWithCountID.ID)
		}
	}

	fmt.Println(len(IDs), " roads left")
	return IDs
}

// GetAllID 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
func (myIDsFilter *IDsFilter) GetAllID() []int {
	allIDs := make([]int, 0, len(myIDsFilter.EventDataDict))
	for k := range myIDsFilter.EventDataDict {
		allIDs = append(allIDs, k)
	}
	return allIDs
}

func (myIDsFilter *IDsFilter) GetInfectEventsIDict(ids []int) map[int][]InfectEvent {
	var infectEventsDict map[int][]InfectEvent
	infectEventsDict = make(map[int][]InfectEvent)
	for _, k := range ids {
		infectEventsDict[k] = myIDsFilter.InfectEventDataDict[k]
	}
	return infectEventsDict
}

// GetDistribution 获取拥堵时间数量的分布
func (myIDsFilter *IDsFilter) GetDistribution() ([]int, int, int) {
	nBin := 20
	eventNumbers := []int{}
	// 这里的Count指的是拥堵的时间段的个数，不是出租车的记录
	maxCount := 0
	minCount := 9000
	for _, EventsWithCountID := range myIDsFilter.EventDataDict {
		eventNumbers = append(eventNumbers, EventsWithCountID.Count)
		if EventsWithCountID.Count > maxCount {
			maxCount = EventsWithCountID.Count
		}
		if EventsWithCountID.Count < minCount {
			minCount = EventsWithCountID.Count
		}
	}

	var distribution = make([]int, nBin)
	interval := math.Ceil(float64(maxCount-minCount) / float64(nBin))

	fmt.Println("calculate distribution", minCount, maxCount, interval)

	for _, n := range eventNumbers {
		bin := math.Floor(float64(n-minCount) / interval)
		if int(bin) >= nBin {
			bin--
		}
		distribution[int(bin)]++
	}

	return distribution, maxCount, minCount
}

// InitialIDsFilter generate an instanfe of IDsFilter
func InitialIDsFilter(filename string) IDsFilter {
	var EventDataDict map[int]EventsWithCountID
	var InfectEventDataDict map[int]InfectEventData
	EventDataDict = make(map[int]EventsWithCountID)
	InfectEventDataDict = make(map[int]InfectEventData)

	file, err := os.Open("../output/" + filename + "Data/" + filename + "Events.txt")
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
		// l, _ := strconv.ParseFloat(attributes[2], 32)

		EventsStr := strings.Split(params[1], ",")
		Events := []int{}
		for _, t := range EventsStr {
			it, _ := strconv.Atoi(t)
			Events = append(Events, it)
		}

		EventDataDict[rid] = EventsWithCountID{count, rid, Events}
		InfectEventDataDict[rid] = GetInfectEvent(Events, 1)
	}
	var AverageEventDataDict AverageEventDataDict
	var outputFileName string
	if filename == "air"{
		outputFileName = "../output/airData/airEventsAverages.json"
	} else if filename == "congestion" || filename == "congestion2"{
		outputFileName = "../output/congestion2Data/congestion2Averages.json"
	} else if filename == "flow" {
		outputFileName = "../output/flowData/flowAverages.json"
	}
	fileAve, errAve := os.Open(outputFileName)
	if errAve != nil {
		log.Fatal(err)
	}
	defer fileAve.Close()
	byteValue, _ := ioutil.ReadAll(fileAve)
	json.Unmarshal(byteValue, &AverageEventDataDict)
	myIDsFilter := IDsFilter{EventDataDict, InfectEventDataDict, AverageEventDataDict}
	return myIDsFilter
}

// 将event从int[]格式转为{start,duration,end}[]格式
func GetInfectEvent(event []int, deltaTime int) []InfectEvent {
	var infectEvent []InfectEvent
	if len(event) == 0 {
		return infectEvent
	}
	startTime := event[0]
	endTime := event[0]
	for i := 1; i < len(event); i++ {
		if event[i]-event[i-1] > deltaTime+1 {
			infectEvent = append(infectEvent, InfectEvent{startTime,
				endTime - startTime + 1, endTime})
			startTime = event[i]
		}
		endTime = event[i]
		if i == len(event)-1 {
			infectEvent = append(infectEvent, InfectEvent{startTime,
				endTime - startTime + 1, endTime})
		}
	}
	return infectEvent
}
