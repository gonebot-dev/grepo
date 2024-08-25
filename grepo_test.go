package grepo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	currentDir, _ := os.Getwd()
	println(currentDir)
	filepath.Walk(currentDir, func(path string, info os.FileInfo, errr error) error {
		if filepath.Dir(path) != currentDir || info.IsDir() {
			return nil
		}
		println(path)
		return nil
	})
}
