package readings

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Trajectories(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// file, err := os.Open("../../../data/traj_hz_20160304/20160304/20160415146.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer file.Close()
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
}
