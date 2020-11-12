package clients

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
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
	if err != nil || len(res.Cookies()) == 0 {
		return fmt.Errorf("Network error or wrong credentials %v", err)
	}

	err = res.Body.Close()
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
func (q *QBittorrent) AddTorrent(path string, linksDir string, category string) error {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)

	fileName := filepath.Base(path)

	part, err := writer.CreateFormFile("torrents", fileName)
	if err != nil {
		return err
	}

	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	part.Write(fileContents)

	writer.WriteField("savepath", linksDir)
	writer.WriteField("category", category)
	writer.WriteField("skip_checking", "true")
	writer.WriteField("paused", "true")

	err = writer.Close()
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/%s", q.prefix, "torrents/add")

	req, err := http.NewRequest("POST", URL, &body)
	if err != nil {
		return err
	}

	req.Header = q.headers.Clone()
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := q.client.Do(req)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
