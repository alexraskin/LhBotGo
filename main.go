package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"

	"github.com/alexraskin/LhBotGo/lhbot"
	"github.com/alexraskin/LhBotGo/lhbot/commands"
	"github.com/alexraskin/LhBotGo/lhbot/database"
)

func main() {
	cfgPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	cfg, err := lhbot.LoadConfig(*cfgPath)
	if err != nil {
		slog.Error("Error while loading config", slog.Any("err", err))
		return
	}

	version := "unknown"
	goVersion := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
		goVersion = info.GoVersion
	}

	slog.Info("Starting LhBot...", slog.String("version", version), slog.String("go_version", goVersion))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	discord, err := disgo.New(cfg.Bot.Token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentMessageContent)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
	)
	if err != nil {
		slog.Error("Error while creating bot client", slog.Any("err", err))
		return
	}
	defer discord.Close(timeoutCtx)

	mongoClient, err := database.New(timeoutCtx, cfg.Mongo.URI)
	if err != nil {
		slog.Error("Error while creating mongo client", slog.Any("err", err))
		return
	}
	defer mongoClient.Disconnect(timeoutCtx)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	b := lhbot.New(cfg, version, goVersion, discord, mongoClient, httpClient, ctx)

	b.Discord.AddEventListeners(commands.New(b), lhbot.MessageHandler(b), lhbot.OnReady(b))

	if err := b.Start(commands.Commands); err != nil {
		slog.Error("Failed to start bot", "error", err)
		return
	}

	slog.Info("LhBot started")
	<-ctx.Done()

	if err := b.Stop(); err != nil {
		slog.Error("Failed to stop bot", "error", err)
		return
	}

	slog.Info("LhBot stopped")
}
