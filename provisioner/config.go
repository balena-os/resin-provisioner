package provisioner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	AppId            string `json:"applicationId,omitempty"`
	UserId           string `json:"userId,omitempty"`
	UserName         string `json:"username,omitempty"`
	DeviceType       string `json:"deviceType,omitempty"`
	ApiEndpoint      string `json:"apiEndpoint,omitempty"`
	RegistryEndpoint string `json:"registryEndpoint,omitempty"`
	VpnEndpoint      string `json:"vpnEndpoint,omitempty"`
	ApiKey           string `json:"apiKey,omitempty"`
	DeviceId         int64  `json:"deviceId,omitempty"`

	InitialRaw map[string]interface{} `json:"-"`
}

func Status(configPath string) (bool, error) {
	if config, err := readConfig(configPath); err != nil {
		if err.Error() == "Supervisor config file does not exist" {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return (config.DeviceId != 0), nil
	}
}

func Url(configPath, domain string) (string, error) {
	if config, err := readConfig(configPath); err != nil {
		return "", err
	} else {
		url := fmt.Sprintf("https://dashboard.%s/apps/%s/devices/%d/summary",
			domain, config.AppId, config.DeviceId)
		return url, nil
	}
}

func readConfig(configPath string) (Config, error) {
	var config Config
	if bytes, err := ioutil.ReadFile(configPath); os.IsNotExist(err) {
		return config, errors.New("Supervisor config file does not exist")
	} else if err != nil {
		return config, err
	} else if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	} else if err := unmarshalRawSafe(string(bytes), &config.InitialRaw); err != nil {
		return config, err
	} else {
		return config, nil
	}
}

// Avoid screwing up numbers when unmarshalling (i.e. decoding).
// See http://stackoverflow.com/a/22346593
func unmarshalRawSafe(str string, out *map[string]interface{}) error {
	dec := json.NewDecoder(strings.NewReader(str))
	dec.UseNumber()
	return dec.Decode(&out)
}

func marshal(config Config) (string, error) {
	var exportedRaw map[string]interface{}

	// HACK: We marshal then unmarshal the config file to get a
	// map[string]interface{} representation of the exported fields so we
	// can more easily overlay the existing raw data.
	if bytes, err := json.Marshal(config); err != nil {
		return "", err
	} else if err := unmarshalRawSafe(string(bytes), &exportedRaw); err != nil {
		return "", err
	}

	// We are mutating this field but it's fine.
	raw := config.InitialRaw
	for name, val := range exportedRaw {
		raw[name] = val
	}

	bytes, err := json.Marshal(raw)

	return string(bytes), err
}
