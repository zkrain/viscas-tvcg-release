package api

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type RoadShape struct {
	Shapes string
	Len    float32
	Lv     int
	ID     int
	Dir    int
	Level  int
}

func RoadShapeQueryer() map[int]RoadShape {
	var m map[int]RoadShape
	m = make(map[int]RoadShape)

	file, err := os.Open("../output/roadsShape.txt") // todo 要改回来
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), "#")
		attributes := strings.Split(params[0], ",")
		shapes := params[1]

		rid, _ := strconv.Atoi(attributes[0])
		len, _ := strconv.ParseFloat(attributes[1], 32)
		lv, _ := strconv.Atoi(attributes[2])
		dir, _ := strconv.Atoi(attributes[3])
		level, _ := strconv.Atoi(attributes[4])
		m[rid] = RoadShape{shapes, float32(len), lv, rid, dir, level}
	}

	return m
}
