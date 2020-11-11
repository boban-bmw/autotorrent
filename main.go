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
	torrentsDirFlag := flag.String("torrentsDir", ".", "Path to directory with .torrent files, relative to current directory")
	downloadsDirFlag := flag.String("downloadsDir", ".", "Path to downloads directory, relative to the current directory")
	linksDirFlag := flag.String("linksDir", "", "Path to the links directory, relative to the current directory")

	flag.Parse()

	if *linksDirFlag == "" {
		log.Fatal("linksDir must be set!")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory", err)
	}

	torrentsDir := filepath.Join(cwd, *torrentsDirFlag)
	downloadsDir := filepath.Join(cwd, *downloadsDirFlag)
	linksDir := filepath.Join(cwd, *linksDirFlag)

	files, err := ioutil.ReadDir(torrentsDir)
	if err != nil {
		log.Fatal("Couldn't read", torrentsDir, err)
	}

	fileNames := make([]string, 0)

	for _, file := range files {
		fileName := file.Name()

		if strings.HasSuffix(fileName, ".torrent") {
			fileNames = append(fileNames, filepath.Join(torrentsDir, fileName))
		}
	}

	torrents := parseTorrents(fileNames)
	downloads := parseDownloads(downloadsDir)

	log.Println(torrents, downloads)

	for _, torrent := range torrents {
		switch t := torrent.(type) {
		case *singleFileTorrent:
			handleSingleFileTorrent(t, downloads, linksDir)
		case *multiFileTorrent:

		}
	}
}
