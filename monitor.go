package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"codeberg.org/meadowingc/auto-wayback/site"
	"codeberg.org/meadowingc/auto-wayback/waybacker"
)

func archivePage(pageUrl string, config *site.Config) {
	linksInSite, err := waybacker.GetSitemapURLs(pageUrl)
	if err != nil {
		log.Printf("Error getting sitemap URLs: %v\n", err)
		return
	}

	log.Printf("Found %v links in sitemap for %v\n", len(linksInSite), pageUrl)

	for _, link := range linksInSite {
		archiveErr := waybacker.RunIfChanged(link, func() error {
			var resp *http.Response
			var err error

			log.Printf("About to send to Wayback Machine: %v\n", link)

			for i := 0; i < 10; i++ { // Retry up in case of timeout
				resp, err = waybacker.SendToWaybackMachine(link, config.AccessKey, config.SecretKey)
				if err != nil {
					log.Printf("Error sending to Wayback Machine: %v\n", err)
					return err
				}

				if resp.StatusCode != 429 {
					break
				}

				log.Println("Received status code 429, sleeping and retrying...")
				time.Sleep(1 * time.Minute)
			}

			if resp.StatusCode == 200 {
				log.Printf("Sent to Wayback Machine: %v\n", link)
			} else {
				log.Printf("Error sending to Wayback Machine: %v\n", resp.Status)
				return fmt.Errorf("error sending to Wayback Machine: %v", resp.Status)
			}

			// be a nice person and don't spam the Wayback Machine
			time.Sleep(10 * time.Second)

			return nil
		})

		if archiveErr != nil {
			log.Printf("Error archiving link: %v\n", archiveErr)
		}
	}
}

func StartMonitorProcess(config *site.Config) {
	for {
		log.Println("Starting monitor process")

		for _, url := range config.URLs {
			log.Printf("Starting archiving for page: %v\n", url)
			archivePage(url, config)
			log.Printf("Finished archiving for page: %v\n", url)
		}

		time.Sleep(time.Duration(config.SleepDays) * 24 * time.Hour)
	}
}
