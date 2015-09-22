package consul

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// consulKey describes the JSON we get back from consul when looking up a key
//
// The format we get back looks like the following. Currently we only need the Value
// [
//    {
//        "CreateIndex": 8,
//        "ModifyIndex": 15,
//        "LockIndex": 0,
//        "Key": "priority1-demand",
//        "Flags": 0,
//        "Value": "OQ=="
//    }
// ]
type consulKey struct {
	Value string
}

// GetValuebyID gets the contents of an item in the Consul KV store, as identified by the item's unique ID Key
//
// input unique ID (Key) of target item 
// output string representation of the stored value
func (d *DemandFromConsul) GetValuebyID(key string) (string, error) {
	// Code to get value from Consul
	// GET http://marathon.force12.io:8500/v1/kv/priority1-demand

	url := d.baseConsulUrl + "/v1/kv/" + key // TODO! Move /v1/kv/ somewhere more sensible

	log.Println("GET demand: " + url)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return "", fmt.Errorf("GET value by ID failed %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GET value by ID failed %s", resp.Status)
	}

	// The key values are returned as an array
	keyData := []consulKey{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&keyData)
	if err != nil {
		return "", err
	}
	if len(keyData) != 1 {
		return "", fmt.Errorf("Expected 1 key, have %d", len(keyData))
	}
	return keyData[0].Value, nil
}
