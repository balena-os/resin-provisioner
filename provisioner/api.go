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
	ret.initSocket()

	return ret
}

func (a *Api) IsProvisioned() (bool, error) {
	if conf, err := a.readConfig(); err != nil {
		return false, fmt.Errorf("Cannot read config: %s", err)
	} else {
		return conf.IsProvisioned(), nil
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
	return nil
}

func (a *Api) ConfigJson() (string, error) {
	if conf, err := a.readConfig(); err != nil {
		return "", fmt.Errorf("Cannot read config: %s", err)
	} else if str, err := stringify(conf); err != nil {
		return "", fmt.Errorf("Cannot stringfy config: %s", err)
	} else {
		return str, nil
	}
}
