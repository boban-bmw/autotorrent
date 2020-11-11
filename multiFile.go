package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func handleMultiFileTorrent(torrent *multiFileTorrent, downloads []node, links string, maxMissing int64) bool {
	matchingFiles := make([]node, 0)

	for _, file := range downloads {
		for _, torrentFile := range torrent.Info.Files {
			if torrentFile.Length == file.info.Size() {
				matchingFiles = append(matchingFiles, file)
				continue
			}
		}

		if len(matchingFiles) == len(torrent.Info.Files) {
			break
		}
	}

	missingFileSize := int64(0)
	totalFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		totalFileSize = totalFileSize + torrentFile.Length

		fileFound := false

		for _, matchingFile := range matchingFiles {
			if torrentFile.Length == matchingFile.info.Size() {
				fileFound = true
				break
			}
		}

		if !fileFound {
			missingFileSize = missingFileSize + torrentFile.Length
		}
	}

	if missingFileSize/totalFileSize > maxMissing/100 {
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
