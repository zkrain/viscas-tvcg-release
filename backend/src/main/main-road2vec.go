package main

import (
	"api"
	"fmt"
	"log"
	"os"
	"readings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	roadShapeQueryer := api.RoadShapeQueryer()
	fmt.Println(len(roadShapeQueryer))

	outputFileName := "output/document2.txt"
	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer outputFile.Close()

	count := 0

	root := "../data/traj_hz_20160304/20160304"
	fileNames := readings.FileNames(root)
	for j, fileName := range fileNames[1:] {
		if j%100 == 0 {
			fmt.Println(j, len(fileNames), count)
		}
		c := readings.GenerateDoc(roadShapeQueryer, fileName, outputFile)
		count += c
	}
}
