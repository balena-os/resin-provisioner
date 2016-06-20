// +build integration

package integrationtest

import (
	"fmt"
	"os"
	"testing"

	"github.com/resin-os/resin-provisioner/provisioner"
)

func TestRegisterDevice(t *testing.T) {
	k := os.Getenv("API_KEY")
	u := os.Getenv("USER_ID")
	a := os.Getenv("APP_ID")
	if k == "" || u == "" || a == "" {
		t.Skip("Skipping integration test, env vars not defined")
	} else {
		var c = provisioner.Config{DeviceType: "intel-edison", ApplicationId: a, ApiKey: k, UserId: u, ApiEndpoint: "https://api.resinstaging.io"}
		api := provisioner.New("./config.json")
		if err := api.RegisterDevice(&c); err != nil {
			t.Error(err)
		} else if c.RegisteredAt == 0 {
			t.Error("RegisteredAt not written to config")
		} else if c.DeviceId == 0 {
			t.Error("DeviceId not written to config")
		}
		fmt.Printf("%+v\n", c)

		// Test that it doesn't fail it's an already registered device
		var c2 = provisioner.Config{DeviceType: "intel-edison", ApplicationId: a, ApiKey: k, UserId: u, ApiEndpoint: "https://api.resinstaging.io", Uuid: c.Uuid}
		if err := api.RegisterDevice(&c2); err != nil {
			t.Error(err)
		} else if c.DeviceId != c2.DeviceId {
			t.Error("Device ids don't match when using the same uuid")
		} else if c2.RegisteredAt == 0 {
			t.Error("RegisteredAt not written to config when using same uuid")
		}
		fmt.Printf("%+v\n", c2)
	}
}
