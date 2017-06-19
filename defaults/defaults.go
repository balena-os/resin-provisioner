package defaults

const (
	SERVICES_ROOT_PATH = "/lib/systemd/system/"
	SUPERVISOR_PATH    = SERVICES_ROOT_PATH + "resin-supervisor.service"
	PREPARE_VPN_PATH   = SERVICES_ROOT_PATH + "prepare-openvpn.service"
	VPN_RESIN_PATH     = SERVICES_ROOT_PATH + "openvpn-resin.service"

	OSRELEASE_PATH = "/etc/os-release"
	RESIN_DOMAIN   = "resin.io"
	CONFIG_PATH    = "/mnt/boot/config.json"
)
