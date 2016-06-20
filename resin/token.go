package resin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func GetUserId(token string) (string, error) {
	var jsonToken map[string]interface{}
	if content := strings.Split(token, "."); len(content) != 3 {
		return "", errors.New("Invalid token")
	} else if parsedContent, e := base64.RawURLEncoding.DecodeString(content[1]); e != nil {
		return "", e
	} else if e = json.Unmarshal(parsedContent, &jsonToken); e != nil {
		return "", e
	} else if id, ok := jsonToken["id"].(float64); !ok {
		return "", errors.New("Invalid id in token")
	} else {
		return strconv.Itoa(int(id)), nil
	}
}
