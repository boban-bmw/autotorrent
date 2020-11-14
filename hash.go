package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func getKey(filePath string, pieceLength int64, offset int64) string {
	return fmt.Sprintf("%v-%v-%v", filePath, pieceLength, offset)
}

func getFileHash(filePath string, pieceLength int64, offset int64) ([sha1.Size]byte, error) {
	cacheKey := getKey(filePath, offset, offset)

	cachedByteArray, found := cache[cacheKey]
	if found {
		var cachedHash [sha1.Size]byte

		copy(cachedHash[:], cachedByteArray)

		return cachedHash, nil
	}

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

	cache[cacheKey] = fileHash[:]

	return fileHash, nil
}

func compareHashSingleFile(potentialMatchPath string, pieceLength int64, pieceHash [sha1.Size]byte, offset int64) (bool, error) {
	fileHash, err := getFileHash(potentialMatchPath, pieceLength, offset)
	if err != nil {
		return false, nil
	}

	return fileHash == pieceHash, nil
}

func compareHashMultiFile(potentialMatchesMap map[string][]match, torrent *multiFileTorrent) ([]match, error) {
	pieceLength := torrent.Info.PieceLength

	matches := []match{}

	for torrentFilePath, potentialMatches := range potentialMatchesMap {
		fileOffset := int64(0)
		fileSize := int64(0)

		for _, torrentFile := range torrent.Info.Files {
			if torrentFilePath == filepath.Join(torrentFile.Path...) {
				fileSize = torrentFile.Length

				break
			}

			fileOffset += torrentFile.Length
		}

		// bordering piece contains both prev file(s) and current file - this is how much of it is from current file
		fileBeginPieceMod := fileOffset % pieceLength

		// file contains at least 1 whole piece - use it to check if it's ok
		if fileBeginPieceMod+pieceLength < fileSize {
			pieceIndex := int64(math.Floor(float64(fileOffset)/float64(pieceLength)) + 1)

			pieceStart := pieceIndex * sha1.Size
			pieceEnd := (pieceIndex + 1) * sha1.Size

			pieceHash := [sha1.Size]byte{}

			copy(pieceHash[:], torrent.Info.pieces[pieceStart:pieceEnd])

			fileReadOffset := int64(pieceIndex)*pieceLength - fileBeginPieceMod

			for _, potentialMatch := range potentialMatches {
				fileMatches, err := compareHashSingleFile(potentialMatch.file.path, pieceLength, pieceHash, fileReadOffset)
				if err != nil {
					continue
				}

				if fileMatches {
					matches = append(matches, potentialMatch)
					break
				}
			}
		} else {
			if len(potentialMatches) >= 1 {
				matches = append(matches, potentialMatches[0])
			}
		}
	}

	return matches, nil
}
