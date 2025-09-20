package lhbot

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexraskin/LhBotGo/internal/ver"
	"github.com/alexraskin/LhBotGo/lhbot/database"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Bot struct {
	cfg        Config
	Version    ver.Version
	Discord    bot.Client
	Mongo      database.MongoClient
	HTTPClient *http.Client
	Ctx        context.Context
	DBName     string
	Collection string
}

func New(cfg Config, version ver.Version, discord bot.Client, mongo database.MongoClient, httpClient *http.Client, ctx context.Context) *Bot {

	s := &Bot{
		cfg:        cfg,
		Version:    version,
		Discord:    discord,
		Mongo:      mongo,
		HTTPClient: httpClient,
		Ctx:        ctx,
		DBName:     "lhbot",
		Collection: "lhbot_collection",
	}

	return s
}

func (b *Bot) Start(commands []discord.ApplicationCommandCreate) error {
	if b.cfg.Bot.SyncCommands {
		go func() {
			slog.Info("Syncing commands")
			if err := handler.SyncCommands(b.Discord, commands, b.cfg.Bot.GuildIDs); err != nil {
				slog.Error("failed to sync commands", "error", err)
			}
		}()
	}

	if err := b.Discord.OpenGateway(b.Ctx); err != nil {
		slog.Error("failed to open gateway", "error", err)
		return err
	}
	return nil
}

func (b *Bot) Stop() error {
	b.Mongo.Disconnect(b.Ctx)
	b.Discord.Close(b.Ctx)
	return nil
}
