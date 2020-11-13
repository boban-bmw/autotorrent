package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type match struct {
	torrentPath string
	file        node
}

func handleMultiFileTorrent(torrent *multiFileTorrent, downloads []node, links string, maxMissing int64) bool {
	potentialMatches := []match{}

	totalFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		totalFileSize += torrentFile.Length

		for _, file := range downloads {
			if torrentFile.Length == file.info.Size() {
				potentialMatches = append(potentialMatches, match{
					torrentPath: filepath.Join(torrentFile.Path...),
					file:        file,
				})
			}
		}
	}

	matches, err := compareHashMultiFile(potentialMatches, torrent)
	if err != nil {
		log.Println("Error comparing multi file hash", torrent, err)
		return false
	}

	missingFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		found := false

		for _, match := range matches {
			if match.torrentPath == filepath.Join(torrentFile.Path...) {
				found = true
				break
			}
		}

		if !found {
			missingFileSize += torrentFile.Length
		}
	}

	if float64(missingFileSize)/float64(totalFileSize) > float64(maxMissing)/float64(100) {
		return false
	}

	err = os.Mkdir(filepath.Join(links, torrent.Info.Name), 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return false
	}

	for _, match := range matches {
		pathFragments := []string{links, torrent.Info.Name}
		pathFragments = append(pathFragments, match.torrentPath)

		completePath := filepath.Join(pathFragments...)
		fileDir := filepath.Dir(completePath)

		err := os.MkdirAll(fileDir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return false
		}

		err = os.Symlink(match.file.path, completePath)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println("Error linking", match.file.path, "->", completePath, err)
			return false
		}
	}

	return true
}
