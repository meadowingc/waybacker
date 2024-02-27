package waybacker

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

const recordFile = "records.json"

type HashRecord struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

func fetchAndHash(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Remove CSRF tokens
	re := regexp.MustCompile(`<input[^>]*name="csrfmiddlewaretoken"[^>]*>`)
	modifiedBody := re.ReplaceAll(body, nil)

	// Hash the modified HTML
	hasher := sha256.New()
	if _, err := io.Copy(hasher, bytes.NewReader(modifiedBody)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func getStoredHash(url string) (string, error) {
	file, err := os.OpenFile(recordFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	records := make(map[string]string)
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	return records[url], nil
}

func updateHashRecord(url, newHash string) error {
	file, err := os.OpenFile(recordFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	records := make(map[string]string)
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		if err != io.EOF {
			return err
		}
	}

	records[url] = newHash

	file, err = os.Create(recordFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(records)
}

func RunIfChanged(link string, callback func() error) error {
	// Parse URL
	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	// Normalize URL
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path = path.Clean(u.Path)

	// Remove www
	u.Host = strings.TrimPrefix(u.Host, "www.")

	// Remove default port
	if (u.Scheme == "http" && u.Port() == "80") || (u.Scheme == "https" && u.Port() == "443") {
		u.Host = u.Hostname()
	}

	// Remove trailing slash
	link = strings.TrimRight(u.String(), "/")

	newHash, err := fetchAndHash(link)
	if err != nil {
		return err
	}

	storedHash, err := getStoredHash(link)
	if err != nil {
		return err
	}

	if newHash != storedHash {
		if err := callback(); err != nil {
			return err
		}

		if err := updateHashRecord(link, newHash); err != nil {
			return err
		}
	}

	return nil
}
