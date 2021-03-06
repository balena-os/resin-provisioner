package provisioner

import (
	"encoding/json"
	"strings"
)

// Avoid screwing up numbers when unmarshalling (i.e. decoding).
// See http://stackoverflow.com/a/22346593
func unmarshalRawSafe(str string, out *map[string]interface{}) error {
	reader := strings.NewReader(str)
	decoder := json.NewDecoder(reader)
	// The key to fixing the issue. Forces decoding into json.Number rather
	// than a possible float64 if it feels like it.
	decoder.UseNumber()

	return decoder.Decode(out)
}

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
	} else if err := unmarshalRawSafe(string(bytes), &exportedRaw); err != nil {
		return "", err
	}

	// We are mutating this field but it's fine.
	raw := conf.InitialRaw
	for name, val := range exportedRaw {
		// TODO: Deal with nested structs correctly.
		raw[name] = val
	}

	bytes, err := json.Marshal(raw)

	return string(bytes), err
}

func parseProvisionOpts(str string) (*ProvisionOpts, error) {
	ret := new(ProvisionOpts)

	return ret, json.Unmarshal([]byte(str), ret)
}

func parseConfig(str string, domain string) (*Config, error) {
	bytes := []byte(str)
	ret := new(Config)
	*ret = DefaultConfig

	// We parse it twice, firstly we populate the struct fileds in the usual
	// way:
	if err := json.Unmarshal(bytes, ret); err != nil {
		return nil, err
	}

	// Next, we populate the 'Raw' fields as map[string]interface{} values:
	if err := unmarshalRawSafe(string(bytes), &ret.InitialRaw); err != nil {
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

	// If DeviceType not specified, attempt to detect it.
	if err := ret.DetectDeviceType(); err != nil {
		return nil, err
	}

	ret.SetDomain(domain)

	return ret, nil
}
