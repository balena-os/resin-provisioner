package provisioner

import (
	"fmt"
	"net"
	"net/http"
)

type Api struct {
	ConfigPath string
	listener   net.Listener
	server     *http.Server
}

type ProvisionOpts struct {
	UserId        string `json:"userId"`
	ApplicationId string `json:"applicationId"`
	ApiKey        string `json:"apikey"`
}

func New(configPath string) *Api {
	ret := &Api{ConfigPath: configPath}
	ret.initServer()

	return ret
}

func (a *Api) IsProvisioned() (bool, error) {
	if conf, err := a.readConfig(); err != nil {
		return false, fmt.Errorf("Cannot read config: %s", err)
	} else if !conf.IsProvisioned() {
		return false, nil
	} else {
		return supervisorDbusRunning()
	}
}

func (a *Api) IsProvisionedJson() (ret string, err error) {
	var provisioned bool

	if provisioned, err = a.IsProvisioned(); err == nil {
		ret = fmt.Sprintf(`{"provisioned": %t}`, provisioned)
	}

	return
}

func (a *Api) Provision(opts *ProvisionOpts) error {
	if provisioned, err := a.IsProvisioned(); err != nil {
		return err
	} else if provisioned {
		return fmt.Errorf("Already provisioned.")
	}

	if !isInteger(opts.UserId) || !isInteger(opts.ApplicationId) ||
		!isValidApiKey(opts.ApiKey) {
		return fmt.Errorf("Invalid options.")
	}

	if conf, err := a.readConfig(); err != nil {
		return fmt.Errorf("Cannot read config: %s", err)
	} else {
		// First check to see whether config.json has changed from
		// underneath us.
		if conf.IsProvisioned() {
			if running, err := supervisorDbusRunning(); err != nil {
				return err
			} else if running {
				return fmt.Errorf("Already provisioned.")
			}
		}

		// Ok, now we go for it.

		conf.UserId = opts.UserId
		conf.ApplicationId = opts.ApplicationId
		conf.ApiKey = opts.ApiKey

		if err := a.writeConfig(conf); err != nil {
			return err
		}

		// Next we need to enable the supervisor systemd service.

		return nil
	}
}

func (a *Api) ConfigJson() (string, error) {
	if conf, err := a.readConfig(); err != nil {
		return "", fmt.Errorf("Cannot read config: %s", err)
	} else if str, err := stringifyConfig(conf); err != nil {
		return "", fmt.Errorf("Cannot stringfy config: %s", err)
	} else {
		return str, nil
	}
}
