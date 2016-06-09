package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/resin-os/resin-provisioner/provisioner"
)

var api *provisioner.Api

func init() {
	// show date/time in log output.
	log.SetFlags(log.LstdFlags)
}

func handleSignals() {
	in := make(chan os.Signal, 1)
	signal.Notify(in,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGABRT,
		syscall.SIGQUIT)

	go func() {
		// Any of the masked signals mean death, so only need to catch
		// once.
		<-in

		api.Cleanup()
		os.Exit(1)
	}()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [socket path] [config.json path]\n",
		os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	socketPath, configPath := os.Args[1], os.Args[2]

	api = provisioner.New(configPath)
	handleSignals()

	log.Printf("Started.")
	api.Serve(socketPath)
}
