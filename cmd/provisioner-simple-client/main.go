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
	fmt.Fprintf(os.Stderr, "usage:     query: %s [config path]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "usage: provision: %s [config path] [user id] [app id] [api key]\n",
		os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 5 && len(os.Args) != 2 {
		usage()
	}

	configPath := os.Args[1]
	api := provisioner.New(configPath)

	// Simply output the provision state.
	if len(os.Args) == 2 {
		if state, err := api.State(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("%s\n", state)
		}

		return
	}

	userId := os.Args[2]
	appId := os.Args[3]
	apiKey := os.Args[4]

	opts := &provisioner.ProvisionOpts{
		UserId: userId, ApplicationId: appId, ApiKey: apiKey}

	if err := api.Provision(opts); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
