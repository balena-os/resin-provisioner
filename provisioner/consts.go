package provisioner

const (
	SUPERVISOR_PATH         = "/lib/systemd/system/resin-supervisor.service"
	UPDATE_RESIN_TIMER_PATH = "/lib/systemd/system/update-resin-supervisor.timer"

	DEFAULT_RESIN_DOMAIN = "resinstaging.io"
)

var DefaultConfig = Config{
	AppUpdatePollInterval: "60000",
	ListenPort:            "48484",
	VpnPort:               "443",
	ApiEndpoint:           "https://api." + DEFAULT_RESIN_DOMAIN,
	VpnEndpoint:           "vpn." + DEFAULT_RESIN_DOMAIN,
	RegistryEndpoint:      "registry." + DEFAULT_RESIN_DOMAIN,
	DeltaEndpoint:         "https://delta." + DEFAULT_RESIN_DOMAIN,
}
