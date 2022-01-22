package api

import (
	"congest/parameters"
	"strconv"
	"strings"

	geo "github.com/kellydunn/golang-geo"
)

func queryRoadsUsingLatlng(targetRoadID int, lat float64, lng float64, RoadShapeQueryer map[int]RoadShape) []RoadShape { // current version fails to consider the direction of the roads.
	pTarget := geo.NewPoint(lat, lng)
	retRoads := []RoadShape{}

	for _, road := range RoadShapeQueryer {
		if road.ID == targetRoadID {
			continue
		}
		roadShape := road.Shapes
		last := len(roadShape) - 1
		latlngsStr := strings.Split(roadShape[:last], ";")

		for _, latlngStr := range latlngsStr {
			latAndLng := strings.Split(latlngStr, ",")
			rLat, _ := strconv.ParseFloat(latAndLng[0], 32)
			rLng, _ := strconv.ParseFloat(latAndLng[1], 32)
			d := pTarget.GreatCircleDistance(geo.NewPoint(rLat, rLng)) * 1000
			if d > parameters.NEIGHBORING_ROADS_DISTANCE_L {
				break
			}
			if d < parameters.NEIGHBORING_ROADS_DISTANCE_S {
				retRoads = append(retRoads, road)
				break
			}
		}
	}

	return retRoads
}

func NeighboringRoadsQuery(targetRoadID int, k int, RoadShapeQueryer map[int]RoadShape) ([]RoadShape, []RoadShape) {
	targetRoad := RoadShapeQueryer[targetRoadID]
	targetRoadShape := targetRoad.Shapes
	last := len(targetRoadShape) - 1
	latlngsStr := strings.Split(targetRoadShape[:last], ";")

	startLatlngStr := latlngsStr[0]
	endLatlngStr := latlngsStr[len(latlngsStr)-1]

	startLatAndLng := strings.Split(startLatlngStr, ",")
	endLatAndLng := strings.Split(endLatlngStr, ",")

	startLat, _ := strconv.ParseFloat(startLatAndLng[0], 32)
	startLng, _ := strconv.ParseFloat(startLatAndLng[1], 32)

	endLat, _ := strconv.ParseFloat(endLatAndLng[0], 32)
	endLng, _ := strconv.ParseFloat(endLatAndLng[1], 32)

	neighboringRoadsFromStartNodes := queryRoadsUsingLatlng(targetRoadID, startLat, startLng, RoadShapeQueryer)
	neighboringRoadsFromEndNodes := queryRoadsUsingLatlng(targetRoadID, endLat, endLng, RoadShapeQueryer)
	return neighboringRoadsFromStartNodes, neighboringRoadsFromEndNodes
}
