package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Embedding struct {
	Rid int     `json:"rid"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
}

type ProjPoint struct {
	Rid  int     `json:"rid"`
	Left float64 `json:"left"`
	Top  float64 `json:"top"`
}

type EmbeddingQuery struct {
	AllProjPoints []ProjPoint
	ProjPointDict map[int]ProjPoint
}

func InitializeEmbeddingQuery(filename string) EmbeddingQuery {
	outputFileName := "../output/" + filename + "Data/coordinates.json"
	file, err := os.Open(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	var data []Embedding
	json.Unmarshal(byteValue, &data)

	points := []ProjPoint{}
	pointDict := make(map[int]ProjPoint)

	maxX := -90000.0
	maxY := -90000.0
	minX := 90000.0
	minY := 90000.0
	paddingRatio := 0.05
	rangeX := 0.0
	rangeY := 0.0

	for _, d := range data {
		if d.X > maxX {
			maxX = d.X
		}
		if d.Y > maxY {
			maxY = d.Y
		}
		if d.X < minX {
			minX = d.X
		}
		if d.Y < minY {
			minY = d.Y
		}

		rangeX = maxX - minX
		rangeY = maxY - minY
	}

	for _, d := range data {
		letf := (d.X - minX + paddingRatio*rangeX) / rangeX / (paddingRatio*2.0 + 1.0)
		top := (d.Y - minY + paddingRatio*rangeY) / rangeY / (paddingRatio*2.0 + 1.0)
		newPoint := ProjPoint{d.Rid, letf, top}
		points = append(points, newPoint)
		pointDict[d.Rid] = newPoint
	}

	myEmbeddingQuery := EmbeddingQuery{points, pointDict}
	return myEmbeddingQuery
}
