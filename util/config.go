package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Bind     string            `json:"bind"`      //绑定地址
	UseCert  bool              `json:"use_cert"`  //是否使用证书
	CertFile string            `json:"cert_file"` //证书文件
	KeyFile  string            `json:"key_file"`  //证书密钥文件
	Auth     map[string]string `json:"auth"`      //鉴权 {path:token}
	Nsq      struct {
		Topic   string `json:"topic"`
		Channel string `json:"channel"`
		Address string `json:"address"`
	} `json:"nsq"`
}

var Cfg Config

func LoadConfig() error {
	//get executable filename without extension
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	name := filepath.Base(ex)
	ext := filepath.Ext(name)
	nameWithoutExt := name[0 : len(name)-len(ext)]
	//get config file path
	cfgPath := nameWithoutExt + ".config.json"
	if len(os.Args) < 2 {
		log.Println("No config file specified, using default config: ", cfgPath)
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
			Auth: map[string]string{
				"": "token",
			},
		}, "", "  ")
		if err != nil {
			log.Println("Marshal default config failed:", err)
			return err
		}
		err = os.WriteFile(cfgPath, cfgbuf, 0644)
		if err != nil {
			log.Println("Write default config failed:", err)
			return err
		}
	}
	err = json.Unmarshal(cfgbuf, &Cfg) //解析配置文件 反序列化json到结构体
	//check config
	//remove '/' in auth path
	changed := false
	for k, v := range Cfg.Auth {
		if len(k) == 0 {
			continue
		}
		if k[0] == '/' {
			Cfg.Auth[k[1:]] = v
			delete(Cfg.Auth, k)
			changed = true
		}
	}
	if changed {
		cfgbuf, err = json.MarshalIndent(Cfg, "", "  ")
		os.WriteFile(cfgPath, cfgbuf, 0644)
		log.Println("Config file changed.")
	}
	if err != nil {
		log.Println("Unmarshal config failed:", err)
		return err
	}
	return nil
}
