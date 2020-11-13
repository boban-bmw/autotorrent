package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/jackpal/bencode-go"
)

type singleFileTorrentInfo struct {
	Length      int64
	Name        string
	PieceLength int    `bencode:"piece length"`
	PiecesRaw   string `bencode:"pieces"`
	pieces      []byte
}

type singleFileTorrent struct {
	Announce string
	Info     singleFileTorrentInfo
	path     string
}

type torrentFile struct {
	Length int64
	Path   []string
}

type multiFileTorrentInfo struct {
	Files       []torrentFile
	Name        string
	PieceLength int    `bencode:"piece length"`
	PiecesRaw   string `bencode:"pieces"`
	pieces      []byte
}

type multiFileTorrent struct {
	Announce string
	Info     multiFileTorrentInfo
	path     string
}

func parseTorrents(fileNames []string) []interface{} {
	torrents := make([]interface{}, 0)

	for _, fileName := range fileNames {
		rawTorrent, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Println("Couldn't read", fileName, err)
			continue
		}

		t := bytes.NewReader(rawTorrent)

		var torrent interface{}

		torrent = &multiFileTorrent{}

		err = bencode.Unmarshal(t, torrent)
		if err != nil {
			log.Println("Couldn't unmarshal multi-file torrent", fileName, err)
		}

		if torrent.(*multiFileTorrent).Info.Files == nil {
			t = bytes.NewReader(rawTorrent)

			torrent = &singleFileTorrent{}

			err = bencode.Unmarshal(t, torrent)
			if err != nil {
				log.Println("Couldn't unmarshal single-file torrent", fileName, err)
				continue
			}

			torrent.(*singleFileTorrent).path = fileName
			torrent.(*singleFileTorrent).Info.pieces = []byte(torrent.(*singleFileTorrent).Info.PiecesRaw)
		} else {
			torrent.(*multiFileTorrent).path = fileName
			torrent.(*multiFileTorrent).Info.pieces = []byte(torrent.(*multiFileTorrent).Info.PiecesRaw)
		}

		torrents = append(torrents, torrent)
	}

	return torrents
}
