package main

import (
	"bufio"
	"congest/api"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var maxLng = 120.29291546053952
var maxLat = 30.354600411193275
var minLng = 119.96724016669839
var minLat = 30.170034989228213

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func check(localMaxLat float64, localMaxLng float64, localMinLat float64, localMinLng float64) bool {

	if minLng < localMinLng && localMinLng < maxLng && minLat < localMaxLat && localMaxLat < maxLat {
		return true
	}

	if minLng < localMaxLng && localMaxLng < maxLng && minLat < localMaxLat && localMaxLat < maxLat {
		return true
	}

	if minLng < localMinLng && localMinLng < maxLng && minLat < localMinLat && localMinLat < maxLat {
		return true
	}

	if minLng < localMaxLng && localMaxLng < maxLng && minLat < localMinLat && localMinLat < maxLat {
		return true
	}

	return false
}

func filterRoadShape() map[int]bool {
	var roadDict map[int]bool
	roadDict = make(map[int]bool)

	outputFileName := "../output/roadsShape-zj-filter.txt"
	f, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer f.Close()

	roadShapeQueryer := api.RoadShapeQueryer() // modify the file of api.RoadShapeQueryer

	// for rid, value := range roadShapeQueryer {
	for rid := 0; rid < 250000; rid++ {
		if _, ok := roadShapeQueryer[rid]; !ok {
			continue
		}
		value := roadShapeQueryer[rid]

		localMaxLng := 0.0
		localMaxLat := 0.0
		localMinLng := 200.0
		localMinLat := 2000.0

		shapes := string(value.Shapes)
		last := len(shapes) - 1
		latlngs := strings.Split(shapes[:last], ";")
		for _, latlngStr := range latlngs {
			// fmt.Println(latlngStr)
			latAndlng := strings.Split(latlngStr, ",")
			lat, err := strconv.ParseFloat(latAndlng[0], 32)
			lng, err := strconv.ParseFloat(latAndlng[1], 32)
			checkError(err)

			localMaxLng = math.Max(lng, localMaxLng)
			localMaxLat = math.Max(lat, localMaxLat)
			localMinLng = math.Min(lng, localMinLng)
			localMinLat = math.Min(lat, localMinLat)
		}

		flag := check(localMaxLat, localMaxLng, localMinLat, localMinLng)
		if flag { // the current road is in the target region
			ridStr := strconv.Itoa(rid)
			lenStr := strconv.FormatFloat(float64(value.Len), 'f', 2, 64)
			lvStr := strconv.Itoa(value.Lv)
			dirStr := strconv.Itoa(value.Dir)
			levelStr := strconv.Itoa(value.Level)

			if value.Level == 4 {
				continue
			}

			roadDict[rid] = true

			fmt.Println(ridStr)
			f.WriteString(ridStr + "," + lenStr + "," + lvStr + "," + dirStr + "," + levelStr + "#")
			f.WriteString(value.Shapes)
			f.WriteString("\n")
		}
	}
	return roadDict
}

func filterCongestionEvents(roadDict map[int]bool) {
	outputFileName := "../output/congestionEvents-zj-filter.txt"
	f, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer f.Close()

	file, err := os.Open("../output/congestionEvents.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), "#")
		attributes := strings.Split(params[0], ",")

		rid, _ := strconv.Atoi(attributes[0])
		eventCount, _ := strconv.Atoi(attributes[1])
		if eventCount > 2000 || eventCount < 100 {
			continue
		}
		if _, ok := roadDict[rid]; !ok {
			continue
		}

		f.WriteString(scanner.Text())
		f.WriteString("\n")
	}
}

func main() {
	roadDict := filterRoadShape()
	filterCongestionEvents(roadDict)
}
