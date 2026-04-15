package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var lhCloudyCommands = discord.SlashCommandCreate{
	Name:        "lhcloudy",
	Description: "LhCloudy Commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "birthday",
			Description: "Get LhCloudy's birthday",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "from",
			Description: "Where LhCloudy is from?",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "youtube",
			Description: "Link to LhCloudy's YouTube channel",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "twitter",
			Description: "Link to LhCloudy's X (Twitter) account",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "tips",
			Description: "Get a link to LhCloudy's Rein tips",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "code",
			Description: "Rein Workshop Code",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "lhfurry",
			Description: "LhCloudy's furry side",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "instagram",
			Description: "LhCloudy's Instagram",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "age",
			Description: "LhCloudy's age",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "interview",
			Description: "LhCloudy's interview",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "socials",
			Description: "LhCloudy's Socials",
		},
	},
}

func (c *commands) onBirthday(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	now := time.Now().UTC()
	currentYear := now.Year()

	nextBirthday := time.Date(currentYear, time.May, 21, 0, 0, 0, 0, time.UTC)
	if now.After(nextBirthday) {
		nextBirthday = time.Date(currentYear+1, time.May, 21, 0, 0, 0, 0, time.UTC)
	}

	ts := nextBirthday.Unix()

	if now.Month() == time.May && now.Day() == 21 {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Today is Cloudy's birthday! 🎉\n\nhttps://tenor.com/bmYbD.gif",
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Cloudy's next birthday is <t:%d:D> (<t:%d:R>)", ts, ts),
	})
}

func (c *commands) onFrom(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "kotka of south eastern finland of the continent of europe",
	})
}

func (c *commands) onYoutube(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "SMÄSH THAT LIKE AND SUBSCRIBE BUTTON -> https://www.youtube.com/channel/UC2CV-HWvIrMO4mUnYtNS-7A",
	})
}

func (c *commands) onTwitter(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "https://x.com/lhcloudy",
	})
}

func (c *commands) onTips(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "W+M1",
	})
}

func (c *commands) onCode(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "rein: XEEAE | other: https://workshop.codes/u/Seita%232315",
	})
}

func (c *commands) onLhfurry(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "https://i.gyazo.com/3ae8376713000ab829a2853d0f31e6f2.png",
	})
}

func (c *commands) onInstagram(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "https://www.instagram.com/lhcloudy/",
	})
}

func (c *commands) onAge(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	birthday := time.Date(1999, time.May, 21, 0, 0, 0, 0, time.Local)

	today := time.Now()
	ageYears := int(today.Sub(birthday).Hours() / (24 * 365.25))

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("%d", ageYears),
	})
}

func (c *commands) onInterview(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "https://www.youtube.com/watch?v=sM_PkcoFgM8&t=1s",
	})
}

func (c *commands) onLinks(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	links := "" +
		"• Twitch: <https://www.twitch.tv/lhcloudy27>\n" +
		"• YouTube: <https://www.youtube.com/channel/UC2CV-HWvIrMO4mUnYtNS-7A>\n" +
		"• Discord: <https://discord.gg/jd6CZSj8jb>\n" +
		"• Twitter: <https://twitter.com/LhCloudy>\n" +
		"• Instagram: <https://www.instagram.com/lhcloudy/>\n" +
		"• Reddit: <https://www.reddit.com/r/overwatchSRpeakCHECK/>"

	embed := discord.Embed{
		Title:       "LhCloudy Links",
		Description: links,
		Color:       embedColor,
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}
