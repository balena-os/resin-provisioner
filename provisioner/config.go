package provisioner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ApplicationId         string  `json:"applicationId"`
	ApiKey                string  `json:"apiKey"`
	UserId                string  `json:"userId"`
	DeviceType            string  `json:"deviceType"`
	RegisteredAt          float64 `json:"registered_at,omitempty"`
	AppUpdatePollInterval string  `json:"appUpdatePollInterval"`
	ListenPort            string  `json:"listenPort"`
	VpnPort               string  `json:"vpnPort"`
	ApiEndpoint           string  `json:"apiEndpoint"`
	VpnEndpoint           string  `json:"vpnEndpoint"`
	RegistryEndpoint      string  `json:"registryEndpoint"`
	DeltaEndpoint         string  `json:"deltaEndpoint"`
	PubnubSubscribeKey    string  `json:"pubnubSubscribeKey"`
	PubnubPublishKey      string  `json:"pubnubPublishKey"`
	MixpanelToken         string  `json:"mixpanelToken"`

	// See json.go/parseConfig() for more details on what this is for.
	InitialRaw map[string]interface{} `json:"-"`
}

func (a *Api) readConfig() (*Config, error) {
	if bytes, err := ioutil.ReadFile(a.ConfigPath); os.IsNotExist(err) {
		// We'll create a new config.json.
		log.Printf("Empty %s, will create new on write.\n", a.ConfigPath)
		return parseConfig("{}")
	} else if err != nil {
		return nil, err
	} else {
		return parseConfig(string(bytes))
	}
}

func (a *Api) writeConfig(conf *Config) error {
	if str, err := stringifyConfig(conf); err != nil {
		return err
	} else {
		return atomicWrite(a.ConfigPath, str)
	}
}

func (c *Config) ProvisionedState() ProvisionedState {
	if c.ApplicationId == "" || c.ApiKey == "" {
		return Unprovisioned
	}

	if c.RegisteredAt == 0 {
		return Provisioning
	}

	return Provisioned
}

// If DeviceType not specified, attempt to determine it and assign.
func (c *Config) DetectDeviceType() error {
	if c.DeviceType != "" {
		return nil
	}

	if deviceType, err := ScanDeviceTypeSlug(); err != nil {
		return err
	} else {
		c.DeviceType = deviceType
	}

	return nil
}

// Get /config from the Resin API specified at c.ApiEndpoint
func (c Config) getConfigFromApi() (map[string]interface{}, error) {
	var conf map[string]interface{}
	conf = make(map[string]interface{})
	if r, err := getUrl(c.ApiEndpoint + "/config"); err != nil {
		return nil, err
	} else if err = json.Unmarshal(r, &conf); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}

// Get and populate mixpanel and pubnub keys from the Resin API
func (c *Config) GetKeysFromApi() error {
	// GET /config from api
	if conf, err := c.getConfigFromApi(); err != nil {
		return fmt.Errorf("Error getting config from Resin API: %s", err)
	} else {
		i := errors.New("Invalid config received from the Resin API")
		if t, ok := conf["mixpanelToken"].(string); !ok {
			return i
		} else if p, ok := conf["pubnub"].(map[string]interface{}); !ok {
			return i
		} else if pk, ok := p["publish_key"].(string); !ok {
			return i
		} else if sk, ok := p["subscribe_key"].(string); !ok {
			return i
		} else if t == "" || pk == "" || sk == "" {
			return i
		} else {
			c.MixpanelToken = t
			c.PubnubPublishKey = pk
			c.PubnubSubscribeKey = sk
			return nil
		}
	}
}

// Read environment-field specified values.
func (c *Config) ReadEnv() {
	// TODO: Deduplicate from defaults.go.
	if domainOverride := os.Getenv(DOMAIN_OVERRIDE_ENV_VAR); domainOverride != "" {
		c.ApiEndpoint = fmt.Sprintf("https://api.%s", domainOverride)
		c.VpnEndpoint = fmt.Sprintf("vpn.%s", domainOverride)
		c.RegistryEndpoint = fmt.Sprintf("registry.%s", domainOverride)
		c.DeltaEndpoint = fmt.Sprintf("https://delta.%s", domainOverride)
	}
}
