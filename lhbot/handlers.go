package lhbot

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func MessageHandler(b *Bot) bot.EventListener {
	return bot.NewListenerFunc(func(e *events.MessageCreate) {
		if e.Message.Author.Bot {
			return
		}

		// hardcoded because i'm lazy and these are the only channels that should be used for commands
		if (e.Message.ChannelID == 935059381802905631 || e.Message.ChannelID == 932412270565269545) && strings.HasPrefix(e.Message.Content, "!") {
			b.Discord.Rest().CreateMessage(e.Message.ChannelID, discord.MessageCreate{
				Content: "Hey, I see you are trying to use the old command syntax. Please use `/` to start a command.",
			})
		}
	})
}

func OnReady(b *Bot) bot.EventListener {
	return bot.NewListenerFunc(func(e *events.Ready) {
		if err := b.Discord.SetPresence(b.Ctx,
			gateway.WithPlayingActivity("Overwatch 2"),
			gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
			slog.Error("Failed to set presence", slog.Any("err", err))
		}
	})
}
