package provisioner

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	pathLib "path"
	"regexp"
	"strconv"
	"strings"
)

var apiKeyRegexp = regexp.MustCompile("[a-zA-Z0-9]+")

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

// Try to be atomic - write output to temporary file, sync it, rename it
// to the target file.
func atomicWrite(path, content string) error {
	// To avoid cross-device rename issues, create the temporary file in
	// config.json's containing directory.
	targetDir := pathLib.Dir(path)

	if tmpFile, err := ioutil.TempFile(targetDir, "provisioner"); err != nil {
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

func reportError(status int, writer http.ResponseWriter, req *http.Request,
	err error, userErr string) {
	log.Printf("ERROR: %s %s: %s (%s)\n", req.Method, req.URL.Path, err,
		userErr)

	writer.WriteHeader(status)
	fmt.Fprintf(writer, "ERROR: %s", userErr)
}

func readPostBodyReportErr(writer http.ResponseWriter, req *http.Request) string {
	// req.Body doesn't need to be closed by us.
	if str, err := readerToString(req.Body); err != nil {
		reportError(500, writer, req, err,
			"Can't convert read to string")

		return ""
	} else {
		return str
	}
}

func isInteger(str string) bool {
	_, err := strconv.Atoi(str)

	return err == nil
}

func isValidApiKey(str string) bool {
	return apiKeyRegexp.Match([]byte(str))
}

func supervisorDbusRunning() (bool, error) {
	if dbus, err := NewDbus(); err != nil {
		return false, err
	} else {
		defer dbus.Close()

		return dbus.SupervisorRunning()
	}
}

func readLines(path string) ([]string, error) {
	var (
		bytes []byte
		err   error
	)

	if bytes, err = ioutil.ReadFile(path); err != nil {
		return nil, fmt.Errorf("Can't read %s: %s", path, err)
	}

	str := string(bytes)
	rawLines := strings.Split(str, "\n")

	ret := make([]string, len(rawLines))
	for i, line := range rawLines {
		ret[i] = strings.TrimSpace(line)
	}

	return ret, nil
}

func getEnvFileFields(path string) (map[string]string, error) {
	var (
		lines []string
		err   error
	)

	if lines, err = readLines(path); err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	for _, line := range lines {
		fields := strings.Split(line, "=")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		val := strings.TrimSpace(fields[1])

		ret[key] = val
	}

	return ret, nil
}

func setEnvFileFields(path string, fields map[string]string) error {
	lines := make([]string, 0, len(fields))

	for key, val := range fields {
		lines = append(lines, fmt.Sprintf("%s=%s", key, val))
	}

	// Add trailing newline too.
	str := strings.Join(lines, "\n") + "\n"

	return atomicWrite(path, str)
}

func setSupervisorTag() error {
	if fields, err := getEnvFileFields(SUPERVISOR_CONF_PATH); err != nil {
		return err
	} else {
		fields["SUPERVISOR_TAG"] = INIT_UPDATER_SUPERVISOR_TAG

		return setEnvFileFields(SUPERVISOR_CONF_PATH, fields)
	}
}
