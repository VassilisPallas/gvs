package config

type Configuration struct {
	GO_BASE_URL string
}

func GetConfig() Configuration {
	return Configuration{
		GO_BASE_URL: "https://go.dev/dl",
	}
}
