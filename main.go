package main

import (
	"codeberg.org/meadowingc/auto-wayback/site"
)

func main() {
	config, err := site.ReadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	go StartMonitorProcess(config)

	site.StartSiteProcess(config)
}
