package api

// GetEvent8785s ..
func GetEvent8785s(roadsEvents [][]int) [][8785]int {
	timeRanges := [][8785]int{}

	for _, roadEvents := range roadsEvents {
		timeRange := [8785]int{0}
		movingIndex := 0
		for index, t := range roadEvents {
			// 1,2,3,6,7,8  ,11,12,13,14,30
			// 0,1,2,3,4,5  ,6, 7, 8, 9, 10
			if index < movingIndex {
				continue
			}

			movingIndex = index
			for (roadEvents[movingIndex] - t) == (movingIndex - index) {
				movingIndex++
				if movingIndex >= len(roadEvents) {
					break
				}
			}

			timeRange[t] = roadEvents[movingIndex-1] + 1
			for i := t; i <= roadEvents[movingIndex-1]; i++ {
				timeRange[i] = 1
			}
		}

		timeRanges = append(timeRanges, timeRange)
	}

	return timeRanges
}

// GetEventsFull return boolean record of events
func GetEventsFull(eventsList [][]int, filename string) [][]int {
	eventsFullList := [][]int{}

	var length int
	if filename == "congestion" || filename == "congestion2" {
		length = 8785
	} else if filename == "flow" {
		length = 7221
	} else {
		length = 10000
	}

	for _, events := range eventsList {
		eventsFull := make([]int, length)

		movingIndex := 0
		for index, t := range events {
			if index < movingIndex {
				continue
			}

			movingIndex = index
			for (events[movingIndex] - t) == (movingIndex - index) {
				movingIndex++
				if movingIndex >= len(events) {
					break
				}
			}

			eventsFull[t] = events[movingIndex-1] + 1
			for i := t; i <= events[movingIndex-1]; i++ {
				eventsFull[i] = 1
			}
		}

		eventsFullList = append(eventsFullList, eventsFull)
	}

	return eventsFullList
}

// ReduceNoise return eventsList after reducing the noise
func ReduceNoise(eventsFullList [][]int, halfWindow int, threshold float64, filename string) [][]int {
	eventsListNoiseReduced := [][]int{}

	var length int
	if filename == "congestion" || filename == "congestion2" {
		length = 8785
	} else if filename == "flow" {
		length = 7221
	} else {
		length = 10000
	}

	for _, eventsFull := range eventsFullList { // times
		eventsNoiseReduced := []int{}

		for i, e := range eventsFull {
			nEvent := 0
			validTime := 0

			if e == 1 {
				eventsNoiseReduced = append(eventsNoiseReduced, i)
				continue
			}

			for j := i - halfWindow; j <= i+halfWindow; j++ {
				if j < 0 || j >= length {
					continue
				}

				if eventsFull[j] == 1 {
					nEvent++
				}
				validTime++
			}
			if ratio := float64(nEvent) / float64(validTime); ratio > threshold {
				eventsNoiseReduced = append(eventsNoiseReduced, i)
			}
		}
		eventsListNoiseReduced = append(eventsListNoiseReduced, eventsNoiseReduced)
	}

	return eventsListNoiseReduced
}
