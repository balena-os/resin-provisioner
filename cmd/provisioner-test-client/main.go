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

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	path := os.Args[1]
	client := NewSocketClient(path)

	if str, err := client.Get("provision"); err == nil {
		log.Printf("GET  /provision: '%s'\n", str)
	} else {
		log.Printf("GET  /provision: ERROR: %s\n", err)
	}

	json := `{"foo": "abc", "bar": 3}`
	if str, err := client.PostJsonString("provision", json); err == nil {
		log.Printf("POST /provision: '%s'\n", str)
	} else {
		log.Printf("POST /provision: ERROR: %s\n", err)
	}
}
