package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Bind       string            `json:"bind"`      //绑定地址
	UseCert    bool              `json:"use_cert"`  //是否使用证书
	CertFile   string            `json:"cert_file"` //证书文件
	KeyFile    string            `json:"key_file"`  //证书密钥文件
	AuthBearer map[string]string `json:"auth"`      //鉴权 {path:token}
}

var cfg Config

func loadConfig() {
	cfgPath := "config.json"
	if len(os.Args) < 2 {
		log.Println("No config file specified, using default config.json")
	} else {
		cfgPath = os.Args[1]
		log.Println("Using config file:", cfgPath)
	}
	cfgbuf, err := os.ReadFile(cfgPath) //读取配置文件
	if err != nil {
		log.Println("Read config file failed:", err)
		//write default config
		cfgbuf, err = json.MarshalIndent(Config{
			Bind:     ":8080",
			UseCert:  false,
			CertFile: "fullchain.cer",
			KeyFile:  "private.key",
			AuthBearer: map[string]string{
				"": "token",
			},
		}, "", "  ")
		os.WriteFile(cfgPath, cfgbuf, 0644)
		if err != nil {
			log.Println("Write default config failed:", err)
			return
		}
	}
	err = json.Unmarshal(cfgbuf, &cfg) //解析配置文件 反序列化json到结构体
	if err != nil {
		log.Println("Unmarshal config failed:", err)
		return
	}
}
