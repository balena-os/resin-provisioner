package provisioner

const (
	SERVICES_ROOT_PATH      = "/lib/systemd/system/"
	MOUNT_OVERLAY_PATH      = SERVICES_ROOT_PATH + "etc-systemd-system-resin.target.wants.mount"
	SUPERVISOR_PATH         = SERVICES_ROOT_PATH + "resin-supervisor.service"
	UPDATE_RESIN_TIMER_PATH = SERVICES_ROOT_PATH + "update-resin-supervisor.timer"
	UPDATE_RESIN_PATH       = SERVICES_ROOT_PATH + "update-resin-supervisor.service"
	OPENVPN_PATH            = SERVICES_ROOT_PATH + "openvpn-resin.service"
	SUPERVISOR_CONF_PATH    = "/etc/resin-supervisor/supervisor.conf"
	RESIN_SERVICES_PATH     = "/etc/resin-connectable.conf"

	DEFAULT_RESIN_DOMAIN        = "resin.io"
	INIT_UPDATER_SUPERVISOR_TAG = "production"

	UUID_BYTE_LENGTH = 31
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
