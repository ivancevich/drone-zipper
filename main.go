package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Files []string `json:files`
}

var (
	buildDate string
)

func main() {
	fmt.Printf("Drone Zipper Plugin built at %s\n", buildDate)

	workspace := plugin.Workspace{}
	vargs := Config{}

	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if len(vargs.Files) == 0 {
		return
	}

	err := zipThem(vargs.Files, workspace.Path, "Foo") // TODO: support custom names

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

// Zip files or directories
func zipThem(files []string, basePath, target string) error {
	zipFile, err := os.Create(filepath.Join(basePath, target+".zip"))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, file := range files {
		current := filepath.Join(basePath, file)

		info, err := os.Stat(current)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return errors.New("File is a directory")
		}

		data, err := os.Open(current)
		if err != nil {
			return err
		}

		f, err := archive.Create(target + "/" + info.Name())
		if err != nil {
			return err
		}

		io.Copy(f, data)
		data.Close()
	}

	return nil
}
