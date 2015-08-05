package main

import (
	//"fmt"
	//"time"
	//"sync"
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	//"math/rand"
	//"net/http"
	//"os"
	//"io/ioutil"
	//"bytes"
)

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
func Decode_ContainerCount(encoded string) int {
	// decode json returned from Consul KV store
	var demand int = 0
	var json_prefix string = "Value\":\""
	stringslice := strings.Split(encoded, json_prefix)
	log.Println("Encoding: " + encoded)
	log.Println("Search for: " + json_prefix)

	if len(stringslice) >= 2 {
		splitstr := strings.Split(stringslice[1], "\"")
		log.Println("split: " + splitstr[0])
		var str64 string
		str64, err1 := base64Decode(splitstr[0])
		if err1 {
			// handle error
			log.Println("base64 decoding error")
		}
		container_count, err2 := strconv.Atoi(str64)
		if err2 != nil {
			// handle error
			log.Println("base64 decoding error")
		}
		demand = container_count
		log.Println("container count: " + strconv.Itoa(demand))
	} else {
		log.Println("Length only " + strconv.Itoa(len(stringslice)))
	}
	return demand
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
}
