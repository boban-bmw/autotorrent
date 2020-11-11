package main

func handleSingleFileTorrent(torrent *singleFileTorrent, downloads []node, links string) {
	torrentSize := torrent.Info.Length

	for _, file := range downloads {
		if file.info.Size() == torrentSize {
			// create link
			// send to dl client
		}
	}
}
