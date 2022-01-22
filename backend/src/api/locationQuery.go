package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Location struct {
	Rid int     `json:"rid"`
	Lng float32 `json:"lng"`
	Lat float32 `json:"lat"`
}

type LocationQueryer struct {
	AllLocation  []Location
	LocationDict map[int]Location
}

func InitializeLocationQueryer(filename string) LocationQueryer {
	var m map[int]Location
	m = make(map[int]Location)

	jsonFile, err := os.Open("../output/" + filename + "Data/locations.json") // todo 要改回来
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var locations []Location
	json.Unmarshal(byteValue, &locations)

	for _, location := range locations {
		m[location.Rid] = location
	}

	myLocationQueryer := LocationQueryer{locations, m}
	return myLocationQueryer
}

// func LocationQueryer(filename string) map[int]Location {
// 	var m map[int]Location
// 	m = make(map[int]Location)

// 	jsonFile, err := os.Open("../output/" + filename + "Data/locations.json") // todo 要改回来
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	var locations []Location
// 	json.Unmarshal(byteValue, &locations)

// 	for _, location := range locations {
// 		m[location.Rid] = location
// 	}
// 	return m
// }
