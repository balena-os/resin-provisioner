package provisioner

import "testing"

var minimalJson = `
{
   "deviceType":"raspberrypi3",
   "pubnubSubscribeKey":"sub-c-abc12345-ab1a-12a1-9876-01ab1abcd8zx",
   "pubnubPublishKey":"pub-c-abc12345-ab1a-12a1-9876-01ab1abcd8zx",
   "mixpanelToken":"abcdefghijklmnopqrstuvwxyz123456",
   "ListenPort":"1234"
}
`

// Ensure the default values are set + not overwritten.
func TestParseMinimalConfigJson(t *testing.T) {
	if conf, err := parseConfig(minimalJson); err != nil {
		t.Fatalf("Parse failed: ERROR: %s", err)
	} else {
		if conf.ListenPort == DefaultConfig.ListenPort {
			t.Error("ListenPort overwriten by default!")
		}

		// TODO: Automate this via reflection.

		if conf.AppUpdatePollInterval != DefaultConfig.AppUpdatePollInterval {
			t.Error("MISMATCH: AppUpdatePollInterval")
		}
		if conf.VpnPort != DefaultConfig.VpnPort {
			t.Error("MISMATCH: VpnPort")
		}
		if conf.ApiEndpoint != DefaultConfig.ApiEndpoint {
			t.Error("MISMATCH: ApiEndpoint")
		}
		if conf.VpnEndpoint != DefaultConfig.VpnEndpoint {
			t.Error("MISMATCH: VpnEndpoint")
		}
		if conf.RegistryEndpoint != DefaultConfig.RegistryEndpoint {
			t.Error("MISMATCH: RegistryEndpoint")
		}
		if conf.DeltaEndpoint != DefaultConfig.DeltaEndpoint {
			t.Error("MISMATCH: DeltaEndpoint")
		}
	}
}
