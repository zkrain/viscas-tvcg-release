package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	hamiltonian "congest/Hamiltonian"
	"congest/api"
	"congest/graph"
)

type GetEmbeddingsResponse struct {
	Code       int
	Embeddings []api.ProjPoint
}

type GetDistributionResponse struct {
	Code         int
	Distribution []int
	MaxCount     int
	MinCount     int
}

type IDResponse struct {
	Code int
	IDs  []int
}

type ChangeBackendDatasetParam struct {
	Filename string `json:"filename"`
}

type FilterIDsByCountParam struct {
	VLow  int `json:"vLow"`
	VHigh int `json:"vHigh"`
}

type ReduceNoiseForEventsParam struct {
	Window    int     `json:"window"`
	Threshold float64 `json:"threshold"`
}

type CodeResponse struct {
	Code int
}

type LocationIDsParam struct {
	LocationIDs []int `json:"locationIDs"`
}

type GetLocationsResponse struct {
	Code      int
	Locations []api.Location
}

type GetLocationsInfoResponse struct {
	Code                int
	Locations           []api.Location
	LocationsEventCount []int
	LocationsStyle      []LocationStyle
}

type GetLocationInfectEvents struct {
	Code             int                       `json:"code"`
	InfectEventsDict map[int][]api.InfectEvent `json:"infectEventsDict"`
}

type EventsResponse struct {
	Code   int
	Events [][]int
}

type InferCascadingPatternParam struct {
	LocationIDs []int `json:"locationIDs"`
	TimeWindow  int   `json:"tw"`
	K           int   `json:"k"`
}

type InferCascadingPatternResponse struct {
	Code                         int                            `json:"code"`
	InfluencesDict               map[int]map[int]api.Influences `json:"influencesDict"`
	InfluencesDictContraryInG    map[string]api.Influences      `json:"influencesDictContraryInG"`
	InfluencesDictContraryNotInG map[string]api.Influences      `json:"influencesDictContraryNotInG"`
	Links                        [][2]int                       `json:"links"`
	Trees                        api.PropagationTrees           `json:"trees"`
}

type LocationAverage struct {
	Code                 int                      `json:"code"`
	AverageEventDataDict api.AverageEventDataDict `json:"averageEventDict"`
}

type StClusteringParam struct {
	LocationIDs []int   `json:"locationIDs"`
	EpsT        float64 `json:"EpsT"`
	EpsS        float64 `json:"EpsS"`
	MinPts      int     `json:"MinPts"`
	DeltaE      float64 `json:"deltaE"`
}

type StClusteringResponse struct {
	Code                          int
	Clusters                      [][]int
	ClusterEventCount             [][]int
	AverageEventCount             []float32
	AverageInfectEventCount       []float32
	Noise                         []int
	ClusterInformationList        []api.ClusterInformation
	SpatialDistributionOfClusters api.SpatialDistributionOfClusters
}

type LocationStyle struct {
	AngRate           float32 `json:"angRate"`
	NormalizedAngRate float32 `json:"normalizedAngRate"`
	AveDuration       float32 `json:"aveDuration"`
}

type EventNumRange struct {
	MinNum int
	MaxNum int
}

type HamiltonianWalkParam struct {
	Vertices [][]int
}

type HamiltonianWalkWithHeightParam struct {
	Vertices [][]int `json:"vertices"`
	Heights  []int   `json:"heights"`
}

type HamiltonianWalkResponse struct {
	Order []int `json:"order"`
}

// global variables, including embeddings, locations, and events, as well as the road shape (special)
var embeddingQuery api.EmbeddingQuery
var IDsFilter api.IDsFilter
var locationQueryer api.LocationQueryer
var g graph.Graph

var eventsListBackup [][]int
var eventsListNoiseReducedBackup [][]int
var IDsBackup []int
var locationsBackup []api.Location
var globalFilename string
var globalTimeSpan map[string]int
var globalEventNumRange map[string]EventNumRange

func initGlobalVariable() {
	globalTimeSpan = make(map[string]int)
	globalTimeSpan["air"] = 8760
	globalTimeSpan["congestion2"] = 8785
	globalTimeSpan["flow"] = 7221

	globalEventNumRange = make(map[string]EventNumRange)
	globalEventNumRange["air"] = EventNumRange{9, 5129}
	globalEventNumRange["congestion2"] = EventNumRange{12, 3952}
	globalEventNumRange["flow"] = EventNumRange{2, 5588}
}

func getLocationInfectEvents(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams LocationIDsParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	locationIDs := reqParams.LocationIDs
	response := GetLocationInfectEvents{200, IDsFilter.GetInfectEventsIDict(locationIDs)}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func getLocationAverage(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	averageDict := IDsFilter.AverageEventDataDict
	response := LocationAverage{200, averageDict}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getEmbeddings return the projection of the selected dataset
func getEmbeddings(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	var response GetEmbeddingsResponse
	response = GetEmbeddingsResponse{200, embeddingQuery.AllProjPoints}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getDistribution calculate the distribution of count
func getDistribution(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	distribution, max, min := IDsFilter.GetDistribution()
	response := GetDistributionResponse{200, distribution, max, min}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getAllID return all ids
func getAllID(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	IDs := IDsFilter.GetAllID()
	response := IDResponse{200, IDs}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// filterIDsByCount
func filterIDsByCount(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams FilterIDsByCountParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	fmt.Println(reqParams.VLow, reqParams.VHigh)

	IDs := IDsFilter.FilterIDsByCount(reqParams.VLow, reqParams.VHigh)

	response := IDResponse{200, IDs}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// changeBackendDataset change dataset
func changeBackendDataset(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams ChangeBackendDatasetParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	fmt.Println("filename", reqParams.Filename)
	globalFilename = reqParams.Filename
	filename := reqParams.Filename
	if filename == "" {
		return
	}

	embeddingQuery = api.InitializeEmbeddingQuery(filename)
	locationQueryer = api.InitializeLocationQueryer(filename)
	IDsFilter = api.InitialIDsFilter(filename)

	response := CodeResponse{200}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getLocations return all locations
func getAllLocations(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	var response GetLocationsResponse
	locations := locationQueryer.AllLocation
	response = GetLocationsResponse{200, locations}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getLocations return locations given a set of ids
func getLocations(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams LocationIDsParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	locationIDs := reqParams.LocationIDs
	var response GetLocationsInfoResponse

	// get shapes
	var locationCounts []int
	var locations []api.Location
	var locationsStyle []LocationStyle
	for _, locationID := range locationIDs {
		locations = append(locations, locationQueryer.LocationDict[locationID])
		infectEvents := IDsFilter.InfectEventDataDict[locationID]
		events := IDsFilter.EventDataDict[locationID].Events
		eventsLen := float32(len(events))
		smallAng := eventsLen / float32(globalTimeSpan[globalFilename])
		normalizedSmallAng := (eventsLen - float32(globalEventNumRange[globalFilename].MinNum)) /
			(float32(globalEventNumRange[globalFilename].MaxNum) - float32(globalEventNumRange[globalFilename].MinNum))
		infectDurationSum := 0
		for _, event := range infectEvents {
			infectDurationSum += event.DurationTime
		}
		var infectEventsLen = float32(len(infectEvents))
		// fmt.Println("debug", infectEvents)
		var infectDurationAve = float32(infectDurationSum) / infectEventsLen
		locationCounts = append(locationCounts, IDsFilter.EventDataDict[locationID].Count)
		locationsStyle = append(locationsStyle, LocationStyle{smallAng, normalizedSmallAng, infectDurationAve})
	}

	locationsBackup = locations

	// get events?
	response = GetLocationsInfoResponse{200, locations, locationCounts, locationsStyle}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// getEventsOfCorrelatedLocations return events by locationID
func getEventsOfCorrelatedLocations(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams LocationIDsParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	locationIDs := reqParams.LocationIDs

	IDsBackup = locationIDs

	var response EventsResponse
	eventsList := [][]int{}
	for _, locationID := range locationIDs {
		eventsList = append(eventsList, IDsFilter.EventDataDict[locationID].Events)
	}

	// store the eventsList in backend
	eventsListBackup = eventsList
	eventsListNoiseReducedBackup = eventsList

	response = EventsResponse{200, eventsList}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// reduceNoiseWithFilter
func reduceNoiseForEvents(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams ReduceNoiseForEventsParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)

	window := reqParams.Window
	threshold := reqParams.Threshold
	halfWindow := (window - 1) / 2

	eventsFullList := api.GetEventsFull(eventsListBackup, globalFilename)
	eventsListNoiseReduced := api.ReduceNoise(eventsFullList, halfWindow, threshold, globalFilename)
	eventsListNoiseReducedBackup = eventsListNoiseReduced

	response := EventsResponse{200, eventsListNoiseReduced}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// inferCascadingPattern corresponds to infer-new.go
func inferCascadingPattern(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams InferCascadingPatternParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)
	locationIDs := reqParams.LocationIDs
	k := reqParams.K
	tw := reqParams.TimeWindow

	if len(locationIDs) == 0 {
		return
	}

	timeSpanLength := globalTimeSpan[globalFilename]
	infectEventsList := [][]api.InfectEvent{}
	basicEpss := []float64{}
	for _, locationID := range locationIDs {
		infectEvents := IDsFilter.InfectEventDataDict[locationID]
		EventsWithCountID := IDsFilter.EventDataDict[locationID]
		infectEventsList = append(infectEventsList, infectEvents)
		basicEps := float64(EventsWithCountID.Count) / float64(timeSpanLength)
		basicEpss = append(basicEpss, basicEps)
	}
	var influencesInG map[int]map[int]api.Influences
	var influencesDictContraryInG map[string]api.Influences
	var influencesDictContraryNotInG map[string]api.Influences
	var links [][2]int
	var trees api.PropagationTrees
	links, influencesInG, trees, influencesDictContraryInG, influencesDictContraryNotInG = api.NetworkInfer(locationIDs, k, tw, locationsBackup, infectEventsList, basicEpss, timeSpanLength)

	response := InferCascadingPatternResponse{200, influencesInG, influencesDictContraryInG, influencesDictContraryNotInG, links, trees}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func stclustering(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams StClusteringParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)

	locationIDs := reqParams.LocationIDs
	locations := []api.Location{}
	for _, locationID := range locationIDs {
		locations = append(locations, locationQueryer.LocationDict[locationID])
	}
	epsT := reqParams.EpsT
	epsS := reqParams.EpsS
	minPts := reqParams.MinPts
	deltaE := reqParams.DeltaE

	if minPts == 0 {
		return
	}

	fmt.Println(epsT, epsS, minPts, deltaE)

	clusters, noise := api.Stdbscan(locations, epsT, epsS, minPts, deltaE, embeddingQuery.ProjPointDict)
	fmt.Println(len(clusters), " clusters has been generated")

	clusterInformationList := api.ContextualizeClustersCompactly(clusters, IDsFilter.EventDataDict, globalFilename)
	spatialDistributionOfClusters := api.GetSpatialDistributionOfClusters(clusters, locationQueryer.LocationDict)
	clusterAverageCounts := api.GetClusterAverageEventCount(clusters, IDsFilter.EventDataDict)
	clusterCounts := api.GetClusterEventCount(clusters, IDsFilter.EventDataDict)
	clusterInfectEvents := api.GetClusterInfectEvent(clusters, IDsFilter.InfectEventDataDict, 1)
	var clusterInfectEventNumberAve []float32
	for _, infectEvent := range clusterInfectEvents {
		var infectEventSum int = 0
		for _, event := range infectEvent {
			infectEventSum += len(event)
		}
		clusterInfectEventNumberAve = append(clusterInfectEventNumberAve, float32(infectEventSum)/float32(len(infectEvent)))
	}
	response := StClusteringResponse{200, clusters, clusterCounts,
		clusterAverageCounts, clusterInfectEventNumberAve, noise,
		clusterInformationList, spatialDistributionOfClusters}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func hamiltonianWalkOptimizer(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams HamiltonianWalkWithHeightParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)

	vertices := reqParams.Vertices
	heights := reqParams.Heights
	// fmt.Println(vertices)
	// vertices := [][]int{{2,3,9},{4,5,9},{3,4,5},{1,2,3,4},{1,3,9},{6,7,8,9},{4,5,8,9},{1,3,8,9},{1}, {4},{3,4,6},{3,11},{6,7,11,2},{2,11,10,9},{11},{9},{8},{8,9}}
	// vertices := [][]int{{2,3,9},{4,5,9},{3,4,5},{1,2,3,4},{1,3,9},{6,7,8,9},{4,5,8,9},{1,3,8,9},{1}, {4},{3,4,6}}
	// vertices := [][]int{{2,3,9},{4,5,9},{3,4,5},{1,2,3,4},{1,3,9},{6,7,8,9},{4,5,8,9}}
	// vertices := [][]int{{2,3,9},{4,5,9},{3,4,5},{1,2,3,4}}
	N := len(vertices)

	if N == 0 {
		return
	}

	matrix := [][]float64{}
	for range vertices {
		row := make([]float64, N)
		matrix = append(matrix, row)
	}
	for i, vertexi := range vertices {
		for j, vertexj := range vertices {
			if i >= j {
				continue
			}
			d := hamiltonian.JaccardDistance(vertexi, vertexj)
			hi := heights[i]
			hj := heights[j]
			hWeight := 0.0
			if hi > hj {
				hWeight = float64(hi) / float64(hj)
			} else {
				hWeight = float64(hj) / float64(hi)
			}
			matrix[i][j] = d * hWeight
			matrix[j][i] = d * hWeight
		}
	}

	var optimalPath []int
	if N < 17 {
		optimalPath = hamiltonian.FastHamiltonianWalk(matrix, uint64(N))
	} else {
		optimalPath = hamiltonian.GA(matrix, uint64(N))
	}
	// optimalPath := hamiltonian.FastHamiltonianWalk(matrix, uint64(N))
	// optimalPath := hamiltonian.GA(matrix, uint64(N))

	response := HamiltonianWalkResponse{optimalPath}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func hamiltonianWalkOptimizerEdge(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams HamiltonianWalkParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)

	vertices := reqParams.Vertices
	// fmt.Println(vertices)
	N := len(vertices)

	if N == 0 {
		return
	}

	matrix := [][]float64{}
	for range vertices {
		row := make([]float64, N)
		matrix = append(matrix, row)
	}
	for i, vertexi := range vertices {
		for j, vertexj := range vertices {
			if i >= j {
				continue
			}
			d := hamiltonian.WeightedJaccardDistance(vertexi, vertexj)
			matrix[i][j] = d
			matrix[j][i] = d
		}
	}

	var optimalPath []int
	if N < 17 {
		optimalPath = hamiltonian.FastHamiltonianWalk(matrix, uint64(N))
	} else {
		optimalPath = hamiltonian.GA(matrix, uint64(N))
	}

	response := HamiltonianWalkResponse{optimalPath}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func hamiltonianWalkOptimizerInfluence(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	defer req.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(req.Body)
	var reqParams HamiltonianWalkParam
	json.Unmarshal([]byte(string(jsonBlob)), &reqParams)

	vertices := reqParams.Vertices
	// fmt.Println(vertices)
	N := len(vertices)

	if N == 0 {
		return
	}

	matrix := [][]float64{}
	for range vertices {
		row := make([]float64, N)
		matrix = append(matrix, row)
	}
	for i, vertexi := range vertices {
		for j, vertexj := range vertices {
			if i >= j {
				continue
			}
			d := hamiltonian.EuclideanDistance(vertexi, vertexj)
			matrix[i][j] = d
			matrix[j][i] = d
		}
	}

	var optimalPath []int
	if N < 18 {
		optimalPath = hamiltonian.FastHamiltonianWalk(matrix, uint64(N))
	} else {
		optimalPath = hamiltonian.MyGA(matrix, uint64(N))
	}

	response := HamiltonianWalkResponse{optimalPath}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

// enableCors allows cross-origin
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// embeddings = api.LoadEmbedding(filename)
	// locationQueryer = api.LocationQueryer(filename)
	// IDsFilter = api.InitialIDsFilter(filename)
	initGlobalVariable()

	http.HandleFunc("/getEmbeddings", getEmbeddings)
	http.HandleFunc("/getDistribution", getDistribution)
	http.HandleFunc("/getAllID", getAllID)
	http.HandleFunc("/filterIDsByCount", filterIDsByCount)
	http.HandleFunc("/changeBackendDataset", changeBackendDataset)
	http.HandleFunc("/getLocations", getLocations)
	http.HandleFunc("/getLocationAverage", getLocationAverage)
	http.HandleFunc("/getAllLocations", getAllLocations)
	http.HandleFunc("/getLocationInfectEvents", getLocationInfectEvents)
	http.HandleFunc("/getEventsOfCorrelatedLocations", getEventsOfCorrelatedLocations)
	http.HandleFunc("/reduceNoiseForEvents", reduceNoiseForEvents)
	// http.HandleFunc("/inferCascadingPatternsPatterns", inferCascadingPatternsPatterns)
	http.HandleFunc("/inferCascadingPattern", inferCascadingPattern)
	http.HandleFunc("/hamiltonianWalkOptimizer", hamiltonianWalkOptimizer)
	http.HandleFunc("/hamiltonianWalkOptimizerEdge", hamiltonianWalkOptimizerEdge)
	http.HandleFunc("/hamiltonianWalkOptimizerInfluence", hamiltonianWalkOptimizerInfluence)
	http.HandleFunc("/stclustering", stclustering)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
