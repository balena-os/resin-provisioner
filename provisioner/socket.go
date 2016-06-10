package provisioner

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) provision(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if conf, err := a.readConfig(); err != nil {
			writer.WriteHeader(404)
			fmt.Fprintf(writer, "Cannot read config: %s", err)
		} else if str, err := stringifyConfig(conf); err != nil {
			writer.WriteHeader(404)
			fmt.Fprintf(writer, "Cannot stringify config: %s", err)
		} else {
			fmt.Fprintf(writer, str)
		}

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

func (a *Api) initSocket() {
	router := mux.NewRouter()

	router.HandleFunc("/provision", a.provision).Methods("GET", "POST")

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
