package httpx

func IsHttpStatusOk(status int) bool {
	return status >= 200 && status < 300
}
