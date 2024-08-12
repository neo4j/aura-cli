package testutils

import (
	"bytes"
	"encoding/json"
	"strings"
)

func FormatJson(unformatted string) (string, error) {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(strings.TrimSpace(unformatted)), "", "\t")
	if err != nil {
		return "", err
	}
	return pretty.String() + "\n", nil
}
