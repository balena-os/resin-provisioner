package provisioner

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/resin-os/resin-provisioner/resin"
)

type ProvisionedState int

type Api struct {
	ConfigPath string
	Domain     string
	listener   net.Listener
	server     *http.Server
}

type ProvisionOpts struct {
	UserId        string `json:"userId"`
	ApplicationId string `json:"applicationId"`
	ApiKey        string `json:"apikey"`
}

const (
	Unknown ProvisionedState = iota
	Unprovisioned
	Provisioning
	Provisioned
)

func (s ProvisionedState) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Unprovisioned:
		return "unprovisioned"
	case Provisioning:
		return "provisioning"
	case Provisioned:
		return "provisioned"
	}

	return "invalid"
}

func New(configPath string) *Api {
	ret := &Api{ConfigPath: configPath}
	ret.initServer()

	return ret
}

// Returns the device provisioned state.
func (a *Api) State() (ProvisionedState, error) {
	if conf, err := a.readConfig(); err != nil {
		return Unknown, fmt.Errorf("Cannot read config: %s", err)
	} else {
		return conf.ProvisionedState(), nil
	}
}

func (a *Api) StateJson() (ret string, err error) {
	var state ProvisionedState

	if state, err = a.State(); err == nil {
		ret = fmt.Sprintf(`{"state": "%s"}`, state)
	}

	return
}

func (a *Api) Provision(opts *ProvisionOpts) error {
	if state, err := a.State(); err != nil {
		return err
	} else if state != Unprovisioned {
		return fmt.Errorf("Cannot provision, device is %s.", state)
	}

	if !isInteger(opts.UserId) || !isInteger(opts.ApplicationId) ||
		!isValidApiKey(opts.ApiKey) {
		return fmt.Errorf("Invalid options.")
	}

	if conf, err := a.readConfig(); err != nil {
		return fmt.Errorf("Cannot read config: %s", err)
	} else {
		// First check to see whether config.json has changed from
		// underneath us.
		if state := conf.ProvisionedState(); state != Unprovisioned {
			return fmt.Errorf("Cannot provision, device is %s.", state)
		}

		// Ok, now we go for it.
		conf.UserId = opts.UserId
		conf.ApplicationId = opts.ApplicationId
		conf.ApiKey = opts.ApiKey
		if err := conf.GetKeysFromApi(); err != nil {
			return err
		}

		// We save the config.json as it is now to persist the uuid
		if err := a.writeConfig(conf); err != nil {
			return err
		}

		// Register the device on resin
		if err := a.RegisterDevice(conf); err != nil {
			return err
		}

		// We write config again with registered_at and deviceId
		if err := a.writeConfig(conf); err != nil {
			return err
		}

		// Next we need to enable the supervisor systemd service.
		if conn, err := NewDbus(); err != nil {
			return err
		} else {
			defer conn.Close()

			return conn.SupervisorEnableStart()
		}
	}
}

// TODO: Use proper pinejs client for all this.
func (a Api) RegisterDevice(c *Config) error {
	registeredAt := time.Now().Unix()
	if c.Uuid == "" {
		if uuid, err := randomHexString(UUID_BYTE_LENGTH); err != nil {
			return err
		} else {
			c.Uuid = uuid
		}
	}

	device := make(map[string]interface{})
	device["user"] = c.UserId
	device["application"] = c.ApplicationId
	device["uuid"] = c.Uuid
	device["device_type"] = c.DeviceType
	device["registered_at"] = registeredAt

	err := resin.CreateOrGetDevice(c.ApiEndpoint, &device, c.ApiKey)
	if err != nil {
		return err
	}
	if deviceId, ok := device["id"].(float64); !ok || deviceId == 0 {
		return errors.New("Device returned from API is invalid")
	} else {
		c.RegisteredAt = int64(registeredAt)
		c.DeviceId = int64(deviceId)
		return nil
	}
}

func (a *Api) ConfigJson() (string, error) {
	if conf, err := a.readConfig(); err != nil {
		return "", fmt.Errorf("Cannot read config: %s", err)
	} else if str, err := stringifyConfig(conf); err != nil {
		return "", fmt.Errorf("Cannot stringfy config: %s", err)
	} else {
		return str, nil
	}
}

func (a *Api) DeviceUrl() (string, error) {
	if conf, err := a.readConfig(); err != nil {
		return "", fmt.Errorf("Cannot read config: %s", err)
	} else if conf.ApplicationId == "" {
		return "", fmt.Errorf("Empty application ID.")
	} else if conf.DeviceId == 0 {
		return "", fmt.Errorf("Empty device ID.")
	} else {
		return fmt.Sprintf("https://dashboard.%s/apps/%s/devices/%d/summary",
			a.Domain, conf.ApplicationId, conf.DeviceId), nil
	}
}
