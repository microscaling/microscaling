package marathon

import (
	"io/ioutil"
	"log"
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////////////
// GetValuebyID
//
// Called to get the contents of an item in the Consul KV store, as identified by the item's unique ID Key
//
// input unique ID (Key) of target item (actually not used, it's hardcoded for now)
// output contents in json format or "" if error
//
func GetValuebyID(key string) string {
	// Code to get value from Consul
	// GET http://marathon.force12.io:8500/v1/kv/priority1-demand
	url := getBaseConsulUrl() + "/v1/kv/priority1-demand"

	log.Println("GET demand: " + url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		// handle error
		log.Printf("GET demand failed %v", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Printf("GET demand failed read body %v", err)
		return ""
	}
	s := string(body)
	log.Println("demand json: " + s)
	return s
}
