package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		Telegram    TelegramConfig
	}
	TelegramConfig struct {
		ApiToken string
		Debug    bool `mapstructure:"debug"`
	}
	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}
)

func Init(configDir string) (*Config, error) {
	if err := parseConfigFile(configDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("telegram", &cfg.Telegram); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Telegram.ApiToken = os.Getenv("API_TOKEN")

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")
}

func parseConfigFile(configDir, env string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}
