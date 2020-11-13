package main

import (
	"crypto/sha1"
	"os"
)

func compareHashSingleFile(potentialMatch node, torrent *singleFileTorrent) (bool, error) {
	file, err := os.Open(potentialMatch.path)
	if err != nil {
		return false, err
	}

	defer file.Close()

	fileSlice := make([]byte, torrent.Info.PieceLength)

	bytesRead, err := file.Read(fileSlice)
	if err != nil {
		return false, err
	}

	if bytesRead != torrent.Info.PieceLength {
		fileSlice = fileSlice[:bytesRead]
	}

	fileHash := sha1.Sum(fileSlice)

	var torrentHashArr [20]byte

	copy(torrentHashArr[:], torrent.Info.pieces[:20])

	return fileHash == torrentHashArr, nil
}
