package config

type Configuration struct {
	// GO_BASE_URL contains the go.dev URL that will be used
	// both for fetching the version list and download the selected one.
	GO_BASE_URL string

	// REQUEST_TIMEOUT contain the timeout in seconds, that will be used on the HTTP client.
	REQUEST_TIMEOUT int
}

func GetConfig() Configuration {
	return Configuration{
		GO_BASE_URL:     "https://go.dev/dl",
		REQUEST_TIMEOUT: 30, // 30 seconds for both fetching versions and downloading version tar,

	}
}
