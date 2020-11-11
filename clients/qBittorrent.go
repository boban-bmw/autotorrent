package clients

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// QBittorrent implements the TorrentClient interface
type QBittorrent struct {
	headers http.Header
	client  http.Client
	prefix  string
}

// Init logs in to the qBittorrent web API
func (q *QBittorrent) Init(config ClientConfig) error {
	q.prefix = fmt.Sprintf("%s/%s", config.URL.String(), "api/v2")

	URL := fmt.Sprintf("%s/%s", q.prefix, "auth/login")

	payload := url.Values{}
	payload.Set("username", config.Username)
	payload.Set("password", config.Password)

	query := payload.Encode()

	req, err := http.NewRequest("POST", URL, strings.NewReader(query))
	if err != nil {
		return err
	}

	referer := fmt.Sprintf("%s://%s", config.URL.Scheme, config.URL.Host)
	req.Header.Add("Referer", referer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := q.client.Do(req)
	if err != nil {
		return err
	}

	q.headers = http.Header{}
	q.headers.Set("Referer", referer)

	cookie := res.Cookies()[0]
	q.headers.Set("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))

	return nil
}

// AddTorrent adds a torrent file to qbt
func (q *QBittorrent) AddTorrent() {}
