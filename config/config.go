package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	NotionToken         string
	NotionUserDBID      string
	NotionOthersDBID    string
	NotionOrderDBID     string
	DiscordToken        string
	DiscordAppID        string
	DiscordLogChannelID string
	ExchangeRateJPYTWD  float64
	WorkerCrontab       string
	Debug               bool
}

func Load() (Config, error) {
	cfg := Config{
		NotionToken:         os.Getenv("NOTION_TOKEN"),
		NotionUserDBID:      os.Getenv("NOTION_USER_DB_ID"),
		NotionOthersDBID:    os.Getenv("NOTION_OTHERS_DB_ID"),
		NotionOrderDBID:     os.Getenv("NOTION_ORDER_DB_ID"),
		DiscordToken:        os.Getenv("DISCORD_TOKEN"),
		DiscordAppID:        os.Getenv("DISCORD_APP_ID"),
		DiscordLogChannelID: os.Getenv("DISCORD_GUILD_LOG_CHANNEL_ID"),
		WorkerCrontab:       os.Getenv("WORKER_CORNTAB"),
		Debug:               os.Getenv("DEBUG") == "1",
	}

	rate, err := strconv.ParseFloat(os.Getenv("EXCHANGE_RATE_JPY_TWD"), 64)
	if err != nil || rate <= 0 {
		return Config{}, fmt.Errorf("EXCHANGE_RATE_JPY_TWD must be a positive number")
	}

	cfg.ExchangeRateJPYTWD = rate

	if cfg.NotionToken == "" {
		return Config{}, fmt.Errorf("NOTION_TOKEN is required")
	}

	if cfg.NotionUserDBID == "" {
		return Config{}, fmt.Errorf("NOTION_USER_DB_ID is required")
	}

	if cfg.NotionOthersDBID == "" {
		return Config{}, fmt.Errorf("NOTION_OTHERS_DB_ID is required")
	}

	if cfg.DiscordToken == "" {
		return Config{}, fmt.Errorf("DISCORD_TOKEN is required")
	}

	if cfg.DiscordAppID == "" {
		return Config{}, fmt.Errorf("DISCORD_APP_ID is required")
	}

	if cfg.DiscordLogChannelID == "" {
		return Config{}, fmt.Errorf("DISCORD_GUILD_LOG_CHANNEL_ID is required")
	}

	if cfg.WorkerCrontab == "" {
		return Config{}, fmt.Errorf("WORKER_CORNTAB is required")
	}

	if cfg.NotionOrderDBID == "" {
		return Config{}, fmt.Errorf("NOTION_ORDER_DB_ID is required")
	}

	if cfg.DiscordAppID == "" {
		return Config{}, fmt.Errorf("DISCORD_APP_ID is required")
	}

	return cfg, nil
}
