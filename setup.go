package rotten_tomatoes

import (
	"log"
)

const (
	init_query = "http://api.rottentomatoes.com/api/public/v1.0.json"
	api_key    = "?apikey=" + rtkey
)

type RtConfig struct {
	Links struct {
		Movies string
	}
	LinkTemplate string `json:"link_template"`
}

var (
	url_setup RtConfig
)

func init() {
	url_setup.Links.Movies = "not initialized"
}

func InitSetup() {
	// Fills out Links.Movies and Links.Lists
	err := GetUrl(init_query+api_key, &url_setup)
	if err != nil {
		log.Println(err)
	}
	//log.Println(url_setup)

	if url_setup.Links.Movies != "" {
		// Fills out LinkTemplate
		err = GetUrl(url_setup.Links.Movies+api_key, &url_setup)
		if err != nil {
			log.Println(err)
		}
	}
}
