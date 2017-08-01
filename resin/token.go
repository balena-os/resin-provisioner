package resin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func GetTwoFactorRequired(token string) (bool, error) {
	if jsonToken, err := parse(token); err != nil {
		return false, err
	} else if required, ok := jsonToken["twoFactorRequired"].(bool); !ok {
		return false, nil
	} else {
		return required, nil
	}
}

func GetUserId(token string) (string, error) {
	if jsonToken, err := parse(token); err != nil {
		return "", err
	} else if id, ok := jsonToken["id"].(float64); !ok {
		return "", errors.New("Invalid id in token")
	} else {
		return strconv.Itoa(int(id)), nil
	}
}

func GetUserName(token string) (string, error) {
	if jsonToken, err := parse(token); err != nil {
		return "", err
	} else if username, ok := jsonToken["username"].(string); !ok {
		return "", errors.New("Invalid username in token")
	} else {
		return username, nil
	}
}

func parse(token string) (map[string]interface{}, error) {
	var jsonToken map[string]interface{}
	if content := strings.Split(token, "."); len(content) != 3 {
		return nil, errors.New("Invalid token")
	} else if parsedContent, err := base64.RawURLEncoding.DecodeString(content[1]); err != nil {
		return nil, err
	} else if err = json.Unmarshal(parsedContent, &jsonToken); err != nil {
		return nil, err
	} else {
		return jsonToken, nil
	}
}
