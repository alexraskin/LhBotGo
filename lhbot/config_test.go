package lhbot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func clearEnvVars(t *testing.T) {
	t.Helper()
	for _, key := range []string{"BOT_TOKEN", "BOT_GUILD_IDS", "BOT_SYNC_COMMANDS", "MONGO_URI"} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.Bot.Token != "" {
		t.Errorf("expected empty token, got %q", cfg.Bot.Token)
	}
	if cfg.Bot.GuildIDs != nil {
		t.Errorf("expected nil guild IDs, got %v", cfg.Bot.GuildIDs)
	}
	if !cfg.Bot.SyncCommands {
		t.Error("expected sync_commands to default to true")
	}
	if cfg.Mongo.URI != "mongodb://localhost:27017" {
		t.Errorf("expected default mongo URI, got %q", cfg.Mongo.URI)
	}
}

func TestLoadConfigFromToml(t *testing.T) {
	clearEnvVars(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	err := os.WriteFile(cfgPath, []byte(`
[bot]
token = "test-token"
guild_ids = [123456789]
sync_commands = false

[mongo]
uri = "mongodb://testhost:27017"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Bot.Token != "test-token" {
		t.Errorf("expected token %q, got %q", "test-token", cfg.Bot.Token)
	}
	if len(cfg.Bot.GuildIDs) != 1 || cfg.Bot.GuildIDs[0] != 123456789 {
		t.Errorf("expected guild ID 123456789, got %v", cfg.Bot.GuildIDs)
	}
	if cfg.Bot.SyncCommands {
		t.Error("expected sync_commands to be false")
	}
	if cfg.Mongo.URI != "mongodb://testhost:27017" {
		t.Errorf("expected mongo URI %q, got %q", "mongodb://testhost:27017", cfg.Mongo.URI)
	}
}

func TestLoadConfigMissingFileUsesDefaults(t *testing.T) {
	clearEnvVars(t)

	cfg, err := LoadConfig("/nonexistent/config.toml")
	if err != nil {
		t.Fatal(err)
	}

	expected := defaultConfig()
	if cfg.Bot.Token != expected.Bot.Token {
		t.Errorf("expected default token, got %q", cfg.Bot.Token)
	}
	if cfg.Mongo.URI != expected.Mongo.URI {
		t.Errorf("expected default mongo URI, got %q", cfg.Mongo.URI)
	}
}

func TestLoadConfigInvalidToml(t *testing.T) {
	clearEnvVars(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	err := os.WriteFile(cfgPath, []byte(`this is not valid toml [[[`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(cfgPath)
	if err == nil {
		t.Error("expected error for invalid TOML, got nil")
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	tests := []struct {
		name   string
		envs   map[string]string
		verify func(t *testing.T, cfg Config)
	}{
		{
			name: "BOT_TOKEN overrides",
			envs: map[string]string{"BOT_TOKEN": "env-token"},
			verify: func(t *testing.T, cfg Config) {
				if cfg.Bot.Token != "env-token" {
					t.Errorf("expected %q, got %q", "env-token", cfg.Bot.Token)
				}
			},
		},
		{
			name: "MONGO_URI overrides",
			envs: map[string]string{"MONGO_URI": "mongodb://envhost:27017"},
			verify: func(t *testing.T, cfg Config) {
				if cfg.Mongo.URI != "mongodb://envhost:27017" {
					t.Errorf("expected %q, got %q", "mongodb://envhost:27017", cfg.Mongo.URI)
				}
			},
		},
		{
			name: "BOT_SYNC_COMMANDS true",
			envs: map[string]string{"BOT_SYNC_COMMANDS": "true"},
			verify: func(t *testing.T, cfg Config) {
				if !cfg.Bot.SyncCommands {
					t.Error("expected sync_commands to be true")
				}
			},
		},
		{
			name: "BOT_SYNC_COMMANDS 1",
			envs: map[string]string{"BOT_SYNC_COMMANDS": "1"},
			verify: func(t *testing.T, cfg Config) {
				if !cfg.Bot.SyncCommands {
					t.Error("expected sync_commands to be true for value '1'")
				}
			},
		},
		{
			name: "BOT_SYNC_COMMANDS false",
			envs: map[string]string{"BOT_SYNC_COMMANDS": "false"},
			verify: func(t *testing.T, cfg Config) {
				if cfg.Bot.SyncCommands {
					t.Error("expected sync_commands to be false")
				}
			},
		},
		{
			name: "BOT_GUILD_IDS single",
			envs: map[string]string{"BOT_GUILD_IDS": "111222333"},
			verify: func(t *testing.T, cfg Config) {
				if len(cfg.Bot.GuildIDs) != 1 || cfg.Bot.GuildIDs[0] != snowflake.ID(111222333) {
					t.Errorf("expected [111222333], got %v", cfg.Bot.GuildIDs)
				}
			},
		},
		{
			name: "BOT_GUILD_IDS multiple comma-separated",
			envs: map[string]string{"BOT_GUILD_IDS": "111, 222, 333"},
			verify: func(t *testing.T, cfg Config) {
				if len(cfg.Bot.GuildIDs) != 3 {
					t.Errorf("expected 3 guild IDs, got %d", len(cfg.Bot.GuildIDs))
				}
			},
		},
		{
			name: "BOT_GUILD_IDS skips invalid entries",
			envs: map[string]string{"BOT_GUILD_IDS": "111,notanumber,333"},
			verify: func(t *testing.T, cfg Config) {
				if len(cfg.Bot.GuildIDs) != 2 {
					t.Errorf("expected 2 valid guild IDs, got %d: %v", len(cfg.Bot.GuildIDs), cfg.Bot.GuildIDs)
				}
			},
		},
		{
			name: "empty env vars do not override",
			envs: map[string]string{},
			verify: func(t *testing.T, cfg Config) {
				expected := defaultConfig()
				if cfg.Bot.Token != expected.Bot.Token {
					t.Errorf("expected default token, got %q", cfg.Bot.Token)
				}
				if cfg.Mongo.URI != expected.Mongo.URI {
					t.Errorf("expected default mongo URI, got %q", cfg.Mongo.URI)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars(t)
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			cfg := defaultConfig()
			applyEnvOverrides(&cfg)
			tt.verify(t, cfg)
		})
	}
}

func TestEnvOverridesTomlValues(t *testing.T) {
	clearEnvVars(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	err := os.WriteFile(cfgPath, []byte(`
[bot]
token = "toml-token"

[mongo]
uri = "mongodb://tomlhost:27017"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("BOT_TOKEN", "env-token")
	t.Setenv("MONGO_URI", "mongodb://envhost:27017")

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Bot.Token != "env-token" {
		t.Errorf("expected env token to override toml, got %q", cfg.Bot.Token)
	}
	if cfg.Mongo.URI != "mongodb://envhost:27017" {
		t.Errorf("expected env mongo URI to override toml, got %q", cfg.Mongo.URI)
	}
}
