package provisioner

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func reportError(status int, writer http.ResponseWriter, req *http.Request, err error) {
	log.Printf("ERROR: %s %s: %s\n", req.Method, req.URL.Path, err)

	writer.WriteHeader(status)
	fmt.Fprintf(writer, err.Error())
}

func (a *Api) provision(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if ret, err := a.getProvision(); err != nil {
			reportError(404, writer, req, err)
		} else {
			fmt.Fprintf(writer, ret)
		}

	case "POST":
		// req.Body doesn't need to be closed by us.
		if str, err := readerToString(req.Body); err != nil {
			writer.WriteHeader(501)
			fmt.Fprintf(writer,
				"Cannot convert reader to string: %s", err)
		} else {
			fmt.Fprintf(writer, "Received: '%s'", str)
		}

	default:
		reportError(400, writer, req, fmt.Errorf("Unspported method %s.", req.Method))
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
