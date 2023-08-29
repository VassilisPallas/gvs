package config

type Configuration struct {
	VERSIONS_URL              string
	DOWNLOAD_VERSION_BASE_URL string
}

func GetConfig() Configuration {
	return Configuration{
		VERSIONS_URL:              "https://go.dev/dl/?mode=json&include=all",
		DOWNLOAD_VERSION_BASE_URL: "https://go.dev/dl",
	}
}
