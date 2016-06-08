package provisioner

import (
	"io"
	"io/ioutil"
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
