package main

import (
	"api"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var roads api.Roads

	file, err := os.Open("output/congestionEvents.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), "#")
		attributes := strings.Split(params[0], ",")

		count, _ := strconv.Atoi(attributes[1])
		if count < 8000 {
			continue
		}
		rid, _ := strconv.Atoi(attributes[0])
		roads = append(roads, api.CongestedRoad{rid, count})
	}

	api.GetTopKCongestedRoads(30, &roads)
}
