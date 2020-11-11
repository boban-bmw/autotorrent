package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// GetTorrentList gets a list of existing torrents in qbt
func (q *QBittorrent) GetTorrentList() ([]Torrent, error) {
	URL := fmt.Sprintf("%s/%s", q.prefix, "torrents/info")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = q.headers.Clone()

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	torrents := []Torrent{}

	err = json.Unmarshal([]byte(body), &torrents)
	if err != nil {
		return nil, err
	}

	return torrents, nil
}
