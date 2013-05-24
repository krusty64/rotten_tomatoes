package rotten_tomatoes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetUrl(url string, response interface{}) error {
	log.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, response)
	if err != nil {
		return err
	}

	return nil
}
