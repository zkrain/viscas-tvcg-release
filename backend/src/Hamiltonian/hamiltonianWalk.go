package hamiltonian

import (
  "fmt"
	"math"
	"math/rand"
	"reflect"

	set "github.com/fatih/set"
)

type Triple struct {
	A float64
	B int
	C int
}

func FastHamiltonianWalk (matrix [][]float64, N uint64) []int {
	dp := [][]Triple{}
	for range matrix {
		row := make([]Triple, (1<<N)+1)
		dp = append(dp, row)
	}
	for i := 0; i < int(N); i++ {
		dp[i][1<<uint64(i)].A = 0.01
	}
	for i := 0; i < (1<<N); i++ {
		for j := 0; j < int(N); j++ {
			if (i & (1<<uint64(j))) > 0.0 {
				for k := 0; k < int(N); k++ {
					if ((i & (1<<uint64(k))) > 0) && (k != j) && (dp[k][i^(1<<uint64(j))].A > 0.0) {
						v := matrix[k][j] + dp[k][i^(1<<uint64(j))].A
						if dp[j][i].A <= 0.000001 {
							dp[j][i].A = v
							dp[j][i].B = k
							dp[j][i].C = i^(1<<uint64(j))
						} else if (dp[j][i].A > 0.0) && (dp[j][i].A > v) {
							dp[j][i].A = v
							dp[j][i].B = k
							dp[j][i].C = i^(1<<uint64(j))
						}
					}
				}
			}
		}
	}

	minCost := 10000000.0
	mini := -1
	var minDp Triple
	for i:= 0; i < int(N); i++ {
		cost := dp[i][(1<<N)-1].A
		if cost < minCost {
			minCost = cost
			minDp = dp[i][(1<<N)-1]
			mini = i
		}
	}

	// fmt.Println("mind ", minCost)

	path := []int{}
	localDp := minDp
	path = append(path, mini)
	for len(path) < int(N) {
		path = append(path, localDp.B)
		localDp = dp[localDp.B][localDp.C]
	}
	return path
}

func EuclideanDistance(vertexi []int, vertexj []int) float64 {
	if len(vertexi) != len(vertexj) {
		fmt.Println("error, len(vertexi) != len(vertexj).")
	}
	sum := 0.0
	for k := 0; k < len(vertexi); k++ {
		vi := vertexi[k]
		vj := vertexj[k]
		v := (vi - vj) * (vi - vj)
		sum += float64(v)
	}

	return math.Sqrt(sum)
}

func WeightedJaccardDistance(vertexi []int, vertexj []int) float64 {
	up := 0
	bottom := 0
	for k := 0; k < len(vertexi); k++ {
		v1 := vertexi[k]
		v2 := vertexj[k]
		if v1 > v2 {
			up += v2
			bottom += v1
		} else {
			up += v1
			bottom += v2
		}
	}
	return 1.0 - float64(up)/float64(bottom)
}

func JaccardDistance(vertexi []int, vertexj []int) float64 {
	vertexSetI := set.New(set.ThreadSafe)
	vertexSetJ := set.New(set.ThreadSafe)
	for _, v := range vertexi {
		vertexSetI.Add(v)
	}
	for _, v := range vertexj {
		vertexSetJ.Add(v)
	}
	intersection := set.Intersection(vertexSetI, vertexSetJ)
	union := set.Union(vertexSetI, vertexSetJ)
	d := 1.0 - float64(intersection.Size()) / float64(union.Size())
	return d * 10.0
}

func HamiltonianWalk(matrix [][]float64, N int) []int {

	optimalPath := []int{}
	min := 10000000.0
	mini := -1
	for i := 0; i < N; i++ {
		// fmt.Println(i, N)
		v := dp(int64(math.Pow(2, float64(N)) - 1), int64(i), matrix, int64(N))
		if v < min {
			min = v
			mini = i
		}
	}
	optimalPath = append(optimalPath, mini)
	retrace(int64(math.Pow(2, float64(N)) - 1), int64(mini), matrix, int64(N), &optimalPath)
	// fmt.Println(optimalPath, min)
	return optimalPath
}

func retrace(mask int64, i int64, matrix [][]float64, N int64, optimalPath *[]int) {
	v := dp(mask, i, matrix, N)
	for j := 0; j < int(N); j++ {
		newMask := mask ^ int64(math.Pow(2, float64(i)))
		newV := dp(newMask, int64(j), matrix, N)
		if (matrix[j][i] + newV) == v {
			// fmt.Println(j)
			*optimalPath = append(*optimalPath, j)
			retrace(newMask, int64(j), matrix, N, optimalPath)
		}
	}
}

func countBit(number int64) int {
	nBit := 0
	for i := 0; i < 64; i ++ {
		if number < (1 << uint64(i)) {
			return nBit
		}
		v := number & (1 << uint64(i))
		if v > 0 {
			nBit++
		}
	}
	return nBit
}

func dp(mask int64, i int64, matrix [][]float64, N int64) float64 {
	nBit := countBit(mask)
	bit := mask & (1 << uint64(i))
	if (nBit == 1) && (bit > 0) {
		return 0
	} else if (nBit > 1) && (bit > 0) {
		min := 1000000.0
		// minNode := -1
		for j := 0; j < int(N); j++ {
			flag := mask & (1 << uint64(j))
			if flag > 0 {
				newMask := mask ^ int64(math.Pow(2, float64(i)))
				v := dp(newMask, int64(j), matrix, N) + matrix[j][i]
				if v < min {
					min = v
				}
			}
		}
		return min
	} else {
		return 33333
	}
}

func GA(matrix [][]float64, N uint64) []int {
	rand.Seed(1)
	k := 0
	nIter := 0
	maxNumberOfPeople := 0
	if N > 50 {
		nIter = 800
		maxNumberOfPeople = 130
		// fmt.Println("N: ", N)
	} else {
		nIter = 400
		maxNumberOfPeople = 100
	}

	// initialize
	firstGeneration := [][]int{}
	numberOfFirstGeneration := 0
	for numberOfFirstGeneration < maxNumberOfPeople {
		firstGeneration = append(firstGeneration, arrRandom(int(N)))
		numberOfFirstGeneration++
	}

	// fmt.Println(firstGeneration)

	currentGeneration := firstGeneration
	for k < nIter {
		judgements, sum, _ := judge(currentGeneration, matrix)

		children := [][]int{}
		for len(children) < maxNumberOfPeople {
			parent1RandomNumber := rand.Float64() * sum
			parent2RandomNumber := rand.Float64() * sum

			parent1Index := getItemIndexByJudge(parent1RandomNumber, judgements)
			parent2Index := getItemIndexByJudge(parent2RandomNumber, judgements)

			parent1 := currentGeneration[parent1Index]
			parent2 := currentGeneration[parent2Index]

			var child []int
			if rand.Float64() > 0.3 {
				i := rand.Intn(int(N) - 1)
				l := rand.Intn(int(N) - i)
				if l == 0 {
					child = parent1
				} else {
					child = cross(&parent1, &parent2, i, l)
					// fmt.Println("child.length", len(child))
				}
			} else {
				child = parent1
			}

			if rand.Float64() > 0.9 {
				i := rand.Intn(int(N) - 1)
				j := rand.Intn(int(N) - 1)
				if i != j {
					swap(&child, i, j)
				}
			}

			children = append(children, child)
		}
		currentGeneration = children

		k++
	}

	bestList := []int{}
	mind := 100000.0
	for _, list := range currentGeneration {
		d := getCostGA(list, matrix)
		if d < mind {
			bestList = list
			mind = d
		}
	}
	// fmt.Println("final mind ", mind)


	return bestList
}

func MyGA(matrix [][]float64, N uint64) []int {
	// generation := arrRandom(int(N))
	generation := []int{}
	for i := int(N) - 1; i >= 0; i-- {
		generation = append(generation, i)
	}
	// for i := 0; i < int(N); i++ {
	// 	generation = append(generation, i)
	// }

	k := 0

	for k < 400 {
		localImprovement := -1.0
		localI := -1
		localJ := -1
		for i := range matrix {
			for j := range matrix {
				if j > i {
					costBefore := getCostGA(generation, matrix)
					swap(&generation, i, j)
					costAfter := getCostGA(generation, matrix)
					improvement := costBefore - costAfter

					if localImprovement < improvement {
						localImprovement = improvement
						localI = i
						localJ = j
					} else {
						swap(&generation, i, j)
					}
				}
			}
		}
		// cost := getCostGA(generation, matrix)
		// fmt.Println(N, k, localImprovement, localI, localJ)
		// fmt.Println(N, k, localImprovement, localI, localJ, cost)
		k++
		if localImprovement < 0.3 && k > 100 {
			return generation
		}
	}

	return generation
}

func getItemIndexByJudge(number float64, judgements []float64) int {
	sum := 0.0
	for index, judge := range judgements {
		sum += judge
		if sum > number {
			return index
		}
	}
	return len(judgements) - 1
}

func getCostGA(list []int, matrix[][]float64) float64 {
	// d := 0.0
	// for i :=0; i < len(list)-2; i++ {
	// 	vertex1 := list[i]
	// 	vertex2 := list[i+1]
	// 	vertex3 := list[i+2]
	// 	ds := []float64{}
	// 	mindlocal := 0.0
	// 	ds = append(ds, matrix[vertex1][vertex2])
	// 	ds = append(ds, matrix[vertex1][vertex3])
	// 	ds = append(ds, matrix[vertex2][vertex3])
	// 	for _, dOfDs := range ds {
	// 		if mindlocal < dOfDs {
	// 			mindlocal = dOfDs
	// 		}
	// 	}

	// 	d += mindlocal

	// }

	// d := 0.0
	// for i :=0; i < len(list)-1; i++ {
	// 	vertex1 := list[i]
	// 	vertex2 := list[i+1]
	// 	d += matrix[vertex1][vertex2]
	// }

	d := 0.0
	for i :=0; i < len(list)-2; i++ {
		vertex1 := list[i]
		vertex2 := list[i+1]
		vertex3 := list[i+2]

		d += matrix[vertex1][vertex2]
		d += 0.5 * matrix[vertex1][vertex3]

	}
	return d
}

func judge(generation [][]int, matrix [][]float64) ([]float64, float64, float64) {
	judgements := []float64{}
	sum := 0.0
	mind := 100000.0
	costs := []float64{}
	for _, list := range generation {
		d := getCostGA(list, matrix)
		costs = append(costs, d)

		if d < mind {
			mind = d
		}
	}

	for _, c := range costs {
		judgements = append(judgements, 1.5/(c*1.02-mind))
		sum += 1.5/(c*1.02-mind)
	}

	return judgements, sum, mind
}

func arrRandom(n int) []int {
	arr := []int{}
	arr = append(arr, rand.Intn(n))

	for (len(arr) < n) {
		num := rand.Intn(n)
		sign := true // 判断生成的数字是否与数组中已有数字重复的标志位；
		for _, v := range arr {
			if (num == v) {
				sign = false
				break
			}
		}
		if(sign) {
			arr = append(arr, num)
		}
	}
	return arr
}

func swap(list *[]int, i int, j int) {
	if i >= len(*list) || j >= len(*list) {
		fmt.Println("error", *list, len(*list), i , j)
		return
	}

	tmp := (*list)[i]
	(*list)[i] = (*list)[j]
	(*list)[j] = tmp
}

func cross(list1 *[]int, list2 *[]int, i int, length int) []int {
	if i > len((*list1)) || i > len((*list2)) || (i+length) > len((*list1)) || (i+length) > len((*list2)) {
		fmt.Println("error 2")
	}
	// i + length <= len(), e.g., 2 + 3 <= 5
	sub := (*list2)[i:i+length]

	newList1 := []int{}
	flag := 0
	for _, item := range (*list1) {
		if flag == i {
			for _, itemInSub := range sub {
				newList1 = append(newList1, itemInSub)
			}
			flag = 10000
		}
		if !(itemExists(sub, item)) {
			newList1 = append(newList1, item)
			flag++
		}
	}
	// fmt.Println("cross ", sub, i, *list1, *list2, i, length, newList1)

	return newList1
}

func itemExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
