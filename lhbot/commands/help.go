package commands

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var helpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "Help commands",
}

func (c *commands) onHelp(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	embed := discord.Embed{
		Title: "LhBot Commands",
		Color: 0x5865F2,
		Fields: []discord.EmbedField{

			{Name: "/lh guess", Value: "submit a guess", Inline: ptr(true)},
			{Name: "/lh latest", Value: "view last 5 guesses", Inline: ptr(true)},
			{Name: "/lh list", Value: "download all guesses CSV", Inline: ptr(true)},
			{Name: "/lh count", Value: "show total guesses", Inline: ptr(true)},
			{Name: "/lh hint", Value: "get a hint", Inline: ptr(true)},

			{Name: "/lhcloudy related commands", Value: "Commands related to LhCloudy", Inline: ptr(false)},

			{Name: "/ow shatter", Value: "shatter a user", Inline: ptr(true)},
			{Name: "/ow reinquote", Value: "random Reinhardt quote", Inline: ptr(true)},

			{Name: "/fun", Value: "random cat, dog, or meme image", Inline: ptr(false)},

			{Name: "/info", Value: "show bot uptime & stats", Inline: ptr(true)},
			{Name: "/help", Value: "show this help message", Inline: ptr(true)},
		},
		Timestamp: ptrTime(e.CreatedAt()),
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}

func ptr(b bool) *bool {
	return &b
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
