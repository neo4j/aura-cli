package utils

import "net/http"

func GetSuccessfulAuthenticationHandler(counter *int) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		*counter++
		res.WriteHeader(200)
		res.Write([]byte(`{"access_token":"<token>","expires_in":3600,"token_type":"bearer"}`))
	}
}
