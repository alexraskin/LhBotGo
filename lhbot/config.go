package lhbot

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

func LoadConfig(cfgPath string) (Config, error) {
	_ = godotenv.Load()

	cfg := defaultConfig()

	file, err := os.Open(cfgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("failed to open config file: %w", err)
		}
	} else {
		defer func() {
			_ = file.Close()
		}()
		if _, err = toml.NewDecoder(file).Decode(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	applyEnvOverrides(&cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("BOT_TOKEN"); v != "" {
		cfg.Bot.Token = v
	}
	if v := os.Getenv("BOT_GUILD_IDS"); v != "" {
		var ids []snowflake.ID
		for _, raw := range strings.Split(v, ",") {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				continue
			}
			id, err := snowflake.Parse(raw)
			if err == nil {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			cfg.Bot.GuildIDs = ids
		}
	}
	if v := os.Getenv("BOT_COMMAND_CHANNEL_IDS"); v != "" {
		var ids []snowflake.ID
		for _, raw := range strings.Split(v, ",") {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				continue
			}
			id, err := snowflake.Parse(raw)
			if err == nil {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			cfg.Bot.CommandChannelIDs = ids
		}
	}
	if v := os.Getenv("BOT_SYNC_COMMANDS"); v != "" {
		cfg.Bot.SyncCommands = v == "true" || v == "1"
	}
	if v := os.Getenv("BOT_LHCLOUDY_ID"); v != "" {
		id, err := snowflake.Parse(v)
		if err == nil {
			cfg.Bot.LhCloudyID = id
		}
	}
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.Mongo.URI = v
	}
}

func defaultConfig() Config {
	return Config{
		Bot: BotConfig{
			Token:        "",
			GuildIDs:     nil,
			SyncCommands: true,
		},
		Mongo: MongoConfig{
			URI: "mongodb://localhost:27017",
		},
	}
}

type Config struct {
	Bot   BotConfig   `toml:"bot"`
	Mongo MongoConfig `toml:"mongo"`
}

func (c Config) String() string {
	return fmt.Sprintf("Bot: %s\nMongo: %s",
		c.Bot,
		c.Mongo,
	)
}

type MongoConfig struct {
	URI string `toml:"uri"`
}

func (c MongoConfig) String() string {
	return fmt.Sprintf("\n URI: %s", c.URI)
}

type BotConfig struct {
	Token             string         `toml:"token"`
	GuildIDs          []snowflake.ID `toml:"guild_ids"`
	SyncCommands      bool           `toml:"sync_commands"`
	CommandChannelIDs []snowflake.ID `toml:"command_channel_ids"`
	LhCloudyID        snowflake.ID   `toml:"lhcloudy_id"`
}

func (c BotConfig) String() string {
	return fmt.Sprintf("\n Token: %s\n GuildIDs: %s\n SyncCommands: %t",
		c.Token,
		c.GuildIDs,
		c.SyncCommands,
	)
}
