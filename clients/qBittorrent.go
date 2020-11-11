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
}

// Init logs in to the qBittorrent web API
func (q *QBittorrent) Init(config ClientConfig) error {
	loginURL := fmt.Sprintf("%s/%s", config.URL.String(), "api/v2/auth/login")

	loginPayload := url.Values{}
	loginPayload.Set("username", config.Username)
	loginPayload.Set("password", config.Password)

	encodedPayload := loginPayload.Encode()

	loginReq, err := http.NewRequest("POST", loginURL, strings.NewReader(encodedPayload))
	if err != nil {
		return err
	}

	referer := fmt.Sprintf("%s://%s", config.URL.Scheme, config.URL.Host)
	loginReq.Header.Add("Referer", referer)
	loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(loginReq)
	if err != nil {
		return err
	}

	q.headers = http.Header{}
	q.headers.Set("Referer", referer)

	cookie := resp.Cookies()[0]
	q.headers.Set("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))

	return nil
}
