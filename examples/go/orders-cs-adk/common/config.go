package common

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAI struct {
		APIKey  string `mapstructure:"api_key"`
		BaseURL string `mapstructure:"base_url"`
		Model   string `mapstructure:"model"`
	} `mapstructure:"openai"`
	Services struct {
		SearchAPIURL string `mapstructure:"search_api_url"`
		PolicyAPIURL string `mapstructure:"policy_api_url"`
		MCPBaseURL   string `mapstructure:"mcp_base_url"`
		MCPCommand   string `mapstructure:"mcp_command"`
	} `mapstructure:"services"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AddConfigPath(".")
	_ = v.ReadInConfig()
	var cfg Config
	_ = v.Unmarshal(&cfg)
	v.SetConfigName("config.local")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AddConfigPath(".")
	if err := v.MergeInConfig(); err == nil {
		_ = v.Unmarshal(&cfg)
	}
	// environment variables intentionally ignored to enforce file-driven config
	cfg.OpenAI.APIKey = strings.TrimSpace(cfg.OpenAI.APIKey)
	cfg.OpenAI.BaseURL = strings.TrimSpace(cfg.OpenAI.BaseURL)
	cfg.OpenAI.BaseURL = strings.Trim(cfg.OpenAI.BaseURL, "`\"' ")
	cfg.OpenAI.Model = strings.TrimSpace(cfg.OpenAI.Model)
	return &cfg, nil
}
