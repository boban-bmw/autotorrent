package clients

// QBittorrent implements the TorrentClient interface
type QBittorrent struct {
	referrer string
	cookie   string
}

// Init logs in to the qBittorrent web API
func (q *QBittorrent) Init(config ClientConfig) {

}
