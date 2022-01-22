package api

import (
	"fmt"
	"math"
)

func Get_Geo_Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344

	return dist
}

func Cluster_Avg(currentCluster []int, embedPointDict map[int]ProjPoint) [2]float64 {
	ret := [2]float64{0.0, 0.0}
	l := len(currentCluster)
	for _, rid := range currentCluster {
		embedPoint := embedPointDict[rid]
		ret[0] += embedPoint.Left
		ret[1] += embedPoint.Top
	}
	ret[0] = ret[0] / float64(l)
	ret[1] = ret[1] / float64(l)
	return ret
}

func Retrieve_Neighbors(d Location, D []Location, EpsT float64, EpsS float64, embedPointDict map[int]ProjPoint) []Location {
	ret := []Location{}

	t := [2]float64{embedPointDict[d.Rid].Left, embedPointDict[d.Rid].Top}

	for _, dOfD := range D {
		if dOfD.Rid != d.Rid {
			sd := Get_Geo_Distance(float64(d.Lat), float64(d.Lng), float64(dOfD.Lat), float64(dOfD.Lng)) // km

			t2 := [2]float64{embedPointDict[dOfD.Rid].Left, embedPointDict[dOfD.Rid].Top}
			td := math.Sqrt((t[0]-t2[0])*(t[0]-t2[0]) + (t[1]-t2[1])*(t[1]-t2[1]))

			if sd <= EpsS && td <= EpsT {
				ret = append(ret, dOfD)
			}
		}
	}
	return ret
}

func Stdbscan(D []Location, EpsT float64, EpsS float64, MinPts int, deltaE float64, embedPointDict map[int]ProjPoint) ([][]int, []int) {

	fmt.Println(len(D), " locations come when conducting clutering")

	clusters := [][]int{}
	currentCluster := []int{}
	currentClusterLabel := 1
	noise := []int{}
	var clusterMap map[int]int
	clusterMap = make(map[int]int)
	Q := []Location{}

	for _, d := range D {
		if _, ok := clusterMap[d.Rid]; !ok {

			X := Retrieve_Neighbors(d, D, EpsT, EpsS, embedPointDict)

			if len(X) < MinPts {
				clusterMap[d.Rid] = -1
				noise = append(noise, d.Rid)
			} else {
				currentClusterLabel++
				if len(currentCluster) > 0 {
					clusters = append(clusters, currentCluster)
				}
				currentCluster = []int{}

				for _, dOfX := range X {
					currentCluster = append(currentCluster, dOfX.Rid)
					clusterMap[dOfX.Rid] = currentClusterLabel
					Q = append(Q, dOfX)
				}

				for len(Q) > 0 {
					currentObj := Q[len(Q)-1]
					Q = Q[0 : len(Q)-1]
					Y := Retrieve_Neighbors(currentObj, D, EpsT, EpsS, embedPointDict)

					if len(Y) >= MinPts {
						for _, dOfY := range Y {

							_, okok := clusterMap[dOfY.Rid]
							if (clusterMap[dOfY.Rid] != -1) && (!okok) {
								avgXY := Cluster_Avg(currentCluster, embedPointDict)

								valueXY := []float64{embedPointDict[dOfY.Rid].Left, embedPointDict[dOfY.Rid].Top} // o.Value
								xd := avgXY[0] - valueXY[0]
								yd := avgXY[1] - valueXY[1]

								if math.Sqrt(xd*xd+yd*yd) <= deltaE {
									fmt.Println("???")
									currentCluster = append(currentCluster, dOfY.Rid)
									clusterMap[dOfY.Rid] = currentClusterLabel
									Q = append(Q, dOfY)
								}
							}

						}
					}
				}

			}
		}
	}

	return clusters, noise
}
