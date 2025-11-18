package config

type Auth struct {
	TokenUrl     string `toml:"token_url"`
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}
