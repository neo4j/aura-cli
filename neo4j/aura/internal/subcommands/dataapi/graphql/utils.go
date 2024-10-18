package graphql

import (
	"encoding/base64"
)

const (
	AuthProviderTypeApiKey = "api-key"
	AuthProviderTypeJwks   = "jwks"
)

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
