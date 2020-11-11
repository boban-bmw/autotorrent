package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/agnivade/levenshtein"
)

func handleSingleFileTorrent(torrent *singleFileTorrent, downloads []node, links string) bool {
	torrentSize := torrent.Info.Length

	matchingFiles := make([]node, 0)

	for _, file := range downloads {
		if file.info.Size() == torrentSize {
			matchingFiles = append(matchingFiles, file)
		}
	}

	var matchingFile node

	if len(matchingFiles) == 0 {
		return false
	} else if len(matchingFiles) == 1 {
		matchingFile = matchingFiles[0]
	} else {
		matchingFile = getBestMatchingFile(matchingFiles, torrent.Info.Name)
	}

	err := os.Symlink(matchingFile.path, filepath.Join(links, torrent.Info.Name))
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println("Error linking", matchingFile.path, "->", torrent.Info.Name, err)
		return false
	}

	return true
}

func getBestMatchingFile(matchingFiles []node, torrentFileName string) node {
	index := 0
	minDistance := 500

	for i, file := range matchingFiles {
		distance := levenshtein.ComputeDistance(file.info.Name(), torrentFileName)

		if distance < minDistance {
			index = i
			minDistance = distance
		}
	}

	return matchingFiles[index]
}
