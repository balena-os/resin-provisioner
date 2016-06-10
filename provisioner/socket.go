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

func readPostBodyReportErr(writer http.ResponseWriter, req *http.Request) string {
	// req.Body doesn't need to be closed by us.
	if str, err := readerToString(req.Body); err != nil {
		reportError(500, writer, req,
			fmt.Errorf("Cannot convert reader to string: %s", err))

		return ""
	} else {
		return str
	}
}

func (a *Api) provisionedHandler(writer http.ResponseWriter, req *http.Request) {
	if str, err := a.IsProvisionedJson(); err != nil {
		reportError(404, writer, req, err)
	} else {
		fmt.Fprintf(writer, str)
	}
}

func (a *Api) provisionHandler(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		a.provisionedHandler(writer, req)
	case "POST":
		if str := readPostBodyReportErr(writer, req); str != "" {
			fmt.Fprintf(writer, "Received: '%s'", str)
		}

	default:
		reportError(400, writer, req,
			fmt.Errorf("Unspported method %s.", req.Method))
	}
}

func (a *Api) configHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		reportError(400, writer, req,
			fmt.Errorf("Unspported method %s.", req.Method))
		return
	}

	if str, err := a.ConfigJson(); err != nil {
		reportError(404, writer, req, err)
	} else {
		fmt.Fprintf(writer, str)
	}
}

func (a *Api) initSocket() {
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
