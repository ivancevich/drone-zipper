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
	Files  []string `json:"files"`
	Name   string   `json:"name"`
	Output string   `json:"output"`
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

	if len(vargs.Name) == 0 {
		vargs.Name = "archive"
	}

	if len(vargs.Output) == 0 {
		vargs.Output = "."
	}

	err := zipThem(vargs.Files, workspace.Path, vargs.Name, vargs.Output)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

func zipThem(files []string, basePath, target, output string) error {
	zipFile, err := os.Create(filepath.Join(basePath, output, target+".zip"))
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
			return errors.New("Directories are not supported")
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate
		header.Name = info.Name()

		fileContents, err := os.Open(current)
		if err != nil {
			return err
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, fileContents)
		if err != nil {
			return err
		}

		err = fileContents.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
