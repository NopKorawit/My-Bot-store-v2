package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Sheet GoogleSheetConfig
	Line  LineConfig
}

type GoogleSheetConfig struct {
	GoogleCredentialsPath string `envconfig:"GOOGLE_CREDENTIALS_PATH" required:"true"`
	SpreadSheetId         string `envconfig:"SPREAD_SHEET_ID" required:"true"`
}

type LineConfig struct {
	LineChannelSecret string `envconfig:"LINE_CHANNEL_SECRET" required:"true"`
	LineChannelToken  string `envconfig:"LINE_CHANNEL_TOKEN" required:"true"`
}

func (cfg *AppConfig) Init() {
	envconfig.MustProcess("", &cfg.Sheet)
	envconfig.MustProcess("", &cfg.Line)
}

func New() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	appCfg := AppConfig{}
	appCfg.Init()

	return &appCfg
}
