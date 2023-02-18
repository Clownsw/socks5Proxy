package main

import (
	"github.com/armon/go-socks5"
	"github.com/go-yaml/yaml"
	"os"
)

type User struct {
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
}

type ServerConfig struct {
	ListenType string  `yaml:"listenType"`
	ListenInfo string  `yaml:"listenInfo"`
	UserInfo   []*User `yaml:"userInfo"`
}

func (serverConfig *ServerConfig) UserInfoToStaticCredentials() socks5.StaticCredentials {
	var result = make(socks5.StaticCredentials)

	for _, userInfo := range serverConfig.UserInfo {
		result[userInfo.UserName] = userInfo.PassWord
	}

	return result
}

func main() {
	content, err := os.ReadFile("application.yml")
	if err != nil {
		panic(err)
	}

	serverConfig := ServerConfig{}
	err = yaml.Unmarshal(content, &serverConfig)
	if err != nil {
		panic(err)
	}

	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{socks5.UserPassAuthenticator{Credentials: serverConfig.UserInfoToStaticCredentials()}},
	}

	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	err = server.ListenAndServe(serverConfig.ListenType, serverConfig.ListenInfo)
	if err != nil {
		panic(err)
	}
}
