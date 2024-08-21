package waybacker

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
)

type Urlset struct {
	Urls []Url `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

func fetchAndParseSitemap(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var urlset Urlset
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(bytes, &urlset)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, u := range urlset.Urls {
		urls = append(urls, u.Loc)
	}

	return urls, nil
}

func GetSitemapURLs(url string) ([]string, error) {
	log.Printf("Getting sitemap URLs for %v\n", url)

	url = strings.TrimSuffix(url, "/")
	var sitemapUrls []string

	// Check if sitemap.xml exists
	urls, err := fetchAndParseSitemap(url + "/sitemap.xml")
	if err == nil {
		sitemapUrls = append(sitemapUrls, urls...)
	} else {
		// If not, check robots.txt for Sitemap entries
		resp, err := http.Get(url + "/robots.txt")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Sitemap:") {
				sitemap := strings.TrimSpace(strings.TrimPrefix(line, "Sitemap:"))
				urls, err := fetchAndParseSitemap(sitemap)
				if err != nil {
					return nil, err
				}
				sitemapUrls = append(sitemapUrls, urls...)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return sitemapUrls, nil
}
