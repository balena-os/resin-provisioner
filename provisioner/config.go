package provisioner

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/resin-os/resin-provisioner/resin"
	"github.com/resin-os/resin-provisioner/util"
)

type Config struct {
	Uuid                  string `json:"uuid"`
	ApplicationId         string `json:"applicationId"`
	ApiKey                string `json:"apiKey"`
	UserId                string `json:"userId"`
	UserName              string `json:"username"`
	DeviceId              int64  `json:"deviceId",omitempty`
	DeviceType            string `json:"deviceType"`
	RegisteredAt          int64  `json:"registered_at,omitempty"`
	AppUpdatePollInterval string `json:"appUpdatePollInterval"`
	ListenPort            string `json:"listenPort"`
	VpnPort               string `json:"vpnPort"`
	ApiEndpoint           string `json:"apiEndpoint"`
	VpnEndpoint           string `json:"vpnEndpoint"`
	RegistryEndpoint      string `json:"registryEndpoint"`
	DeltaEndpoint         string `json:"deltaEndpoint"`
	PubnubSubscribeKey    string `json:"pubnubSubscribeKey"`
	PubnubPublishKey      string `json:"pubnubPublishKey"`
	MixpanelToken         string `json:"mixpanelToken"`

	// See json.go/parseConfig() for more details on what this is for.
	InitialRaw map[string]interface{} `json:"-"`
}

func (a *Api) readConfig() (*Config, error) {
	if bytes, err := ioutil.ReadFile(a.ConfigPath); os.IsNotExist(err) {
		// We'll create a new config.json.
		return parseConfig("{}", a.Domain)
	} else if err != nil {
		return nil, err
	} else {
		return parseConfig(string(bytes), a.Domain)
	}
}

func (a *Api) writeConfig(conf *Config) error {
	if str, err := stringifyConfig(conf); err != nil {
		return err
	} else {
		return util.AtomicWrite(a.ConfigPath, str)
	}
}

func (c *Config) ProvisionedState() ProvisionedState {
	if c.ApplicationId == "" {
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

	if deviceType, err := util.ScanDeviceTypeSlug(util.OSRELEASE_PATH); err != nil {
		return err
	} else {
		c.DeviceType = deviceType
	}

	return nil
}

// Get /config from the Resin API specified at c.ApiEndpoint
func (c Config) getConfigFromApi() (map[string]interface{}, error) {
	return resin.GetConfig(c.ApiEndpoint)
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

func (c *Config) SetDomain(domain string) {
	// TODO: Deduplicate from defaults.go.
	if domain != "" {
		c.ApiEndpoint = fmt.Sprintf("https://api.%s", domain)
		c.VpnEndpoint = fmt.Sprintf("vpn.%s", domain)
		c.RegistryEndpoint = fmt.Sprintf("registry.%s", domain)
		c.DeltaEndpoint = fmt.Sprintf("https://delta.%s", domain)
	}
}
