package resin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	//pinejs "github.com/resin-io/pinejs-client-go"
)

// Simple http GET request helper
func getUrl(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}

// Simple http POST helper
func postUrl(url string, bodyType string, body []byte) ([]byte, int, error) {
	resp, err := http.Post(url, bodyType, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return b, resp.StatusCode, nil
}

func isHttpSuccess(status int) bool {
	return status/100 == 2
}

func pineQueryEscape(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}
func Login(endpoint, email, password string) (token string, err error) {

	return "t", errors.New("Not implemented")
}

func Signup(endpoint, email, password string) (token string, err error) {
	fmt.Printf("Trying to sign up with email: %s and password: %s\n", email, password)
	return "t", errors.New("Not implemented")
}

func GetApps(endpoint, token string) (apps []map[string]interface{}, err error) {
	//client := ppinejs.NewClient(c., "secretapikey")
	return
}

func CreateApp(endpoint, name, token string) (id string, err error) {
	return
}

func GetApiKey(endpoint, appId, token string) (apiKey string, err error) {
	return
}

// TODO: use pinejs client
func CreateOrGetDevice(endpoint string, device *map[string]interface{}, apikey string) error {
	body, err := json.Marshal(device)
	if err != nil {
		return err
	}
	u := endpoint + "/v1/device?apikey=" + apikey
	resp, status, err := postUrl(u, "application/json", body)
	if err != nil {
		return err
	} else if !isHttpSuccess(status) {
		// If device already exists
		uuid, ok := (*device)["uuid"].(string)
		if !ok {
			return errors.New("Invalid uuid")
		}
		if strings.Contains(string(resp), `"uuid" must be unique`) || strings.Contains(string(resp), `Data is referenced by uuid`) {
			u := endpoint + `/v1/device?` + pineQueryEscape(`$filter=uuid eq '`+uuid+`'&apikey=`+apikey)
			if resp, status, err = getUrl(u); err != nil {
				return err
			} else if !isHttpSuccess(status) {
				return fmt.Errorf("Error getting device from API: %d %s", status, resp)
			} else {
				d := make(map[string]interface{})
				if err = json.Unmarshal(resp, &d); err != nil {
					return err
				}
				if arr, ok := d["d"].([]interface{}); !ok {
					return errors.New("Invalid object returned from API")
				} else if len(arr) != 1 {
					return errors.New("Invalid object returned from API")
				} else if dev, ok := arr[0].(map[string]interface{}); !ok {
					return errors.New("Invalid object returned from API")
				} else {
					*device = dev
				}
			}
		} else {
			return fmt.Errorf("Error when registering: %d %s", status, resp)
		}
	} else if err = json.Unmarshal(resp, device); err != nil {
		return err
	}
	return nil
}

func GetConfig(endpoint string) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	if r, status, err := getUrl(endpoint + "/config"); err != nil {
		return nil, err
	} else if !isHttpSuccess(status) {
		return nil, fmt.Errorf("Error status from Resin API: %d", status)
	} else if err = json.Unmarshal(r, &conf); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}
