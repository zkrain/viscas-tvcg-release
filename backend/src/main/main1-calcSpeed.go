package main

import (
	"congest/readings"
	"fmt"
)

func main() { // calculate speeed of each road segment
	for i := 0; i < 244234; i += 10000 {
		fmt.Println(i)
		network := readings.Network(i, i+10000)
		network.OutputRoadShape()
		root := "../data/traj_hz_20160304/20160304"
		fileNames := readings.FileNames(root)
		for j, fileName := range fileNames[1:] {
			if j%100 == 0 {
				fmt.Println(i, j, len(fileNames))
			}
			network.ParseTraj(fileName)
		}
		network.OutputVelocities()
	}
}
