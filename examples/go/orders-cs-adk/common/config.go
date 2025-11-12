package common

import (
    "os"

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
        MCPBaseURL    string `mapstructure:"mcp_base_url"`
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

    if cfg.OpenAI.APIKey == "" {
        cfg.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
    }
    if cfg.OpenAI.BaseURL == "" {
        cfg.OpenAI.BaseURL = os.Getenv("OPENAI_API_BASE")
    }
    if cfg.OpenAI.Model == "" {
        cfg.OpenAI.Model = os.Getenv("OPENAI_MODEL")
    }
    if cfg.Services.SearchAPIURL == "" {
        cfg.Services.SearchAPIURL = os.Getenv("SEARCH_API_URL")
    }
    if cfg.Services.PolicyAPIURL == "" {
        cfg.Services.PolicyAPIURL = os.Getenv("POLICY_API_URL")
    }
    if cfg.Services.MCPBaseURL == "" {
        cfg.Services.MCPBaseURL = os.Getenv("MCP_BASE_URL")
    }
    if cfg.Server.Port == "" {
        cfg.Server.Port = os.Getenv("PORT")
    }
    return &cfg, nil
}
