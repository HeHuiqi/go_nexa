package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var configPath string = "accounts.json"

type MainConfig struct {
	Mnenmonic        string `json:"mnenmonic"`
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	Address          string `json:"address"`
	ChangePrivateKey string `json:"changePrivateKey"`
	ChangePublicKey  string `json:"changePublicKey"`
	ChangeAddress    string `json:"changeAddress"`
}

// type Configs map[string]json.RawMessage

type Configs []MainConfig

var conf *MainConfig
var confs Configs

var instanceOnce sync.Once

// 从配置文件中载入json字符串
func LoadConfig(path string) (Configs, *MainConfig) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config conf failed: ", err)
	}
	allConfigs := make(Configs, 4)
	err = json.Unmarshal(buf, &allConfigs)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	mainConfig := &allConfigs[0]

	return allConfigs, mainConfig
}

// 初始化 可以运行多次
func SetConfig(path string) {
	allConfigs, mainConfig := LoadConfig(path)
	configPath = path
	conf = mainConfig
	confs = allConfigs
}

// 初始化，只能运行一次
func Init(path string) *MainConfig {
	if conf != nil && path != configPath {
		log.Printf("the config is already initialized, oldPath=%s, path=%s", configPath, path)
	}
	instanceOnce.Do(func() {
		allConfigs, mainConfig := LoadConfig(path)
		configPath = path
		conf = mainConfig
		confs = allConfigs
	})

	return conf
}

// 初始化配置文件 为 struct 格式
func Instance() *MainConfig {
	if conf == nil {
		Init(configPath)
	}
	return conf
}

// 初始化配置文件 为 []格式
func AllConfig() Configs {
	if conf == nil {
		Init(configPath)
	}
	return confs
}

// 获取配置文件路径
func ConfigPath() string {
	return configPath
}

func ConfigInit() {
	path := ConfigPath()
	pwd, _ := os.Getwd()
	path = pwd + "/" + path
	Init(path)

}
