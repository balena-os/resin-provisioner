package provisioner

import "encoding/json"

func stringify(val interface{}) (ret string, err error) {
	var bytes []byte

	if bytes, err = json.Marshal(val); err == nil {
		ret = string(bytes)
	}

	return
}

func parseProvisionOpts(str string) (*ProvisionOpts, error) {
	ret := new(ProvisionOpts)

	return ret, json.Unmarshal([]byte(str), ret)
}

func parseConfig(str string) (*Config, error) {
	ret := new(Config)

	return ret, json.Unmarshal([]byte(str), ret)
}
