package provisioner

import (
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
