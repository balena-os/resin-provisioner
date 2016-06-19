package provisioner

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type ProvisionedState int

type Api struct {
	ConfigPath string
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

		if err := a.RegisterDevice(conf); err != nil {
			return err
		}

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

func (a Api) RegisterDevice(c *Config) error {
	newUuid := false
	registeredAt := time.Now().Unix()
	if c.Uuid == "" {
		newUuid = true
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

	body, err := json.Marshal(device)
	if err != nil {
		return err
	}
	url := c.ApiEndpoint + "/v1/device?apikey=" + c.ApiKey
	resp, status, err := postUrl(url, "application/json", body)
	if err != nil {
		return err
	} else if !isHttpSuccess(status) {
		// If device already exists and we're not generating the uuid
		if strings.Contains(string(resp), `"uuid" must be unique`) && !newUuid {
			url = c.ApiEndpoint + `/v1/device?$filter=uuid eq '` + c.Uuid + `'&apikey=` + c.ApiKey
			if resp, status, err = getUrl(url); err != nil {
				return err
			} else if !isHttpSuccess(status) {
				return fmt.Errorf("Error getting device from API: %d %s", status, resp)
			}
		} else {
			return fmt.Errorf("Error when registering: %d %s", status, resp)
		}
	}
	if err = json.Unmarshal(resp, &device); err != nil {
		return err
	}
	if deviceId, ok := device["id"].(float64); !ok {
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
