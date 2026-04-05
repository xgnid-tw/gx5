package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NotionToken         string
	NotionUserDBID      string
	NotionOrderDBID     string
	DiscordToken        string
	DiscordAppID        string
	DiscordGuildID      string
	DiscordLogChannelID string
	ExchangeRateJPYTWD  float64
	TagRoleMap          map[string]string
}

func Load() (Config, error) {
	cfg := Config{
		NotionToken:         os.Getenv("NOTION_TOKEN"),
		NotionUserDBID:      os.Getenv("NOTION_USER_DB_ID"),
		NotionOrderDBID:     os.Getenv("NOTION_ORDER_DB_ID"),
		DiscordToken:        os.Getenv("DISCORD_TOKEN"),
		DiscordAppID:        os.Getenv("DISCORD_APP_ID"),
		DiscordGuildID:      os.Getenv("DISCORD_GUILD_ID"),
		DiscordLogChannelID: os.Getenv("DISCORD_GUILD_LOG_CHANNEL_ID"),
	}
	cfg.TagRoleMap = parseTagRoleMap(os.Getenv("TAG_ROLE_MAP"))

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

	if cfg.DiscordToken == "" {
		return Config{}, fmt.Errorf("DISCORD_TOKEN is required")
	}

	if cfg.DiscordAppID == "" {
		return Config{}, fmt.Errorf("DISCORD_APP_ID is required")
	}

	if cfg.DiscordGuildID == "" {
		return Config{}, fmt.Errorf("DISCORD_GUILD_ID is required")
	}

	if cfg.DiscordLogChannelID == "" {
		return Config{}, fmt.Errorf("DISCORD_GUILD_LOG_CHANNEL_ID is required")
	}

	if cfg.NotionOrderDBID == "" {
		return Config{}, fmt.Errorf("NOTION_ORDER_DB_ID is required")
	}

	return cfg, nil
}

const tagRoleMapParts = 2

func parseTagRoleMap(raw string) map[string]string {
	m := make(map[string]string)
	if raw == "" {
		return m
	}

	for pair := range strings.SplitSeq(raw, ",") {
		parts := strings.SplitN(pair, "=", tagRoleMapParts)
		if len(parts) == tagRoleMapParts {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return m
}
