package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/resin-os/resin-provisioner/provisioner"
	"github.com/resin-os/resin-provisioner/resin"
	"github.com/spf13/cobra"
)

func init() {
	// show date/time in log output.
	log.SetFlags(log.LstdFlags)
}

var api *provisioner.Api
var domain string

func readInput() (input string, err error) {
	i := bufio.NewReader(os.Stdin)
	in, err := i.ReadString('\n')
	if err != nil {
		return
	}
	input = strings.Trim(in, "\n")
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Compare(a, e) == 0 {
			return true
		}
	}
	return false
}

func prompt(options []string, promptMessage string) (input string, err error) {
	for {
		fmt.Printf(promptMessage)
		if input, err = readInput(); err != nil {
			return
		} else if len(options) == 0 {
			return
		} else if contains(options, input) {
			return
		}
	}
}

func login() (token string, err error) {
	fmt.Println("Logging in...")
	for {
		if email, e := prompt(nil, "email: "); err != nil {
			return "", e
		} else {
			fmt.Printf("password: ")
			if p, e := gopass.GetPasswdMasked(); e != nil {
				return "", e
			} else {
				password := string(p)
				if token, e := resin.Login("https://api."+domain, email, password); e != nil {
					return "", e
				} else if token != "" {
					return token, nil
				} else {
					fmt.Println("Wrong email or password, please try again.")
				}
			}
		}
	}
}

func signup() (token string, err error) {
	fmt.Println("Creating new user...")
	if email, e := prompt(nil, "email: "); err != nil {
		return "", e
	} else {
		for {
			fmt.Printf("password: ")
			if p, e := gopass.GetPasswdMasked(); e != nil {
				return "", e
			} else {
				fmt.Printf("confirm password: ")
				if c, e := gopass.GetPasswdMasked(); e != nil {
					return "", e
				} else {
					password := string(p)
					confirm := string(c)
					if password == confirm {
						token, err = resin.Signup("https://api."+domain, email, password)
						if err == nil && token == "" {
							return "", errors.New("Signup failed")
						} else {
							return token, nil
						}
					} else {
						fmt.Println("Passwords don't match, please try again.")
					}
				}
			}
		}
	}
}

func authenticate() (token string, err error) {
	fmt.Println("Welcome to resin.io")
	fmt.Printf(`Please log in or sign up:
	1) Log in
	2) Sign up
`)
	if input, e := prompt([]string{"1", "2"}, "> "); err != nil {
		return "", e
	} else {
		switch input {
		case "1":
			return login()
		case "2":
			return signup()
		}
	}

	return "", nil
}

func createApp(token string) (string, error) {
	for {
		if name, e := prompt(nil, "application name: "); e != nil {
			return "", e
		} else if name != "" {
			return resin.CreateApp("https://api."+domain, name, token)
		}
	}
}

func getOrCreateApp(token string) (string, error) {
	apps, err := resin.GetApps("https://api."+domain, token)
	if err != nil {
		return "", err
	}
	appIds := make([]string, len(apps))
	appList := ""
	options := make([]string, len(apps)+1)
	options[0] = "1"
	for i, app := range apps {
		appId, ok := app["id"].(float64)
		if !ok {
			return "", errors.New("Invalid app id from API")
		}
		appIds[i] = strconv.Itoa(int(appId))
		options[i+1] = strconv.Itoa(i + 2)
		appName, ok := app["app_name"].(string)
		if !ok {
			return "", errors.New("Invalid app list from API")
		}
		appList += fmt.Sprintf("\t%d) %s\n", i+2, appName)
	}
	fmt.Printf(`Choose an app for this device, or create one:
	1) Create new app
`)
	fmt.Printf(appList)
	if input, e := prompt(options, "> "); err != nil {
		return "", e
	} else {
		switch input {
		case "1":
			return createApp(token)
		default:
			i, _ := strconv.Atoi(input)
			i -= 2
			return appIds[i], nil
		}
	}

	return "", nil
}

func main() {
	var configPath string
	var dryRun bool

	rootCmd := &cobra.Command{
		Use:   "resin-provison",
		Short: "Provision this device on resin.io",
		Long: `
This command will register this device on resin.io and
start the Resin Supervisor to allow you to push applications
to this device.
It will prompt you to log in or sign up on resin and select/create
an application for this device to run.
See https://resin.io for more information about how resin.io can
help you manage device fleets.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			api = provisioner.New(configPath)
			api.Domain = domain
			if token, err := authenticate(); err != nil {
				return err
			} else if appId, err := getOrCreateApp(token); err != nil {
				return err
			} else if userId, err := resin.GetUserId(token); err != nil {
				return err
			} else if apiKey, err := resin.GetApiKey("https://api."+domain, appId, token); err != nil {
				return err
			} else {
				opts := &provisioner.ProvisionOpts{
					UserId: userId, ApplicationId: appId, ApiKey: apiKey}

				if dryRun {
					fmt.Printf("Your apikey is %s\n", apiKey)
					fmt.Printf("Ready to provision a device on appId %s for userId %s\n", appId, userId)
					return nil
				}
				if err := api.Provision(opts); err != nil {
					return err
				}

				// Since we're just returning a device URL no
				// point in worrying about the error.
				if url, err := api.DeviceUrl(); err == nil {
					fmt.Println("Your device is now provisioned and is "+
						"downloading and installing the resin supervisor.")
					fmt.Println("Your device will show as configuring during "+
						"this process, appearing online once it's complete.")
					fmt.Printf("\nYou can access the device at:\n%s\n", url)
				}

				return nil
			}
		},
	}

	p := os.Getenv("CONFIG_PATH")
	if p == "" {
		p = "/mnt/conf/config.json"
	}
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "resin.io", "Domain of the API server in which the device will register")
	rootCmd.PersistentFlags().StringVarP(&configPath, "path", "p", p, "Path for supervisor's config.json")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dryrun", "r", false, "Dry run (do not provision)")

	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Find out if this device is provisioned",
		RunE: func(cmd *cobra.Command, args []string) error {
			api := provisioner.New(configPath)
			if state, err := api.State(); err != nil {
				return err
			} else {
				fmt.Printf("This device is %s\n", state)
				return nil
			}
		},
	}
	rootCmd.AddCommand(cmdStatus)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
