package provisioner

const (
	SUPERVISOR_PATH         = "/lib/systemd/system/resin-supervisor.service"
	UPDATE_RESIN_TIMER_PATH = "/lib/systemd/system/update-resin-supervisor.timer"
	PREPARE_OPENVPN_PATH    = "/lib/systemd/system/prepare-openvpn.service"

	DEFAULT_RESIN_DOMAIN = "resinstaging.io"

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
