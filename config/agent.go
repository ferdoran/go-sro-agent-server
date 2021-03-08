package config

import (
	"github.com/ferdoran/go-sro-framework/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	AgentHost                 = "agent.host"
	AgentPort                 = "agent.port"
	AgentSecret               = "agent.secret"
	AgentModuleId             = "agent.module_id"
	AgentDataPath             = "agent.data_path"
	AgentPrelinkedNavdataFile = "agent.prelinked_navdata_file"

	GatewayModuleId = "gateway.module_id"
	GatewaySecret   = "gateway.secret"

	GameTimeTicksPerSecond = "game.time.ticks_per_second"
	GameTimeDaySpeed       = "game.time.day_speed"
)

func Initialize() {
	config.Initialize()

	setDefaultValues()
	bindEnvAliases()

	logrus.Info("agent config initialized")
}

func bindEnvAliases() {
	viper.BindEnv(AgentHost, "AGENT_HOST")
	viper.BindEnv(AgentPort, "AGENT_PORT")
	viper.BindEnv(AgentSecret, "AGENT_SECRET")
	viper.BindEnv(AgentModuleId, "AGENT_MODULE_ID")
	viper.BindEnv(AgentDataPath, "AGENT_DATA_PATH")
	viper.BindEnv(AgentPrelinkedNavdataFile, "AGENT_PRELINKED_NAVDATA_FILE")

	viper.BindEnv(GatewaySecret, "GATEWAY_SECRET")
	viper.BindEnv(GatewayModuleId, "GATEWAY_MODULE_ID")

	viper.BindEnv(GameTimeTicksPerSecond, "GAME_TIME_TICKS_PER_SECOND")
	viper.BindEnv(GameTimeDaySpeed, "GAME_TIME_DAY_SPEED")
}

func setDefaultValues() {
	viper.SetDefault(AgentHost, "127.0.0.1")
	viper.SetDefault(AgentPort, 15882)
	viper.SetDefault(AgentSecret, "agent-server")
	viper.SetDefault(AgentModuleId, "AgentServer")
	viper.SetDefault(AgentDataPath, "./Data")
	viper.SetDefault(AgentPrelinkedNavdataFile, "./prelinked_navdata.gob")

	viper.SetDefault(GatewaySecret, "gateway-server")
	viper.SetDefault(GatewayModuleId, "GatewayServer")

	viper.SetDefault(GameTimeTicksPerSecond, 10)
	viper.SetDefault(GameTimeDaySpeed, 1.0)
}
