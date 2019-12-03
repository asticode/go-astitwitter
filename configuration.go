package astitwitter

import (
	"flag"

	astihttp "github.com/asticode/go-astitools/http"
)

// Flags
var (
	APIKey       = flag.String("twitter-api-key", "", "the api key")
	APISecretKey = flag.String("twitter-api-secret-key", "", "the api secret key")
)

// Configuration represents the lib's configuration
type Configuration struct {
	APIKey       string `toml:"api_key"`
	APISecretKey string `toml:"api_secret_key"`
	Sender       astihttp.SenderOptions
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		APIKey:       *APIKey,
		APISecretKey: *APISecretKey,
	}
}
