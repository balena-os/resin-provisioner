package provisioner

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Provisioner) provision(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprintf(writer, "{}")

	case "POST":
		// req.Body doesn't need to be closed by us.
		if str, err := readerToString(req.Body); err != nil {
			writer.WriteHeader(501)
			fmt.Fprintf(writer, "Cannot decode: %s", err)
		} else {
			fmt.Fprintf(writer, "Received: '%s'", str)
		}

	default:
		writer.WriteHeader(400)
	}
}

func (p *Provisioner) initApi() {
	router := mux.NewRouter()

	router.HandleFunc("/provision", p.provision).Methods("GET", "POST")

	p.server = &http.Server{Handler: router}
}

func (p *Provisioner) Serve(path string) (err error) {
	if err = checkSocket(path); err != nil {
		return
	}

	if p.listener, err = net.Listen("unix", path); err == nil {
		err = p.server.Serve(p.listener)
	}

	return
}

// Clean up the socket after use.
func (p *Provisioner) Cleanup() {
	if p != nil && p.listener != nil {
		// This will clear down the socket file.
		p.listener.Close()
	}
}
