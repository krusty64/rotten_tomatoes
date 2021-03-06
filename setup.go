package rotten_tomatoes

import (
	"log"
	"net/url"
	"sync"
	"utils/json"
)

const (
	init_query = "http://api.rottentomatoes.com/api/public/v1.0.json"
	api_key    = "?apikey=" + rtkey // rtkey is in key.go
)

type Config struct {
	Links struct {
		Movies string
	}
	LinkTemplate string `json:"link_template"`
	LinkUrl      *url.URL
}

func (c *Config) AddKey(inUrl string) (string, error) {
	u, err := url.Parse(inUrl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("apikey", rtkey)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

var globalConfigLock sync.Mutex
var globalConfig *Config

func GetConfig() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	globalConfigLock.Lock()
	defer globalConfigLock.Unlock()

	// Fills out Links.Movies and Links.Lists
	var config Config
	err := json.FromUrl(init_query+api_key, &config)
	if err != nil {
		log.Println(err)
		return nil
	}

	if config.Links.Movies != "" {
		// Fills out LinkTemplate
		err = json.FromUrl(config.Links.Movies+api_key, &config)
		if err != nil {
			log.Println(err)
		} else {
			u, err := url.Parse(config.LinkTemplate)
			if err != nil {
				log.Println(err)
			} else {
				q := u.Query()
				q.Set("page_limit", "10")
				q.Set("page", "1")
				q.Set("apikey", rtkey)
				u.RawQuery = q.Encode()

				config.LinkUrl = u

				globalConfig = &config
			}
		}
	}

	return globalConfig
}
