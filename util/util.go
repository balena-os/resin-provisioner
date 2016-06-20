package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pathLib "path"
	"strings"
)

const (
	OSRELEASE_PATH = "/etc/os-release"
)

// Try to be atomic - write output to temporary file, sync it, rename it
// to the target file.
func AtomicWrite(path, content string) error {
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
func ReadLines(path string) ([]string, error) {
	var (
		bytes []byte
		err   error
	)

	if bytes, err = ioutil.ReadFile(path); err != nil {
		return nil, fmt.Errorf("Can't read %s: %s", path, err)
	}

	str := strings.TrimSpace(string(bytes))
	rawLines := strings.Split(str, "\n")

	ret := make([]string, len(rawLines))
	for i, line := range rawLines {
		ret[i] = strings.TrimSpace(line)
	}

	return ret, nil
}

func GetEnvFileFields(path string) (map[string]string, error) {
	var (
		lines []string
		err   error
	)

	if lines, err = ReadLines(path); err != nil {
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

func SetEnvFileFields(path string, fields map[string]string) error {
	lines := make([]string, 0, len(fields))

	for key, val := range fields {
		lines = append(lines, fmt.Sprintf("%s=%s", key, val))
	}

	// Add trailing newline too.
	str := strings.Join(lines, "\n") + "\n"

	return AtomicWrite(path, str)
}

func getOsReleaseSlug(path string) (string, error) {
	if fields, err := GetEnvFileFields(path); err != nil {
		return "", err
	} else {
		slug := fields["SLUG"]
		if slug == "" {
			log.Printf("Could not find 'SLUG' field in %s, defaulting to raspberry-pi\n",
				path)
			slug = "raspberry-pi"
		}

		return slug, nil
	}
}

func ScanDeviceTypeSlug(path string) (string, error) {
	return getOsReleaseSlug(path)
}
