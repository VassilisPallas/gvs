package config

type Configuration struct {
	GO_BASE_URL     string
	REQUEST_TIMEOUT int
}

func GetConfig() Configuration {
	return Configuration{
		GO_BASE_URL:     "https://go.dev/dl",
		REQUEST_TIMEOUT: 30, // 30 seconds for both fetching versions and downloading version tar
	}
}
