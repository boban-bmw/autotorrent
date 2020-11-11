package clients

import "net/url"

// ClientConfig contains configuration parameters for a torrent client
type ClientConfig struct {
	Username string
	Password string
	URL      url.URL
	Label    string
}

// TorrentClient exposes all methods we need from a torrent client
type TorrentClient interface {
	Init(config ClientConfig)
}
