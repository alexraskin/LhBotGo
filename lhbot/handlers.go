package lhbot

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

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

		content := strings.TrimSpace(e.Message.Content)

		if strings.HasPrefix(content, "!shatter") {

			if len(e.Message.Mentions) == 0 {
				b.Discord.Rest().CreateMessage(e.Message.ChannelID, discord.MessageCreate{
					Content: "You need to mention someone to shatter them! `!shatter @username`",
				})
				return
			}

			target := e.Message.Mentions[0]

			var lhCloudBlockMessages = []string{
				"Blocked.. immune to your shatter!",
				"LhCloudy is immune to your shatter!",
				"Blocked - MTD",
				"ez block... L + ratio",
				"sr peak check?",
			}

			rand.New(rand.NewSource(time.Now().UnixNano()))

			if target.ID.String() == "127122091139923968" {
				b.Discord.Rest().CreateMessage(e.Message.ChannelID, discord.MessageCreate{
					Content: lhCloudBlockMessages[rand.Intn(len(lhCloudBlockMessages))],
				})
				return
			}

			options := []string{"backline deleted", "interrupted mid-shatter", "spawn", "landed...", "hit", "was blocked by", "miss"}
			choice := options[rand.Intn(len(options))]

			var message string
			switch choice {
			case "backline deleted":
				message = fmt.Sprintf("You shattered %s’s entire backline. 💀", target.Mention())
			case "interrupted mid-shatter":
				message = fmt.Sprintf("%s slept you mid-shatter. 💤", target.Mention())
			case "spawn":
				message = fmt.Sprintf("You shattered %s, but your team was still in spawn. 🕳️", target.Mention())
			case "landed...":
				message = fmt.Sprintf("You landed the shatter on %s, but no one followed up. 🤨", target.Mention())
			case "hit":
				message = fmt.Sprintf("Your shatter hit %s! 💥🔨", target.Mention())
			case "was blocked by":
				message = fmt.Sprintf("Your shatter was blocked by %s, the enemy mercy typed MTD. 🧱", target.Mention())
			case "miss":
				message = "You shattered no one, so it missed. Your team is now flaming you, and the enemy mercy typed MTD."
			}

			b.Discord.Rest().CreateMessage(e.Message.ChannelID, discord.MessageCreate{
				Content: message,
			})
			return
		}

		// hardcoded because i'm lazy and these are the only channels that should be used for commands
		if (e.Message.ChannelID == 935059381802905631 || e.Message.ChannelID == 932412270565269545) && strings.HasPrefix(e.Message.Content, "!") {
			b.Discord.Rest().CreateMessage(e.Message.ChannelID, discord.MessageCreate{
				Content: "`!shatter` is supported, but all other commands should use `/` instead",
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
