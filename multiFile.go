package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/agnivade/levenshtein"
)

type match struct {
	torrentPath string
	file        node
	levDistance int
}

type byLevDistance []match

func (m byLevDistance) Len() int {
	return len(m)
}
func (m byLevDistance) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m byLevDistance) Less(i, j int) bool {
	return m[i].levDistance < m[j].levDistance
}

func handleMultiFileTorrent(torrent *multiFileTorrent, downloads map[int64][]node, links string, maxMissing int64) bool {
	potentialMatches := map[string][]match{}

	totalFileSize := int64(0)

	for _, torrentFile := range torrent.Info.Files {
		totalFileSize += torrentFile.Length

		sizeMatches, ok := downloads[torrentFile.Length]
		if !ok {
			continue
		}

		fullPath := filepath.Join(torrentFile.Path...)

		for _, sizeMatch := range sizeMatches {
			potentialMatches[fullPath] = append(potentialMatches[fullPath], match{
				torrentPath: fullPath,
				file:        sizeMatch,
				levDistance: levenshtein.ComputeDistance(sizeMatch.info.Name(), filepath.Base(fullPath)),
			})
		}
	}

	for _, matchesFragment := range potentialMatches {
		sort.Sort(byLevDistance(matchesFragment))

		if len(potentialMatches) > 11 {
			matchesFragment = matchesFragment[:11]
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
