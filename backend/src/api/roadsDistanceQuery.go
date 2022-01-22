package api

import (
	"strconv"
	"strings"

	geo "github.com/kellydunn/golang-geo"
)

// CalculateDistance calculate the distance between two road: from the end of uRoad to the start of dRoad
func CalculateDistance(uRoad RoadShape, dRoad RoadShape) float64 {
	// uRoad: upstream road
	// dRoad: downstream road
	// The vehicles fisrt pass uRoad and then dRoad !!! important

	// dir: 1-bidirection, 2-forward, 3-backward, 4-noway

	if uRoad.Dir >= 4 || dRoad.Dir >= 4 {
		return 1000.0
	}

	// started node stored in shape string
	uLast := len(uRoad.Shapes) - 1
	uLatlngsStr := strings.Split(uRoad.Shapes[:uLast], ";")

	uRoadNode1Str := uLatlngsStr[0]
	uRoadNode2Str := uLatlngsStr[len(uLatlngsStr)-1]
	uRoadNode1LatAndLng := strings.Split(uRoadNode1Str, ",")
	uRoadNode2LatAndLng := strings.Split(uRoadNode2Str, ",")

	uRoadNode1Lat, _ := strconv.ParseFloat(uRoadNode1LatAndLng[0], 32)
	uRoadNode1Lng, _ := strconv.ParseFloat(uRoadNode1LatAndLng[1], 32)
	uRoadNode2Lat, _ := strconv.ParseFloat(uRoadNode2LatAndLng[0], 32)
	uRoadNode2Lng, _ := strconv.ParseFloat(uRoadNode2LatAndLng[1], 32)

	// started node stored in shape string
	dLast := len(dRoad.Shapes) - 1
	dLatlngsStr := strings.Split(dRoad.Shapes[:dLast], ";")

	dRoadNode1Str := dLatlngsStr[0]
	dRoadNode2Str := dLatlngsStr[len(dLatlngsStr)-1]
	dRoadNode1LatAndLng := strings.Split(dRoadNode1Str, ",")
	dRoadNode2LatAndLng := strings.Split(dRoadNode2Str, ",")

	dRoadNode1Lat, _ := strconv.ParseFloat(dRoadNode1LatAndLng[0], 32)
	dRoadNode1Lng, _ := strconv.ParseFloat(dRoadNode1LatAndLng[1], 32)
	dRoadNode2Lat, _ := strconv.ParseFloat(dRoadNode2LatAndLng[0], 32)
	dRoadNode2Lng, _ := strconv.ParseFloat(dRoadNode2LatAndLng[1], 32)

	var uPoint *geo.Point
	var dPoint *geo.Point

	// fmt.Println(uRoad.Dir, dRoad.Dir)

	if uRoad.Dir == 1 {
		uPoint = geo.NewPoint((uRoadNode1Lat+uRoadNode2Lat)/2, (uRoadNode1Lng+uRoadNode2Lng)/2)
	} else if uRoad.Dir == 2 {
		uPoint = geo.NewPoint(uRoadNode2Lat, uRoadNode2Lng)
	} else if uRoad.Dir == 3 {
		uPoint = geo.NewPoint(uRoadNode1Lat, uRoadNode1Lng)
	}

	if dRoad.Dir == 1 {
		dPoint = geo.NewPoint((dRoadNode1Lat+dRoadNode2Lat)/2, (dRoadNode1Lng+dRoadNode2Lng)/2)
	} else if dRoad.Dir == 2 {
		dPoint = geo.NewPoint(dRoadNode1Lat, dRoadNode1Lng)
	} else if dRoad.Dir == 3 {
		dPoint = geo.NewPoint(dRoadNode2Lat, dRoadNode2Lng)
	}
	// fmt.Println(dPoint.Lat(), dPoint.Lng(), uPoint.Lat(), uPoint.Lng())

	d := dPoint.GreatCircleDistance(uPoint)
	return d
}
