package marathon

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
)

// consulKey describes the JSON we get back from consul when looking up a key
type consulKey struct {
	Value string
}

////////////////////////////////////////////////////////////////////////////////////////
// Decode_ContainerCount
//
// Called to decode a json snippet to return the container count
//
// input encoded json
// returned container count
//
// Json format as below. Note: Value is base64 encoded for Unicode support
//[
//    {
//        "CreateIndex": 8,
//        "ModifyIndex": 15,
//        "LockIndex": 0,
//        "Key": "priority1-demand",
//        "Flags": 0,
//        "Value": "OQ=="
//    }
//]
//
//
func DecodeContainerCount(encoded string) int {
	// decode json returned from Consul KV store
	key := []consulKey{}
	err := json.Unmarshal([]byte(encoded), &key)
	if err != nil {
		log.Printf("Failed to decode container count. %v", err)
		// TODO: Error strategy!
		return 0
	}

	if len(key) != 1 {
		log.Printf("Failed to decode container count. Expected single entry array, have %d", len(key))
		return 0
	}

	data, err := base64.StdEncoding.DecodeString(key[0].Value)
	if err != nil {
		log.Printf("Failed to base64decode container count. %v", err)
		return 0
	}

	containerCount, err := strconv.Atoi(string(data))
	if err != nil {
		log.Printf("Failed to convert container count. %v", err)
		return 0
	}
	log.Printf("Container count: %d", containerCount)
	return containerCount
}
