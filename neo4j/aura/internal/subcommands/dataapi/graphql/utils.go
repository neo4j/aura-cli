package graphql

import (
	"encoding/base64"
)

const (
	SecurityAuthProviderTypeApiKey = "api-key"
	SecurityAuthProviderTypeJwks   = "jwks"
)

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
