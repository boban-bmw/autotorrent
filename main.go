package main

import (
	"autotorrent/clients"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Torrents       string `short:"t" long:"torrents" description:"Path to directory with .torrent files, relative to current directory" default:"."`
	Downloads      string `short:"d" long:"downloads" description:"Path to downloads directory, relative to current directory" default:"."`
	Links          string `short:"l" long:"links" description:"Path to links directory, relative to current directory" required:"true"`
	ClientUsername string `short:"u" long:"username" description:"Torrent client username" required:"true"`
	ClientPassword string `short:"p" long:"password" description:"Torrent client password" required:"true"`
	ClientURL      string `long:"url" description:"Torrent client URL" required:"true"`
	ClientCategory string `short:"c" long:"category" description:"Category for the added torrent" required:"true"`
	ClientID       string `short:"i" long:"client" description:"Id of the torrent client" required:"true" choice:"qbt"`
	MaxMissing     int64  `short:"m" long:"max-missing-percent" description:"Maximum missing percentage in a torrent" default:"5"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln("Couldn't parse flags", err)
	}

	clientURL, err := url.Parse(opts.ClientURL)
	if err != nil {
		log.Fatalln("Couldn't parse client URL", err)
	}

	client, err := clients.GetClient(clients.ClientConfig{
		Username: opts.ClientUsername,
		Password: opts.ClientPassword,
		URL:      *clientURL,
		Category: opts.ClientCategory,
	}, opts.ClientID)
	if err != nil {
		log.Fatalln("Couldn't connect to client", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't get current working directory", err)
	}

	torrentsDir := filepath.Join(cwd, opts.Torrents)
	downloadsDir := filepath.Join(cwd, opts.Downloads)
	linksDir := filepath.Join(cwd, opts.Links)

	files, err := ioutil.ReadDir(torrentsDir)
	if err != nil {
		log.Fatalln("Couldn't read", torrentsDir, err)
	}

	fileNames := make([]string, 0)

	for _, file := range files {
		fileName := file.Name()

		if strings.HasSuffix(fileName, ".torrent") {
			fileNames = append(fileNames, filepath.Join(torrentsDir, fileName))
		}
	}

	torrents := parseTorrents(fileNames)
	downloads := parseDownloads(downloadsDir)

	for _, torrent := range torrents {
		filesFound := false
		path := ""

		switch t := torrent.(type) {
		case *singleFileTorrent:
			filesFound = handleSingleFileTorrent(t, downloads, linksDir)
			path = t.path
		case *multiFileTorrent:
			filesFound = handleMultiFileTorrent(t, downloads, linksDir, opts.MaxMissing)
			path = t.path
		}

		if filesFound {
			err = client.AddTorrent(path, linksDir, opts.ClientCategory)
			if err != nil {
				log.Println("Error adding torrent", torrent, err)
			}
		}
	}
}
