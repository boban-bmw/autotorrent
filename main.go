package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Torrents  string `short:"t" long:"torrents" description:"Path to directory with .torrent files, relative to current directory" default:"."`
	Downloads string `short:"d" long:"downloads" description:"Path to downloads directory, relative to current directory" default:"."`
	Links     string `short:"l" long:"links" description:"Path to links directory, relative to current directory" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln("Couldn't parse flags", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't get current working directory", err)
	}

	torrentsDir := filepath.Join(cwd, opts.Torrents)
	downloadsDir := filepath.Join(cwd, opts.Downloads)
	linksDir := filepath.Join(cwd, opts.Links)

	files, err := ioutil.ReadDir(torrentsDir)
	if err != nil {
		log.Fatalln("Couldn't read", torrentsDir, err)
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

	for _, torrent := range torrents {
		// TODO: check if torrent is already added in client

		filesFound := false

		switch t := torrent.(type) {
		case *singleFileTorrent:
			filesFound = handleSingleFileTorrent(t, downloads, linksDir)
		case *multiFileTorrent:
			filesFound = handleMultiFileTorrent(t, downloads, linksDir)
		}

		if filesFound {
			// TODO: add torrent to client
		}
	}
}
