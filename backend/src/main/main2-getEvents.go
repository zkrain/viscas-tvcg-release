package main

import (
	"bufio"
	"congest/parameters"
	"congest/readings"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func ParseRoads(trajFileName string, f *os.File) {
	file, err := os.Open(trajFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), ",")

		roadSegmentId := params[0]
		lv, _ := strconv.Atoi(params[1])
		l, _ := strconv.ParseFloat(params[2], 32)

		if lv < parameters.CONGESTED_LIMITED_SPEED_THRESHOLD || l < parameters.ROAD_LEN_THRESHOLD_IN_GETEVENTS {
			continue
		}
		// timespan := len(params[2:])

		congestedTimestamps := []int{}

		for t, vStr := range params[3:] {
			if len(vStr) > 0 { // has record => not free flow
				v, err := strconv.Atoi(vStr)
				checkError(err)
				if v <= parameters.CONGESTED_SPEED_THRESHOLD && (float32(v) <= float32(lv)*0.5) {
					congestedTimestamps = append(congestedTimestamps, t)
				}
			}
		}

		f.WriteString(roadSegmentId + "," + strconv.Itoa(len(congestedTimestamps)) + "," + params[2] + "#")
		f.WriteString(arrayToString(congestedTimestamps, ","))
		f.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() { // congestion event detection
	outputFileName := "output/congestionEvents-zj.txt"
	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer outputFile.Close()

	root := "output"
	fileNames := readings.FileNames(root)
	for _, fileName := range fileNames[1:] {
		if string(fileName[7:9]) != "v_" {
			continue
		}
		fmt.Println("fileName", fileName)

		ParseRoads(fileName, outputFile)
	}
}
