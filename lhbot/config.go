package lhbot

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/disgoorg/snowflake/v2"
)

func LoadConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	cfg := defaultConfig()
	if _, err = toml.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
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
	Token        string         `toml:"token"`
	GuildIDs     []snowflake.ID `toml:"guild_ids"`
	SyncCommands bool           `toml:"sync_commands"`
}

func (c BotConfig) String() string {
	return fmt.Sprintf("\n Token: %s\n GuildIDs: %s\n SyncCommands: %t",
		c.Token,
		c.GuildIDs,
		c.SyncCommands,
	)
}
