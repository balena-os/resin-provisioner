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

func New(configPath string) *Api {
	ret := &Api{ConfigPath: configPath}
	ret.initSocket()

	return ret
}

func (a *Api) getProvision() (string, error) {
	if conf, err := a.readConfig(); err != nil {
		return "", fmt.Errorf("Cannot read config: %s", err)
	} else if str, err := stringify(conf); err != nil {
		return "", fmt.Errorf("Cannot stringfy config: %s", err)
	} else {
		return str, nil
	}
}
