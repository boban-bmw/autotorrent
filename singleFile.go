package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func handleSingleFileTorrent(torrent *singleFileTorrent, downloads []node, links string) bool {
	torrentSize := torrent.Info.Length

	var matchingFile node

	for _, file := range downloads {
		if file.info.Size() == torrentSize {
			matchFound, err := compareHashSingleFile(file, torrent)
			if err != nil {
				log.Println("An error ocurred comparing hashes", err)
				continue
			}

			if matchFound {
				matchingFile = file
			}
		}
	}

	err := os.Symlink(matchingFile.path, filepath.Join(links, torrent.Info.Name))
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println("Error linking", matchingFile.path, "->", torrent.Info.Name, err)
		return false
	}

	return true
}
