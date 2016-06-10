package provisioner

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func checkSocket(path string) error {
	// The socket file not existing means we can create it.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Otherwise, remove the socket ready for re-creation. If we fail here,
	// just pass on the error.
	return os.Remove(path)
}

func readerToString(r io.Reader) (ret string, err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadAll(r); err == nil {
		ret = string(bytes)
	}

	return
}

func atomicWrite(path, content string) error {
	// Try to be atomic - write output to temporary file, sync it, rename it
	// to the target file.

	// "" is the dir - defaults to system temp dir, "provisioner" is prefix.
	if tmpFile, err := ioutil.TempFile("", "provisioner"); err != nil {
		return err
	} else {
		name := tmpFile.Name()

		// We ignore the error so removing the now renamed file isn't an
		// issue. In error cases we clean up.
		defer os.Remove(name)

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			return err
		} else if err := tmpFile.Sync(); err != nil {
			return err
		} else if err := tmpFile.Close(); err != nil {
			return err
		}

		return os.Rename(name, path)
	}
}

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
