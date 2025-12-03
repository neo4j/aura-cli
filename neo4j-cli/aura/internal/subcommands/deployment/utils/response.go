package utils

func IsSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}
