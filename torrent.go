package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/jackpal/bencode-go"
)

type singleFileTorrentInfo struct {
	Length      int
	Name        string
	PieceLength int    `bencode:"piece length"`
	PiecesRaw   string `bencode:"pieces"`
	pieces      []byte
}

type singleFileTorrent struct {
	Announce string
	Info     singleFileTorrentInfo
}

type torrentFile struct {
	Length int
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
		}

		switch tt := torrent.(type) {
		case *singleFileTorrent:
			tt.Info.pieces = []byte(tt.Info.PiecesRaw)
		case *multiFileTorrent:
			tt.Info.pieces = []byte(tt.Info.PiecesRaw)
		default:
			log.Printf("Unknown type %T encountered\n", tt)
			continue
		}

		torrents = append(torrents, torrent)
	}

	return torrents
}
