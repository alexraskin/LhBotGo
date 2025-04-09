package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var overwatchCommands = discord.SlashCommandCreate{
	Name:        "ow",
	Description: "Overwatch Commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "shatter",
			Description: "Shatter another user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to shatter",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "reinquote",
			Description: "Get a random Reinhardt quote",
		},
	},
}

func (c *commands) onShatter(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var userID snowflake.ID
	var err error

	for _, option := range data.Options {
		if option.Name == "user" {
			var idStr string
			if err := json.Unmarshal(option.Value, &idStr); err != nil {
				slog.Error("Error unmarshalling user ID", "error", err)
				return e.CreateMessage(discord.MessageCreate{
					Content: "An error occurred - KEKL",
				})
			}
			userID, err = snowflake.Parse(idStr)
			if err != nil {
				slog.Error("Error parsing user ID", "error", err)
				return e.CreateMessage(discord.MessageCreate{
					Content: "An error occurred - KEKL",
				})
			}
			break
		}
	}

	targetUser, ok := data.Resolved.Users[userID]
	if !ok {
		slog.Error("Could not resolve user from Discord API", "userID", userID)
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
		})
	}

	var lhCloudBlockMessages = []string{
		"Blocked.. immune to your shatter!",
		"LhCloudy is immune to your shatter!",
		"Blocked - MTD",
		"ez block... L + ratio",
		"sr peak check?",
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	if targetUser.ID.String() == "127122091139923968" {
		return e.CreateMessage(discord.MessageCreate{
			Content: lhCloudBlockMessages[rand.Intn(len(lhCloudBlockMessages))],
		})
	}

	options := []string{"hit", "miss", "was blocked by"}
	choice := options[rand.Intn(len(options))]

	var message string
	switch choice {
	case "hit":
		message = fmt.Sprintf("Your shatter hit %s! 💥🔨", targetUser.Mention())
	case "was blocked by":
		message = fmt.Sprintf("Your shatter was blocked by %s, the enemy mercy typed MTD. 🧱", targetUser.Mention())
	case "miss":
		message = "You shattered no one, so it missed. Your team is now flaming you, and the enemy mercy typed MTD."
	}

	return e.CreateMessage(discord.MessageCreate{Content: message})
}

var reinQuotes = []string{
	"Reinhardt, at your service.",
	"I'm getting cooked alive in here! Ughhhh!",
	"Who's ready for some karaoke?",
	"Ah, my friends! What's with all this standing around? There's a glory to be won! We shall prove ourselves in glorious combat",
	"Steel yourselves! Push. Them. Back! Just a little longer! Bear down! Hold fast! They will not breach our defense!",
	"That's how it's done! Ah, impressive, if I do say so myself! To see justice done is its own reward!",
	"There's no glory in a stopped payload! The payload has come to a halt. Let's get it moving!",
	"You kids today with your techno music. You should listen to the classics, like Hasselhoff!",
	"We are out of time! Attack! Make every second count! Crush their defenses!",
	"Taking the objective! Join me in my glory! I am taking the objective! Try and stop me!",
	"Haha! Finally, I have beaten the champion at arm wrestling!",
	"We fought a terrible battle here. Many crusaders lost their lives. This is the hometown of my master, Balderich. He was born here...and he died here. Too much blood was spilled in my country during the war...",
	"So this is the base of the great D.va! The city's walls must be splattered with posters of her!",
	"Are you chicken?",
	"Unstoppable!",
	"Easy does it.",
	"I AM ON FIRE! Come here and get burned",
	"Barrier active!",
	"Don't worry, my friends. I am your shield!",
	"Ha! Get behind me!",
	"Come out and face me!",
	"Are you ready?! Here I come!",
	"I feel powerful!",
	"I am unstoppable!",
	"Hammer down!",
	"Earthshatter, ready! My ultimate is ready! Join me!",
	"Bring. It. On! I live for this!",
	"Hahaha! Is. That. All?",
	"I salute you!",
	"Catchphrase!",
	"Honor and glory.",
	"Let me show you how it's done.",
	"Are you afraid to fight me?",
	"Ooo, glittering prizes!",
	"Respect your elders.",
	"I'm the ultimate crushing machine.",
	"What do we have here?",
	"Forgive and forget, like snow from yesteryear.",
	"Bring me another!",
	"You're on my naughty list.",
	"Crusader online.",
	"Fortune favors the bold.",
	"Precision German engineering.",
	"This old dog still knows a few tricks.",
	"100% German power!",
	"You shame yourself!",
	"Smashing!",
}

func (c *commands) onQuote(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := reinQuotes[rand.Intn(len(reinQuotes))]
	return e.CreateMessage(discord.MessageCreate{
		Content: quote,
	})
}
