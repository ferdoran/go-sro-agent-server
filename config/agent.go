package config

import (
	"encoding/json"
	"github.com/ferdoran/go-sro-framework/config"
	"github.com/sirupsen/logrus"
	"os"
)

type AgentConfig struct {
	config.Config
	GameTimeConfig GameTimeConfig `json:"game_time_config"`
}

func LoadConfig(configFile string) {
	logrus.Printf("loading config: %s\n", configFile)
	file, err := os.Open(configFile)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	cfg := AgentConfig{}
	err = decoder.Decode(&cfg)

	if err != nil {
		panic(err.Error())
	}

	GlobalConfig = cfg
	config.GlobalConfig = cfg.Config
}

var GlobalConfig AgentConfig
