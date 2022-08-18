package config

import (
	"os"

	"github.com/joho/godotenv"
	// "github.com/kelseyhightower/envconfig"
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

// func (cfg *AppConfig) Init() {
// 	envconfig.MustProcess("", &cfg.Sheet)
// 	envconfig.MustProcess("", &cfg.Line)
// }

func New() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	appCfg := AppConfig{}
	appCfg.Line.LineChannelSecret = os.Getenv("CHANNEL_SECRET")
	appCfg.Line.LineChannelToken = os.Getenv("CHANNEL_TOKEN")
	appCfg.Sheet.GoogleCredentialsPath = os.Getenv("GOOGLE_CREDENTIALS_PATH")
	appCfg.Sheet.SpreadSheetId = os.Getenv("SPREAD_SHEET_ID")

	return &appCfg
}
