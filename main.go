package main

import (
	"flag"
	"fmt"

	"github.com/luyuhuang/subsocks/auth"
	"github.com/luyuhuang/subsocks/client"
	acc "github.com/luyuhuang/subsocks/control/access"
	"github.com/luyuhuang/subsocks/control/rule"
	"github.com/luyuhuang/subsocks/log"
	"github.com/luyuhuang/subsocks/server"
)

var (
	accessFlag bool
	username   string
	password   string

	Prefix string

	TOKEN_URL   = "https://zitadel-go-p36jcg.us1.zitadel.cloud/oauth/v2/token"
	PROJECT_ID  = "347030683681641570"
	PROTOCOL    = "http"
	HTTP_PATH   = "/proxy"
	LISTEN_PORT = "1030"
	ACCESS_PORT = "1080"
)

func main() {
	flag.BoolVar(&accessFlag, "acc", false, "whether to start access mode")
	flag.StringVar(&username, "u", "", "username in zitadel")
	flag.StringVar(&password, "p", "", "password in zitadel")
	flag.Parse()

	if accessFlag {
		Prefix = "[Access Mode]"
		log.Info(Prefix, "subsocks starting...")
	} else {
		Prefix = "[Client Mode]"
		log.Info(Prefix, "subsocks starting...")
	}

	if username == "" {
		log.Error(Prefix, "subsocks param [username] is empty")
		return
	}

	if password == "" {
		log.Error(Prefix, "param [password] is empty")
		return
	}

	jwtInfo, err := auth.GetJWTInfo(TOKEN_URL, username, password, PROJECT_ID)
	if err != nil {
		log.Error(Prefix, "Get JWT failed, token_url:", TOKEN_URL, "error:", err)
		return
	}

	if err = jwtInfo.GetIDTokenClaims(); err != nil {
		log.Error(Prefix, "load ID token claims failed, error:", err)
		return
	}

	if accessFlag {
		cli := client.NewClient(fmt.Sprintf("127.0.0.1:%d", LISTEN_PORT))
		cli.Config.ServerProtocol = PROTOCOL
		cli.Config.HTTPPath = HTTP_PATH

		accessTree := jwtInfo.AccessTree()
		accessTree.FillInfo()
		ruleInfos := accessTree.ListRule()

		m := make(map[string]*rule.Info)
		for _, ruleInfo := range ruleInfos {
			ruleInfoCopy := ruleInfo
			m[ruleInfo.ServiceInfo.Host] = &ruleInfoCopy
		}

		r, err := client.NewRulesFromStructMap(m)
		if err != nil {
			log.Error(Prefix, "Load rules failed, error:", err)
			return
		}

		cli.Rules = r
		cli.JWTInfo = jwtInfo
		cli.AccessTree = &accessTree
		if err = cli.Serve(); err != nil {
			log.Error(Prefix, "Launch client failed, error:", err)
		}
	} else {
		metadata := jwtInfo.IDTokenClaims.GetMetadata()
		info, err := acc.GetInfoById(metadata.Id)
		if err != nil {
			log.Error(Prefix, "Get access info failed, error", err)
			return
		}

		log.Info(Prefix, "Get access info success, info:", info)

		ser := server.NewServer(PROTOCOL, fmt.Sprintf("0.0.0.0:%d", ACCESS_PORT))
		ser.Config.HTTPPath = HTTP_PATH

		if err := ser.Serve(); err != nil {
			log.Error(Prefix, "Launch server failed, error:", err)
		}
	}
}
