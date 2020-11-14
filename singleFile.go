package main

import (
	"crypto/sha1"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func handleSingleFileTorrent(torrent *singleFileTorrent, downloads map[int64][]node, links string) bool {
	potentialMatches, ok := downloads[torrent.Info.Length]
	if !ok {
		return false
	}

	var firstPiece [sha1.Size]byte

	copy(firstPiece[:], torrent.Info.pieces[:20])

	var matchingFile *node

	for _, potentialMatch := range potentialMatches {
		matchFound, err := compareHashSingleFile(potentialMatch.path, torrent.Info.PieceLength, firstPiece, 0)
		if err != nil {
			log.Println("An error ocurred comparing hashes", err)
			continue
		}

		if matchFound {
			matchingFile = &potentialMatch
			break
		}
	}

	if matchingFile == nil {
		return false
	}

	err := os.Symlink(matchingFile.path, filepath.Join(links, torrent.Info.Name))
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println("Error linking", matchingFile.path, "->", torrent.Info.Name, err)
		return false
	}

	return true
}
