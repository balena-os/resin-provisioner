package provisioner

const (
	SUPERVISOR_PATH         = "/lib/systemd/system/resin-supervisor.service"
	UPDATE_RESIN_TIMER_PATH = "/lib/systemd/system/update-resin-supervisor.timer"

	DEFAULT_RESIN_DOMAIN = "resinstaging.io"
)

var DefaultConfig = Config{
	// Default to standard ethernet connection.

	Files: map[string]string{
		"network/settings":       "[global]\nOfflineMode=false\n\n[WiFi]\nEnable=true\nTethering=false\n\n[Wired]\nEnable=true\nTethering=false\n\n[Bluetooth]\nEnable=true\nTethering=false",
		"network/network.config": "[service_home_ethernet]\nType = ethernet\nNameservers = 8.8.8.8,8.8.4.4",
	},
	AppUpdatePollInterval: "60000",
	ListenPort:            "48484",
	VpnPort:               "443",
	ApiEndpoint:           "https://api." + DEFAULT_RESIN_DOMAIN,
	VpnEndpoint:           "vpn." + DEFAULT_RESIN_DOMAIN,
	RegistryEndpoint:      "registry." + DEFAULT_RESIN_DOMAIN,
	DeltaEndpoint:         "https://delta." + DEFAULT_RESIN_DOMAIN,
}
