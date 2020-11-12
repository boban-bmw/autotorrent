package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func handleMultiFileTorrent(torrent *multiFileTorrent, downloads []node, links string, maxMissing int64) bool {
	fileMap := make(map[string]node)

	totalFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		totalFileSize = totalFileSize + torrentFile.Length

		for _, file := range downloads {
			if torrentFile.Length == file.info.Size() {
				fileMap[filepath.Join(torrentFile.Path...)] = file
				break
			}
		}

		if len(fileMap) == len(torrent.Info.Files) {
			break
		}
	}

	missingFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		if _, ok := fileMap[filepath.Join(torrentFile.Path...)]; !ok {
			missingFileSize = missingFileSize + torrentFile.Length
		}
	}

	if float64(missingFileSize)/float64(totalFileSize) > float64(maxMissing)/float64(100) {
		return false
	}

	err := os.Mkdir(filepath.Join(links, torrent.Info.Name), 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return false
	}

	for torrentFilePath, file := range fileMap {
		pathFragments := []string{links, torrent.Info.Name}
		pathFragments = append(pathFragments, torrentFilePath)

		completePath := filepath.Join(pathFragments...)
		fileDir := filepath.Dir(completePath)

		err := os.MkdirAll(fileDir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return false
		}

		err = os.Symlink(file.path, completePath)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println("Error linking", file.path, "->", completePath, err)
			return false
		}
	}

	return true
}
