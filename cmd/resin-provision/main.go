package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/resin-os/resin-provisioner/defaults"
	"github.com/resin-os/resin-provisioner/provisioner"
	"github.com/resin-os/resin-provisioner/resin"
	"github.com/spf13/cobra"
)

var domain, configPath, application, email, password, token string
var dryRun bool

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {

	rootCmd := &cobra.Command{
		Use:   "resin-provision",
		Short: "Provision this device on resin.io",
		Long: `
This tool is a resin component for converting an unmanaged resinOS device into a
managed resinOS device. This is achieved by setting the configuration needed for
the supervisor to provision the device against the resin.io servers.
See https://resin.io for more information about how resin.io can help
you manage device fleets.`,
	}

	oneshotCmd := &cobra.Command{
		Use:   "oneshot",
		Short: "One-shot mode",
		Long: `
Provision the device with a single command, for example:
./resin-provision oneshot -e email@resin.io -p secret_password -a testApplication`,
		PreRunE:  validate,
		RunE:     oneshot,
		PostRunE: success,
	}

	interactiveCmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode",
		Long: `
Provision the device with a series of interactive prompts. Use
this option if you need to create a new account or application.`,
		RunE:     interactive,
		PostRunE: success,
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check whether this device is provisioned",
		RunE:  status,
	}

	rootCmd.AddCommand(oneshotCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(statusCmd)

	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", defaults.RESIN_DOMAIN, "Domain the device will provision to")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", defaults.CONFIG_PATH, "Supervisor config path")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dryrun", "r", false, "Dry run (do not provision)")
	oneshotCmd.Flags().StringVarP(&application, "application", "a", "", "Application the device will provision to")
	oneshotCmd.Flags().StringVarP(&email, "email", "e", "", "Email addres of the account the device will provision to")
	oneshotCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the account the device will provision to")
	oneshotCmd.Flags().StringVarP(&token, "token", "t", "", "User AUTH token")

	return rootCmd.Execute()
}

func validate(cmd *cobra.Command, args []string) error {
	if application == "" {
		return errors.New("Application is required")
	} else if token == "" {
		if email == "" {
			return errors.New("Email is required unless an AUTH token is used")
		} else if password == "" {
			return errors.New("Password is required unless an AUTH token is used")
		}
	}

	return nil
}

func oneshot(cmd *cobra.Command, args []string) error {
	if status, err := provisioner.Status(configPath); err != nil {
		return err
	} else if status {
		return nil
	} else if token, err := getToken(); err != nil {
		return err
	} else if token == "" {
		return errors.New("Wrong email or password, please try again")
	} else if appId, err := getApp(token, application); err != nil {
		return err
	} else if userId, err := resin.GetUserId(token); err != nil {
		return err
	} else if userName, err := resin.GetUserName(token); err != nil {
		return err
	} else if apiKey, err := resin.GetApiKey("https://api."+domain, appId, token); err != nil {
		return err
	} else {
		return provisioner.Provision(appId, apiKey, userId, userName, configPath, domain, dryRun)
	}
}

func interactive(cmd *cobra.Command, args []string) error {
	if status, err := provisioner.Status(configPath); err != nil {
		return err
	} else if status {
		return nil
	} else if token, err := authenticate(); err != nil {
		return err
	} else if appId, err := getOrCreateApp(token); err != nil {
		return err
	} else if userId, err := resin.GetUserId(token); err != nil {
		return err
	} else if userName, err := resin.GetUserName(token); err != nil {
		return err
	} else if apiKey, err := resin.GetApiKey("https://api."+domain, appId, token); err != nil {
		return err
	} else {
		return provisioner.Provision(appId, apiKey, userId, userName, configPath, domain, dryRun)
	}
}

func success(cmd *cobra.Command, args []string) error {
	fmt.Println("Your device is now provisioned and is " +
		"downloading and installing the resin supervisor")
	fmt.Println("Your device may show as configuring during " +
		"this process, appearing online once it's complete")

	return nil
}

func status(cmd *cobra.Command, args []string) error {
	if status, err := provisioner.Status(configPath); err != nil {
		return err
	} else if status {
		if url, err := provisioner.Url(configPath, domain); err != nil {
			return err
		} else {
			fmt.Println("Your device is provisioned")
			fmt.Printf("You can access the device at: %s\n", url)
		}
	} else {
		fmt.Println("Your device is not provisioned")
	}

	return nil
}

func getToken() (string, error) {
	if token == "" {
		return resin.Login("https://api."+domain, email, password)
	} else {
		return token, nil
	}
}

func getOrCreateApp(token string) (string, error) {
	apps, err := resin.GetApps("https://api."+domain, token)
	if err != nil {
		return "", err
	}

	options := make([]string, 0)
	list := make([]string, 0)

	options = append(options, "1")
	list = append(list, "Create new app")

	for index, app := range apps {
		appName, ok := app["app_name"].(string)
		if !ok {
			return "", errors.New("Invalid app name from API")
		}

		list = append(list, appName)
		options = append(options, strconv.Itoa(index+2))
	}

	fmt.Println("Choose an app for this device, or create one:")
	for index := range options {
		fmt.Printf("%s) %s", options[index], list[index])

		if index < len(options)-1 {
			fmt.Printf(" \n")
		}
	}

	if input, err := prompt(options, "\n> "); err != nil {
		return "", err
	} else {
		switch input {
		case "1":
			return createApp(token)
		default:
			if index, err := strconv.Atoi(input); err != nil {
				return "", err
			} else {
				return getApp(token, list[index-1])
			}
		}
	}
}

func getApp(token, appName string) (string, error) {
	apps, err := resin.GetApps("https://api."+domain, token)
	if err != nil {
		return "", err
	}

	for _, app := range apps {
		if app["app_name"].(string) == appName {
			appId, ok := app["id"].(float64)
			if !ok {
				return "", errors.New("Invalid app ID from API")
			}

			return strconv.Itoa(int(appId)), nil
		}
	}

	return "", errors.New("Application not found")
}

func createApp(token string) (string, error) {
	for {
		if name, err := prompt(nil, "application name: "); err != nil {
			return "", err
		} else if name != "" {
			return resin.CreateApp("https://api."+domain, name, token)
		}
	}
}

func authenticate() (string, error) {
	fmt.Println("Welcome to resin.io")
	fmt.Printf(`Please log in or sign up:
1) Log in
2) Sign up`)

	if input, err := prompt([]string{"1", "2"}, "\n> "); err != nil {
		return "", err
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

func login() (string, error) {
	fmt.Println("Logging in...")
	for {
		if email, err := prompt(nil, "email: "); err != nil {
			return "", err
		} else {
			fmt.Printf("password: ")
			if p, err := gopass.GetPasswdMasked(); err != nil {
				return "", err
			} else {
				password := string(p)
				if token, err := resin.Login("https://api."+domain, email, password); err != nil {
					return "", err
				} else if token != "" {
					return token, nil
				} else {
					fmt.Println("Wrong email or password, please try again")
				}
			}
		}
	}
}

func signup() (string, error) {
	fmt.Println("Creating new user...")
	if email, err := prompt(nil, "email: "); err != nil {
		return "", err
	} else {
		for {
			fmt.Printf("password: ")
			if p, err := gopass.GetPasswdMasked(); err != nil {
				return "", err
			} else {
				fmt.Printf("confirm password: ")
				if c, err := gopass.GetPasswdMasked(); err != nil {
					return "", err
				} else {
					password := string(p)
					confirm := string(c)
					if password == confirm {
						token, err := resin.Signup("https://api."+domain, email, password)
						if err == nil && token == "" {
							return "", errors.New("Signup failed")
						} else {
							return token, nil
						}
					} else {
						fmt.Println("Passwords don't match, please try again")
					}
				}
			}
		}
	}
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
		} else {
			fmt.Printf("Not a valid option")
		}
	}
}

func readInput() (string, error) {
	i := bufio.NewReader(os.Stdin)
	if in, err := i.ReadString('\n'); err != nil {
		return "", err
	} else {
		return strings.Trim(in, "\n"), nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Compare(a, e) == 0 {
			return true
		}
	}
	return false
}
