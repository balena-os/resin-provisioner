package provisioner

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) initServer() {
	router := mux.NewRouter()

	router.HandleFunc("/provisioned", a.provisionedHandler).Methods("GET")
	router.HandleFunc("/provision", a.provisionHandler).Methods("GET", "POST")
	router.HandleFunc("/config", a.configHandler).Methods("GET")

	a.server = &http.Server{Handler: router}
}

func (a *Api) Serve(path string) (err error) {
	if err = checkSocket(path); err != nil {
		return
	}

	if a.listener, err = net.Listen("unix", path); err == nil {
		err = a.server.Serve(a.listener)
	}

	return
}

// Clean up the socket after use.
func (a *Api) Cleanup() {
	if a != nil && a.listener != nil {
		// This will clear down the socket file.
		a.listener.Close()
	}
}
