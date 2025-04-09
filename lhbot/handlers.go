package lhbot

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

func (b *Bot) StartTasks(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	b.statusTask(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down status task")
			return
		case <-ticker.C:
			b.statusTask(ctx)
		}
	}
}

func (b *Bot) statusTask(ctx context.Context) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	statuses := []string{
		"Reinhardt", "Overwatch 2", "LhCloudy27",
	}

	b.Discord.SetPresence(ctx,
		gateway.WithPlayingActivity(statuses[rand.Intn(len(statuses))]),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
		gateway.WithAfk(false),
	)
}
