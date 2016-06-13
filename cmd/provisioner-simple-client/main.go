package main

import (
	"fmt"
	"log"
	"os"
)

var testJson = `{
	"userId": "12345",
	"ApplicationId": "45678",
	"ApiKey": "aNdjvg398djeh389eEEHEjejkhsaxhsZ"
}
`

func init() {
	// show date/time in log output.
	log.SetFlags(log.LstdFlags)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [socket path]\n", os.Args[0])
	os.Exit(1)
}

func tryGet(client *SocketClient, path string) {
	if str, err := client.Get(path); err == nil {
		log.Printf("GET  /%s: %s", path, str)
	} else {
		log.Printf("GET  /%s: ERROR: %s\n", path, err)
	}
}

func tryPost(client *SocketClient, path, json string) {
	if str, err := client.PostJsonString(path, json); err == nil {
		log.Printf("POST /%s: %s", path, str)
	} else {
		log.Printf("POST /%s: ERROR: %s\n", path, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	path := os.Args[1]
	client := NewSocketClient(path)

	tryGet(client, "provisioned")
	tryGet(client, "provision")
	tryGet(client, "config")

	tryPost(client, "provision", testJson)
	//tryPost(client, "provisioned", testJson)
	//tryPost(client, "config", testJson)
}
