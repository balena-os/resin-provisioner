package provisioner

import (
	"fmt"
	"log"
	"net/http"
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
