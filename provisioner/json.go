package provisioner

import (
	"encoding/json"
	"io/ioutil"
)

type ProvisionOpts struct {
	UserId        string `json:"userId"`
	ApplicationId string `json:"applicationId"`
	ApiKey        string `json:"apikey"`
}

type Config struct {
	ApplicationId string  `json:"applicationId"`
	ApiKey        string  `json:"apikey"`
	UserId        string  `json:"userId"`
	Username      string  `json:"username"`
	DeviceType    string  `json:"deviceType"`
	Uuid          string  `json:"uuid,omitempty"`
	RegisteredAt  float64 `json:"registered_at,omitempty"`
	DeviceId      float64 `json:"deviceId,omitempty"`
}

func parseProvisionOpts(str string) (opts *ProvisionOpts, err error) {
	opts = new(ProvisionOpts)
	err = json.Unmarshal([]byte(str), opts)

	return
}

func (a *Api) readConfig() (conf *Config, err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadFile(a.ConfigPath); err == nil {
		conf = new(Config)
		err = json.Unmarshal(bytes, conf)
	}

	return
}

func stringifyConfig(conf *Config) (ret string, err error) {
	var bytes []byte

	if bytes, err = json.Marshal(conf); err == nil {
		ret = string(bytes)
	}

	return
}

func (a *Api) writeConfig(conf *Config) error {
	if str, err := stringifyConfig(conf); err != nil {
		return err
	} else {
		return atomicWrite(a.ConfigPath, str)
	}
}
