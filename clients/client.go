package clients

import (
	"fmt"
	"net/url"
)

// ClientConfig contains configuration parameters for a torrent client
type ClientConfig struct {
	Username string
	Password string
	URL      url.URL
	Category string
}

// TorrentClient exposes all methods we need from a torrent client
type TorrentClient interface {
	Init(config ClientConfig) error
	AddTorrent(path string, linksDir string, category string) error
}

// GetClient initializes a TorrentClient
func GetClient(config ClientConfig, id string) (TorrentClient, error) {
	var client TorrentClient

	switch id {
	case "qbt":
		client = &QBittorrent{}
	}

	if client != nil {
		err := client.Init(config)
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	return nil, fmt.Errorf("Unknown client type %v", id)
}
