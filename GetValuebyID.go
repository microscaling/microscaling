package main

import (
	//"fmt"
	//"time"
	//"sync"
	"log"
	//"strings"
	//"strconv"
	//"math/rand"
	"io/ioutil"
	"net/http"
	"os"
	//"bytes"
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
	var demandvalue string = ""
	var str string
	var port string
	port = os.Getenv("MARATHON_CONSUL_PORT")
	str = os.Getenv("MARATHON_ADDRESS")
	str = str + port
	if port == "" {
		port = "8500"
	}
	if str == "" {
		str = "http://marathon.force12.io:" + port
	}

	str += "/v1/kv/priority1-demand"

	log.Println("GET demand: " + str)
	resp, err := http.Get(str)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp == nil {
		// handle error
		log.Println("GET demand failed ")
	} else {
		body, err0 := ioutil.ReadAll(resp.Body)
		if err0 != nil {
			// handle error
			log.Println("GET demand failed read body ")
		} else {
			s := string(body[:])
			log.Println("demand json: " + s)
			demandvalue = s
		}
	}
	return demandvalue
}
