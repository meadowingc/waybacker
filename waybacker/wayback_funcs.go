package waybacker

import (
	"net/http"
	"net/url"
	"strings"
)

func SendToWaybackMachine(targetURL, accessKey, accessSecret string) (*http.Response, error) {
	waybackURL := "https://web.archive.org/save"
	data := url.Values{}
	data.Set("url", targetURL)
	data.Set("capture_all", "1")
	data.Set("capture_outlinks", "1")
	data.Set("skip_first_archive", "1")
	data.Set("delay_wb_availability", "1")
	data.Set("if_not_archived_within", "20d")

	req, err := http.NewRequest("POST", waybackURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "LOW "+accessKey+":"+accessSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
