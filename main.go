package main

import (
	"autotorrent/clients"
	"encoding/json"
	"errors"
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

var cache map[string][]byte = map[string][]byte{}
var cacheName string = "_autotorrent.cache"

func init() {
	rawCache, err := ioutil.ReadFile(cacheName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Println("Error while reading cache", err)
	}

	if len(rawCache) != 0 {
		err = json.Unmarshal(rawCache, &cache)
		if err != nil {
			log.Println("Error unmarshaling from cache", err)
		}
	}
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
		trackerDir := ""

		switch t := torrent.(type) {
		case *singleFileTorrent:
			trackerDir, err = createTrackerDir(linksDir, t.Announce)
			if err != nil {
				log.Fatal("Couldn't create tracker folder", t.Announce, err)
			}

			filesFound = handleSingleFileTorrent(t, downloads, trackerDir)
			path = t.path
		case *multiFileTorrent:
			trackerDir, err = createTrackerDir(linksDir, t.Announce)
			if err != nil {
				log.Fatal("Couldn't create tracker folder", t.Announce, err)
			}

			filesFound = handleMultiFileTorrent(t, downloads, trackerDir, opts.MaxMissing)
			path = t.path
		}

		if filesFound {
			err = client.AddTorrent(path, trackerDir, opts.ClientCategory)
			if err != nil {
				log.Println("Error adding torrent", torrent, err)
			}
		}
	}

	rawCache, err := json.Marshal(cache)
	if err != nil {
		log.Println("Error marshaling cache", err)
	}

	ioutil.WriteFile(cacheName, rawCache, 0755)
}

func createTrackerDir(linksDir string, trackerURL string) (string, error) {
	tracker, err := url.Parse(trackerURL)
	if err != nil {
		log.Println("Couldn't parse tracker name", trackerURL)
		return "", err
	}

	trackerDir := filepath.Join(linksDir, tracker.Hostname())

	err = os.Mkdir(trackerDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}

	return trackerDir, nil
}
