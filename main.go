package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirFlag := flag.String("dir", ".", "Path to directory with .torrent files, relative to current folder")

	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory", err)
	}

	dir := filepath.Join(cwd, *dirFlag)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Couldn't read", dir, err)
	}

	fileNames := make([]string, 0)

	for _, file := range files {
		fileName := file.Name()

		if strings.HasSuffix(fileName, ".torrent") {
			fileNames = append(fileNames, filepath.Join(dir, fileName))
		}
	}

	torrents := parseTorrents(fileNames)

	log.Println(torrents)
}
