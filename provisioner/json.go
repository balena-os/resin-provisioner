package provisioner

import (
	"encoding/json"
	"log"
)

func stringifyConfig(conf *Config) (string, error) {
	var exportedRaw map[string]interface{}

	// This forms part of the overall approach to avoid accidentally
	// stripping new JSON fields, see parseConfig() comments for more
	// details.

	// HACK: We marshal then unmarshal the config file to get a
	// map[string]interface{} representation of the exported fields so we
	// can more easily overlay the existing raw data.
	if bytes, err := json.Marshal(conf); err != nil {
		return "", err
	} else if err := json.Unmarshal(bytes, &exportedRaw); err != nil {
		return "", err
	}

	// We are mutated this field but it's fine.
	raw := conf.InitialRaw
	for name, val := range exportedRaw {
		if _, has := raw[name]; has {
			log.Printf("ivg: %s\n", name)
			raw[name] = val
		} else {
			log.Printf("WARNING: JSON: Exported '%s' not found in raw data.\n")
		}
	}

	bytes, err := json.Marshal(raw)

	return string(bytes), err
}

func parseProvisionOpts(str string) (*ProvisionOpts, error) {
	ret := new(ProvisionOpts)

	return ret, json.Unmarshal([]byte(str), ret)
}

func parseConfig(str string) (*Config, error) {
	bytes := []byte(str)
	ret := new(Config)

	// We parse it twice, firstly we populate the struct fileds in the usual
	// way:
	if err := json.Unmarshal(bytes, ret); err != nil {
		return nil, err
	}

	// Next, we populate the 'Raw' fields as map[string]interface{} values:
	if err := json.Unmarshal(bytes, &ret.InitialRaw); err != nil {
		return nil, err
	}

	// These fields will NOT be synced by default, however on stringify we
	// will overlay InitialRaw with values taken from the fields and
	// generate a combination of the two.

	// This way we avoid stripping newly created fields in config.json when
	// we only wanted to update existing known ones.

	// IMPORTANT: This won't work for any nested structs as we simply
	// overwrite exported fields in the generated output, more work would
	// need to be done to deal with that which isn't hugely worth it for the
	// provisioner.

	return ret, nil
}
