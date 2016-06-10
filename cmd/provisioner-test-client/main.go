package main

import (
	"fmt"
	"log"
	"os"
)

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
		log.Printf("GET  /%s: '%s'\n", path, str)
	} else {
		log.Printf("GET  /%s: ERROR: %s\n", path, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	path := os.Args[1]
	client := NewSocketClient(path)

	tryGet(client, "provision")

	json := `{"foo": "abc", "bar": 3}`
	if str, err := client.PostJsonString("provision", json); err == nil {
		log.Printf("POST /provision: '%s'\n", str)
	} else {
		log.Printf("POST /provision: ERROR: %s\n", err)
	}
}
