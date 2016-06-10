package provisioner

import "io/ioutil"

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

func (a *Api) readConfig() (*Config, error) {
	if bytes, err := ioutil.ReadFile(a.ConfigPath); err != nil {
		return nil, err
	} else {
		return parseConfig(string(bytes))
	}
}

func (a *Api) writeConfig(conf *Config) error {
	if str, err := stringify(conf); err != nil {
		return err
	} else {
		return atomicWrite(a.ConfigPath, str)
	}
}
