package main

import (
	"fmt"
	"log"
	"os"

	"github.com/resin-os/resin-provisioner/provisioner"
	"github.com/spf13/cobra"
)

func init() {
	// show date/time in log output.
	log.SetFlags(log.LstdFlags)
}

func authenticate() (string, error) {
	return "", nil
}

func getOrCreateApp(token string) (string, error) {
	return "", nil
}

func getUserId(token string) (string, error) {
	return "", nil
}

func getApiKey(token, appId string) (string, error) {
	return "", nil
}

func main() {
	var configPath string
	var domain string

	rootCmd := &cobra.Command{
		Use: "resin-provison"
		Short: "Provision this device on resin.io"
		Long: `
			This command will register this device on resin.io and
			start the Resin Supervisor to allow you to push applications
			to this device.
			It will prompt you to log in or sign up on resin and select/create
			an application for this device to run.
			See https://resin.io for more information about how resin.io can
			help you manage device fleets.
		`
		RunE: func(cmd *cobra.Command, args []string) error {
			if token, err := authenticate(); err != nil {
				return err
			} else if appId, err := getOrCreateApp(token); err != nil {
				return err
			} else if userId, err := getUserId(token); err != nil {
				return err
			} else if apiKey, err := getApiKey(token, appId); err != nil {
				return err
			} else {
				api := provisioner.New(configPath)
				opts := &provisioner.ProvisionOpts{
					UserId: userId, ApplicationId: appId, ApiKey: apiKey}

				if err := api.Provision(opts); err != nil {
					return err
				}
				return nil
			}
		}
	}

	p := os.Getenv("CONFIG_PATH")
	if p == "" {
		p = "/mnt/conf/config.json"
	}
	rootCmd.PersistentFlags().StringVarP(&domain,"domain", "d", "resin.io", "Domain of the API server in which the device will register")
	rootCmd.PersistentFlags().StringVarP(&configPath,"path", "p", p, "Path for supervisor's config.json")

	cmdStatus := &cobra.Command{
		Use: "status"
		Short: "Find out if this device is provisioned"
		RunE: func(cmd *cobra.Command, args []string) error {
			api := provisioner.New(configPath)
			if state, err := api.State(); err != nil {
				return err
			} else {
				fmt.Printf("This device is %s\n", state)
			}
		}
	}
	rootCmd.AddCommand(cmdStatus)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
