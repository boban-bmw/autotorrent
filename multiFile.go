package main

import (
	"errors"
	"log"
	"math"
	"os"
	"path/filepath"
)

func handleMultiFileTorrent(torrent *multiFileTorrent, downloads []node, links string, maxMissingFiles int) bool {
	matchingFiles := make([]node, 0)

	for _, file := range downloads {
		for _, torrentFile := range torrent.Info.Files {
			if torrentFile.Length == file.info.Size() {
				matchingFiles = append(matchingFiles, file)
			}
		}
	}

	if math.Abs(float64(len(matchingFiles)-len(torrent.Info.Files))) > float64(maxMissingFiles) {
		return false
	}

	err := os.Mkdir(filepath.Join(links, torrent.Info.Name), 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return false
	}

	for _, matchingFile := range matchingFiles {
		err := os.Symlink(matchingFile.path, filepath.Join(links, torrent.Info.Name, matchingFile.info.Name()))
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println("Error linking", matchingFile.path, "->", torrent.Info.Name, err)
			return false
		}
	}

	return true
}
