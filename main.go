package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/luyuhuang/subsocks/auth"
	"github.com/pelletier/go-toml"
)

func main() {
	var configPath string
	var showVersion bool
	flag.StringVar(&configPath, "c", "", "configuration file, default to 'config.toml'")
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.Parse()

	if showVersion {
		fmt.Println("Subsocks", Version)
		return
	}

	if configPath == "" {
		configPath = "config.toml"
		log.Printf("Using default configuration 'config.toml'")
	}

	config, err := toml.LoadFile(configPath)
	if err != nil {
		log.Fatalf("Load configuration failed: %s", err)
	}

	var jwtInfo *auth.JWTInfo

	if authConfig, ok := config.Get("auth").(*toml.Tree); ok {
		if tokenUrl, ok := authConfig.Get("token_url").(string); ok && tokenUrl != "" {
			if clientId, ok := authConfig.Get("client_id").(string); ok && clientId != "" {
				if clientSecret, ok := authConfig.Get("client_secret").(string); ok && clientSecret != "" {
					if projectId, ok := authConfig.Get("project_id").(string); ok && projectId != "" {
						log.Println("Auth configuration is valid, keep loading JWTInfo...")

						jwtInfo, err = auth.GetJWTInfo(tokenUrl, clientId, clientSecret, projectId)
						if err != nil {
							log.Fatalf("Get JWT info from %s failed: %s", tokenUrl, err)
						}

						if err = jwtInfo.GetAccessTokenClaims(); err != nil {
							log.Fatalf("load access token claims failed: %s", err)
						}
						jwtInfo.AccessTokenClaims.PrintJWTClaims()

						if err = jwtInfo.GetIDTokenClaims(); err != nil {
							log.Fatalf("load ID token claims failed: %s", err)
						}
						jwtInfo.IDTokenClaims.PrintJWTClaims()
					} else {
						log.Fatalf("Missing 'project_id' in '[auth] section'")
					}
				} else {
					log.Fatalf("Missing 'client_secret' in '[auth] section'")
				}
			} else {
				log.Fatalf("Missing 'client_id' in '[auth] section'")
			}
		} else {
			log.Fatalf("Missing 'token_url' in '[auth] section'")
		}
	}

	if jwtInfo == nil {
		log.Fatalf("jwtInfo is nil'")
	}
	log.Printf("JWTInfo.AccessTokenClaims: %v\n", jwtInfo.AccessTokenClaims.GetAccessInfos())
	log.Printf("JWTInfo.IDTokenClaims: %v\n", jwtInfo.IDTokenClaims.GetAccessInfos())

	if c, ok := config.Get("client").(*toml.Tree); ok {
		launchClient(c, jwtInfo)
	} else if s, ok := config.Get("server").(*toml.Tree); ok {
		launchServer(s, jwtInfo)
	} else {
		log.Fatalf("No valid configuration '[client]' or '[server]'")
	}
}
