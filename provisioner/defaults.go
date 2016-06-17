package provisioner

const (
	SERVICES_ROOT_PATH      = "/lib/systemd/system/"
	SUPERVISOR_PATH         = SERVICES_ROOT_PATH + "resin-supervisor.service"
	UPDATE_RESIN_TIMER_PATH = SERVICES_ROOT_PATH + "update-resin-supervisor.timer"
	UPDATE_RESIN_PATH       = SERVICES_ROOT_PATH + "update-resin-supervisor.service"
	PREPARE_OPENVPN_PATH    = SERVICES_ROOT_PATH + "prepare-openvpn.service"
	OSRELEASE_PATH          = "/etc/os-release"
	SUPERVISOR_CONF_PATH    = "/etc/supervisor.conf"

	DEFAULT_RESIN_DOMAIN        = "resinstaging.io"
	INIT_UPDATER_SUPERVISOR_TAG = "production"

	PUBNUB_SUBSCRIBE_KEY_ENV_VAR = "RESIN_PUBNUB_SUBSCRIBE_KEY"
	PUBNUB_PUBLISH_KEY_ENV_VAR   = "RESIN_PUBNUB_PUBLISH_KEY"
	MIXPANEL_TOKEN_ENV_VAR       = "RESIN_MIXPANEL_TOKEN"
	DOMAIN_OVERRIDE_ENV_VAR      = "RESIN_DOMAIN_OVERRIDE"
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
