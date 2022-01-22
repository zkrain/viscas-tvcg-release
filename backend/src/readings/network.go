package readings

import (
	"bufio"
	"congest/parameters"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	coordTransform "github.com/qichengzx/coordtransform"
)

var referenceTime, err = time.Parse(RFC3339FullDate, "2016/03/01 00:00:00")

const RFC3339FullDate = "2006/01/02 15:04:05"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getTimeBin(t time.Time) int {
	diffMinutes := t.Sub(referenceTime).Minutes()
	return int(diffMinutes / 10)
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

type Intersection struct {
	ID  int
	Lat float32
	Lng float32
}

type VelocityRecord struct {
	Count int
	SumV  float32
	MeanV float32
}

type RoadSegment struct {
	ID    int
	Len   float32
	Vs    [8785]VelocityRecord // Vs have a length with 8785
	lv    int                  // limited speed
	Shape [][2]float32
	Dir   int
	Level int
}

// func (roadSegment *RoadSegment) addV(v float32) []float32 {
// 	roadSegment.Vs = append(roadSegment.Vs, v)
// 	return roadSegment.Vs
// }

func (roadSegment *RoadSegment) addVR(v float32, timeBin int) {
	if float32(v) > (float32(roadSegment.lv) * 1.5) { // speed is abnormal
		return
	}

	oldVR := roadSegment.Vs[timeBin]
	newCount := oldVR.Count + 1
	newSumV := oldVR.SumV + v
	newVR := VelocityRecord{newCount, newSumV, 0}
	roadSegment.Vs[timeBin] = newVR
}

func (roadSegment *RoadSegment) addShape(latStr string, lngStr string) [][2]float32 {
	lat, err := strconv.ParseFloat(latStr, 32)
	lng, err := strconv.ParseFloat(lngStr, 32)
	lat, err = strconv.ParseFloat(fmt.Sprintf("%.6f", lat), 32)
	lng, err = strconv.ParseFloat(fmt.Sprintf("%.6f", lng), 32)
	checkError(err)

	lng, lat = coordTransform.GCJ02toWGS84(lng, lat)

	latlng := [2]float32{}
	latlng[0] = float32(lat)
	latlng[1] = float32(lng)

	roadSegment.Shape = append(roadSegment.Shape, latlng)
	return roadSegment.Shape
}

type RoadNetwork struct {
	Intersections []Intersection
	RoadSegments  []RoadSegment
	ridLower      int
	ridUpper      int
}

func (roadNetwork *RoadNetwork) addIntersection(intersection Intersection) []Intersection {
	roadNetwork.Intersections = append(roadNetwork.Intersections, intersection)
	return roadNetwork.Intersections
}

func (roadNetwork *RoadNetwork) addRoadSegment(roadSegment RoadSegment) []RoadSegment {
	roadNetwork.RoadSegments = append(roadNetwork.RoadSegments, roadSegment)
	return roadNetwork.RoadSegments
}

func (roadNetwork *RoadNetwork) ParseTraj(trajFileName string) { // debug id: 163353
	// count163353 := 0

	file, err := os.Open(trajFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		params := strings.Split(scanner.Text(), ",")
		if len(params[0]) != 1 {
			rid, err := strconv.Atoi(params[2])
			checkError(err)
			if rid >= roadNetwork.ridUpper || roadNetwork.ridLower > rid {
				// ridLower == rid 时，应该计算
				continue
			}
			t1, err := time.Parse(RFC3339FullDate, params[0])
			t2, err := time.Parse(RFC3339FullDate, params[1])
			checkError(err)

			len := roadNetwork.RoadSegments[rid-roadNetwork.ridLower].Len
			// lv := roadNetwork.RoadSegments[rid-roadNetwork.ridLower].lv
			ts := t2.Sub(t1).Seconds()
			if ts <= 1 || len < parameters.ROAD_LEN_THRESHOLD {
				continue
			}

			medianTime := t1.Add(t2.Sub(t1) / 2)
			timeBin := getTimeBin(medianTime)
			v := 3.6 * len / float32(ts)
			// if rid == 163353 { // debug
			// 	count163353++
			// 	fmt.Println(rid, len, lv, ts, count163353, v, scanner.Text())
			// }
			roadNetwork.RoadSegments[rid-roadNetwork.ridLower].addVR(v, timeBin)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (roadNetwork *RoadNetwork) OutputRoadShape() {
	// different from OutputVelocities, this function only generetes one file
	outputFileName := "output/roadsShape-zj.txt"
	f, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	defer f.Close()

	for i, roadSegment := range roadNetwork.RoadSegments {
		if roadSegment.Len < parameters.ROAD_LEN_THRESHOLD {
			continue
		}
		if i%1000 == 0 {
			fmt.Println(i, len(roadNetwork.RoadSegments))
		}

		lenStr := strconv.FormatFloat(float64(roadSegment.Len), 'f', 2, 64)
		lvStr := strconv.Itoa(roadSegment.lv)
		ridStr := strconv.Itoa(roadSegment.ID)
		dirStr := strconv.Itoa(roadSegment.Dir)
		levelStr := strconv.Itoa(roadSegment.Level)

		f.WriteString(ridStr + "," + lenStr + "," + lvStr + "," + dirStr + "," + levelStr + "#")
		for _, latlng := range roadSegment.Shape {
			f.WriteString(fmt.Sprintf("%f", latlng[0]))
			f.WriteString(",")
			f.WriteString(fmt.Sprintf("%f", latlng[1]))
			f.WriteString(";")
		}
		f.WriteString("\n")
	}
}

func (roadNetwork *RoadNetwork) OutputVelocities() {
	ridUpperStr := strconv.Itoa(roadNetwork.ridUpper)
	ridLowerStr := strconv.Itoa(roadNetwork.ridLower)
	outputFileName := "output/v_10_zj_" + ridLowerStr + "_" + ridUpperStr + ".txt"
	f, _ := os.Create(outputFileName)
	defer f.Close()

	for _, roadSegment := range roadNetwork.RoadSegments {
		if roadSegment.Len < parameters.ROAD_LEN_THRESHOLD {
			continue
		}

		vList := []int{}
		validCountForRoad := 0

		for _, vr := range roadSegment.Vs {
			var v int
			if vr.Count > 0 {
				validCountForRoad += vr.Count
				v = int(vr.SumV / float32(vr.Count))
			} else {
				v = roadSegment.lv
			}
			vList = append(vList, v)
		}
		if validCountForRoad < parameters.ROAD_RECORD_COUNT_THRESHOLD {
			continue
		}

		lenStr := strconv.FormatFloat(float64(roadSegment.Len), 'f', 6, 64)
		lvStr := strconv.Itoa(roadSegment.lv)
		roadVelocitiesStr := strings.ReplaceAll(arrayToString(vList, ","), lvStr, "")
		ridStr := strconv.Itoa(roadSegment.ID)
		f.WriteString(ridStr + "," + lvStr + "," + lenStr + ",")
		f.WriteString(roadVelocitiesStr)
		f.WriteString("\n")
	}
}

// 183749 244234
// 244234 每1w一批？

func Network(ridLower int, ridUpper int) RoadNetwork {
	file, err := os.Open("../data/Road_Network_HZ_2016Q1.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	lineCount := 0
	rid := -1
	intersections := []Intersection{}
	roadSegments := []RoadSegment{}
	roadnetwork := RoadNetwork{intersections, roadSegments, ridLower, ridUpper}

	intersectionRe := regexp.MustCompile(`(\d)*\.?(\d)*`)
	for scanner.Scan() {
		if lineCount > 0 && lineCount <= 183749 { // intersections
			params := intersectionRe.FindAllString(scanner.Text(), -1)
			id, err := strconv.Atoi(params[0])
			lat, err := strconv.ParseFloat(params[1], 32)
			lng, err := strconv.ParseFloat(params[2], 32)
			checkError(err)

			roadnetwork.addIntersection(Intersection{id, float32(lat), float32(lng)})
		} else if lineCount > (183749 + 1) { // road segments
			if (lineCount-183749-1)%2 == 1 { // road segments attributes except for latlngs
				params := strings.Split(scanner.Text(), "	")
				rid++
				if rid >= ridUpper || ridLower > rid {
					// ridLower == rid 时，应该计算
					continue
				}

				l, err := strconv.ParseFloat(params[3], 32)

				// 在读文件时【不】过滤掉短的路

				if err != nil {
					log.Fatal(err)
				}
				var Vs = [8785]VelocityRecord{}
				lv, _ := strconv.Atoi(params[6])
				dir, _ := strconv.Atoi(params[5])
				level, _ := strconv.Atoi(params[7])
				roadnetwork.addRoadSegment(RoadSegment{rid, float32(l), Vs, lv, [][2]float32{}, dir, level})
			} else { // latlng for road segments
				params := strings.Split(scanner.Text(), ",")
				for _, latlngStr := range params {
					latAndlng := strings.Split(latlngStr, " ")
					// fmt.Println(latlngStr)
					latStr := latAndlng[0]
					lngStr := latAndlng[1]
					lRoadSegments := len(roadnetwork.RoadSegments)
					roadnetwork.RoadSegments[lRoadSegments-1].addShape(latStr, lngStr)
				}
			}
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(len(roadnetwork.roadSegments))

	return roadnetwork
}
