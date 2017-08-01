package resin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	pinejs "github.com/resin-io/pinejs-client-go"
	"github.com/resin-os/resin-provisioner/defaults"
	"github.com/resin-os/resin-provisioner/util"
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

// Simple http POST helper with an auth token
func postWithToken(url, token, bodyType string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", bodyType)
	resp, err := http.DefaultClient.Do(req)
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
func authPost(endpoint, path string, b map[string]string) (token string, err error) {
	body, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	if resp, status, err := postUrl(endpoint+path, "application/json", body); err != nil {
		return "", err
	} else if !isHttpSuccess(status) {
		return "", nil
	} else {
		return string(resp), nil
	}
}

func authPostWithToken(endpoint, path, challengeToken string, b map[string]string) (token string, err error) {
	body, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	if resp, status, err := postWithToken(endpoint+path, challengeToken, "application/json", body); err != nil {
		return "", err
	} else if !isHttpSuccess(status) {
		return "", nil
	} else {
		return string(resp), nil
	}
}

func Login(endpoint, email, password string) (token string, err error) {
	b := map[string]string{"username": email, "password": password}
	return authPost(endpoint, "/login_", b)
}

func Signup(endpoint, email, password string) (token string, err error) {
	b := map[string]string{"email": email, "password": password}
	return authPost(endpoint, "/user/register", b)
}

func TwoFactorChallenge(endpoint, challengeToken, code string) (token string, err error) {
	b := map[string]string{"code": code}
	return authPostWithToken(endpoint, "/auth/totp/verify", challengeToken, b)
}

func GetApps(endpoint, token string) (apps []map[string]interface{}, err error) {
	var deviceType string
	client := pinejs.NewClientWithToken(endpoint+"/v1", token)
	apps = []map[string]interface{}{map[string]interface{}{"pinejs": "application"}}
	if deviceType, err = util.ScanDeviceTypeSlug(defaults.OSRELEASE_PATH); err != nil {
		return nil, fmt.Errorf("Could not get device type: %s", err)
	}
	deviceTypeFilter := fmt.Sprintf("device_type eq '%s'", deviceType)
	err = client.List(&apps, pinejs.NewQueryOptions(pinejs.Filter, deviceTypeFilter)...)
	return
}

func CreateApp(endpoint, name, token string) (id string, err error) {
	client := pinejs.NewClientWithToken(endpoint+"/v1", token)
	app := make(map[string]interface{})
	app["pinejs"] = "application"
	app["app_name"] = name
	t, e := util.ScanDeviceTypeSlug(defaults.OSRELEASE_PATH)
	if e != nil {
		return "", fmt.Errorf("Could not get device type: %s", e)
	}
	app["device_type"] = t
	if err := client.Create(&app); err != nil {
		return "", fmt.Errorf("Could not create application: %s", err)
	}
	appId, ok := app["id"].(float64)
	if !ok {
		return "", errors.New("Invalid app id from API")
	}
	return strconv.Itoa(int(appId)), nil
}

func GetApiKey(endpoint, appId, token string) (apiKey string, err error) {
	resp, status, err := postWithToken(endpoint+"/application/"+appId+"/generate-api-key", token, "application/json", []byte("{}"))
	if err != nil {
		return "", err
	} else if !isHttpSuccess(status) {
		return "", fmt.Errorf("Error getting apikey: %d %s", status, resp)
	} else {
		return strings.Trim(string(resp), `"`), nil
	}
}

func CreateOrGetDevice(endpoint string, device *map[string]interface{}, apikey string) error {
	client := pinejs.NewClient(endpoint+"/v1", apikey)
	(*device)["pinejs"] = "device"
	if err := client.Create(device); err != nil {
		if strings.Contains(err.Error(), `"uuid" must be unique`) || strings.Contains(err.Error(), `Data is referenced by uuid`) {
			uuid, ok := (*device)["uuid"].(string)
			if !ok {
				return errors.New("Invalid uuid")
			}
			devices := []map[string]interface{}{map[string]interface{}{"pinejs": "device"}}
			if err := client.List(&devices, pinejs.NewQueryOptions(pinejs.Filter, "uuid eq '"+uuid+"'")...); err != nil {
				return err
			} else if len(devices) != 1 {
				return errors.New("Invalid object returned from API")
			} else {
				*device = devices[0]
			}
		} else {
			return err
		}
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
