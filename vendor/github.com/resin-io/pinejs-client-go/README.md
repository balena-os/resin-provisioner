# Pine.js Go Library

This is a simple Go library for interacting with [pine.js][pine].
Also includes object definitions to interact with the [resin.io][resin] API.

## Usage

```go
import (
	pinejs "github.com/resin-io/pinejs-client-go"
	"github.com/resin-io/pinejs-client-go/resin"
)

// Using a struct
func GetDatDevice() {
	var device resin.Device{Id: 1234}
	pineClient := pinejs.NewClient("https://api.resinstaging.io/ewa", "secretapikey")

	if err := pineClient.Get(&device); err != nil {
		log.Fatalln(err)
	} else {
		// device contains the device object with id 1234.
	}
}

// Using a map
func GetDatDeviceOnAMap() {
	device := map[string]interface{}{"pinejs": "device", "id": 1234}
	pineClient := pinejs.NewClient("https://api.resinstaging.io/ewa", "secretapikey")

	if err := pineClient.Get(&device); err != nil {
		log.Fatalln(err)
	} else {
		// device contains the device object with id 1234.
	}
}
```

More examples in [test/test.go](./test/test.go)

[pine]:https://bitbucket.org/rulemotion/pinejs/overview
[resin]:https://resin.io
