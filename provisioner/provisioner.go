package provisioner

import (
	"net"
	"net/http"
)

type Provisioner struct {
	listener net.Listener
	server   *http.Server
}

func New() *Provisioner {
	ret := &Provisioner{}
	ret.initApi()

	return ret
}
