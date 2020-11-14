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

func parseDownloads(path string) map[int64][]node {
	downloads := map[int64][]node{}

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error while walking across", path, err)
			return nil
		}

		if !info.IsDir() {
			downloads[info.Size()] = append(downloads[info.Size()], node{
				info: info,
				path: path,
			})
		}

		return nil
	})

	return downloads
}
