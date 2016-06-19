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
		}
		fmt.Printf("%+v", c)
	}
}
