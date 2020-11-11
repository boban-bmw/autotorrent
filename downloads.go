package main

import (
	"log"
	"os"
	"path/filepath"
)

type node struct {
	info os.FileInfo
	path string
}

func parseDownloads(path string) []node {
	downloads := make([]node, 0)

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error while walking across", path, err)
			return nil
		}

		downloads = append(downloads, node{
			info: info,
			path: path,
		})

		return nil
	})

	return downloads
}
