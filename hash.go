package main

import (
	"crypto/sha1"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
)

func getFileHash(filePath string, pieceLength int64, offset int64) ([sha1.Size]byte, error) {
	fileHash := [sha1.Size]byte{}

	file, err := os.Open(filePath)
	if err != nil {
		return fileHash, err
	}

	defer file.Close()

	fileSlice := make([]byte, pieceLength)

	bytesRead, err := file.ReadAt(fileSlice, offset)
	if err != nil && !errors.Is(err, io.EOF) {
		return fileHash, err
	}

	if int64(bytesRead) != pieceLength {
		fileSlice = fileSlice[:bytesRead]
	}

	fileHash = sha1.Sum(fileSlice)

	return fileHash, nil
}

func compareHashSingleFile(potentialMatchPath string, pieceLength int64, pieceHash [sha1.Size]byte, offset int64) (bool, error) {
	fileHash, err := getFileHash(potentialMatchPath, pieceLength, offset)
	if err != nil {
		return false, nil
	}

	return fileHash == pieceHash, nil
}

func compareHashMultiFile(potentialMatches []match, torrent *multiFileTorrent) ([]match, error) {
	pieceLength := torrent.Info.PieceLength

	matches := []match{}

	for _, potentialMatch := range potentialMatches {
		fileOffset := int64(0)

		for _, torrentFile := range torrent.Info.Files {
			if potentialMatch.torrentPath == filepath.Join(torrentFile.Path...) {
				break
			}

			fileOffset += torrentFile.Length
		}

		// bordering piece contains both prev file(s) and current file - this is how much of it is from current file
		fileBeginPieceMod := fileOffset % pieceLength

		// easy case - file contains at least 1 whole piece
		if fileBeginPieceMod+pieceLength < potentialMatch.file.info.Size() {
			pieceIndex := int64(math.Floor(float64(fileOffset)/float64(pieceLength)) + 1)

			pieceStart := pieceIndex * sha1.Size
			pieceEnd := (pieceIndex + 1) * sha1.Size

			pieceHash := [sha1.Size]byte{}

			copy(pieceHash[:], torrent.Info.pieces[pieceStart:pieceEnd])

			fileReadOffset := int64(pieceIndex)*pieceLength - fileBeginPieceMod

			fileMatches, err := compareHashSingleFile(potentialMatch.file.path, pieceLength, pieceHash, fileReadOffset)
			if err != nil {
				continue
			}

			if fileMatches {
				matches = append(matches, potentialMatch)
			}
		}
	}

	return matches, nil
}
