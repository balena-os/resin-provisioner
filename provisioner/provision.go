package provisioner

import (
	"fmt"

	"github.com/resin-os/resin-provisioner/defaults"
	"github.com/resin-os/resin-provisioner/util"
)

func Provision(appId, apiKey, userId, userName, configPath, domain string, dryRun bool) error {
	if deviceType, err := util.ScanDeviceTypeSlug(defaults.OSRELEASE_PATH); err != nil {
		return err
	} else if config, err := readConfig(configPath); err != nil {
		return err
	} else {
		config.AppId = appId
		config.ApiKey = apiKey
		config.UserId = userId
		config.UserName = userName
		config.DeviceType = deviceType
		config.ApiEndpoint = fmt.Sprintf("https://api.%s", domain)
		config.RegistryEndpoint = fmt.Sprintf("registry.%s", domain)
		config.VpnEndpoint = fmt.Sprintf("vpn.%s", domain)

		if marshalledConfig, err := marshal(config); err != nil {
			return err
		} else if dryRun {
			fmt.Printf("Supervisor config: %v\n", marshalledConfig)
			return nil
		} else if err := util.AtomicWrite(configPath, marshalledConfig); err != nil {
			return err
		} else if conn, err := newDbus(); err != nil {
			return err
		} else {
			defer conn.Close()
			if err := conn.restartUnit(defaults.SUPERVISOR_PATH); err != nil {
				return err
			} else if err := conn.restartUnit(defaults.PREPARE_VPN_PATH); err != nil {
				return err
			} else if err := conn.restartUnit(defaults.VPN_RESIN_PATH); err != nil {
				return err
			}
		}
	}

	return nil
}
