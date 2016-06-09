package provisioner

import (
	"net"
	"net/http"
)

type Provisioner struct {
	ConfigPath string
	listener net.Listener
	server   *http.Server
}

func New(configPath string) *Provisioner {
	ret := &Provisioner{ConfigPath:configPath}
	ret.initApi()

	return ret
}
