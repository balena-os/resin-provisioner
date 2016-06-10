package provisioner

import (
	"fmt"
	"net/http"
)

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
		if str := readPostBodyReportErr(writer, req); str == "" {
			return
		} else if opts, err := parseProvisionOpts(str); err != nil {
			reportError(404, writer, req, err)
		} else if err := a.Provision(opts); err != nil {
			reportError(404, writer, req, err)
		}

	default:
		// Shouldn't be possible.
		reportError(400, writer, req,
			fmt.Errorf("Unspported method %s.", req.Method))
	}
}

func (a *Api) configHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		// Shouldn't be possible.
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
