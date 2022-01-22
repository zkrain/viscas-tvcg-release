package readings

import (
	"os"
	"path/filepath"
)

func FileNames(root string) []string {
	var files []string

	// root := "../data/traj_hz_20160304/20160304"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return files
	// for _, file := range files {
	// 	fmt.Println(file)
	// }
}
