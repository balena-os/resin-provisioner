package main

import (
	"fmt"
	"log"
	"os"

	"github.com/resin-os/resin-provisioner/provisioner"
)

func init() {
	// show date/time in log output.
	log.SetFlags(log.LstdFlags)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [config path] [user id] [application id] [api key]\n",
		os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 5 {
		usage()
	}

	configPath := os.Args[1]
	userId := os.Args[2]
	appId := os.Args[3]
	apiKey := os.Args[4]

	opts := &provisioner.ProvisionOpts{
		UserId: userId, ApplicationId: appId, ApiKey: apiKey}

	api := provisioner.New(configPath)

	if err := api.Provision(opts); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
