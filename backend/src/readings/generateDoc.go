package readings

import (
	"bufio"
	"congest/api"
	"congest/parameters"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateDoc(roadShapeQueryer map[int]api.RoadShape, trajFileName string, outputFile *os.File) int {
	trajFile, err := os.Open(trajFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer trajFile.Close()

	count := 0

	scanner := bufio.NewScanner(trajFile)
	sentenceSlice := []int{}

	for scanner.Scan() {
		params := strings.Split(scanner.Text(), ",")
		if len(params[0]) != 1 {
			rid, _ := strconv.Atoi(params[2])
			if _, ok := roadShapeQueryer[rid]; !ok { // if the road segment is outside the target region
				continue
			}

			t1, err := time.Parse(RFC3339FullDate, params[0])
			t2, err := time.Parse(RFC3339FullDate, params[1])
			checkError(err)
			ts := t2.Sub(t1).Seconds()
			l := roadShapeQueryer[rid].Len
			level := roadShapeQueryer[rid].Level
			if ts <= 2 || l < parameters.ROAD_LEN_THRESHOLD || level == 4 {
				continue
			}

			sentenceSlice = append(sentenceSlice, rid)
			count++

		} else { // new trip
			if len(sentenceSlice) > 5 { // save sentence
				outputFile.WriteString(arrayToString(sentenceSlice, " "))
				outputFile.WriteString("\n")
			}

			sentenceSlice = []int{}
		}
	}
	return count
}
