package config

import (
	"fmt"
	"os"
)

type Config struct {
	NotionToken         string
	NotionUserDBID      string
	NotionOthersDBID    string
	DiscordToken        string
	DiscordLogChannelID string
	WorkerCrontab       string
	Debug               bool
}

func Load() (Config, error) {
	cfg := Config{
		NotionToken:         os.Getenv("NOTION_TOKEN"),
		NotionUserDBID:      os.Getenv("NOTION_USER_DB_ID"),
		NotionOthersDBID:    os.Getenv("NOTION_OTHERS_DB_ID"),
		DiscordToken:        os.Getenv("DISCORD_TOKEN"),
		DiscordLogChannelID: os.Getenv("DISCORD_GUILD_LOG_CHANNEL_ID"),
		WorkerCrontab:       os.Getenv("WORKER_CORNTAB"),
		Debug:               os.Getenv("DEBUG") == "1",
	}

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

	if cfg.DiscordLogChannelID == "" {
		return Config{}, fmt.Errorf("DISCORD_GUILD_LOG_CHANNEL_ID is required")
	}

	if cfg.WorkerCrontab == "" {
		return Config{}, fmt.Errorf("WORKER_CORNTAB is required")
	}

	return cfg, nil
}
