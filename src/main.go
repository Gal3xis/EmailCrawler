package main

import (
	"EmailCrawler/src/conf"
)

var Config = conf.Config{}

func main() {
	// Config einlesen
	err := Config.Read()
	if err != nil {
		return
	}
}
