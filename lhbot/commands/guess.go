package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"strings"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/alexraskin/LhBotGo/lhbot/database"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var guessCommands = discord.SlashCommandCreate{
	Name:        "lh",
	Description: "LhBot Guess Commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "guess",
			Description: "Take a guess at what LH means",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "guess",
					Description: "Your guess",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "latest",
			Description: "Get the last 5 guesses",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List all the guesses",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "count",
			Description: "Shows the number of guesses",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "hint",
			Description: "Get a hint",
		},
	},
}

func (c *commands) onHint(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var hint []string = []string{
		"It is in English",
		"Made by 10 year old finnish lad",
		"Clever",
		"Masaa can be bribed",
		"It's not long hammer",
	}

	randomHint := hint[rand.Intn(len(hint))]
	return e.CreateMessage(discord.MessageCreate{
		Content: randomHint,
	})
}

func (c *commands) onLatest(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guesses, err := c.Bot.Mongo.GetLatestGuesses(c.Bot.Ctx, c.Bot.DBName, c.Bot.Collection, 5)
	if err != nil {
		slog.Error("Error getting guesses", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Error getting guesses",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	var latestGuesses string
	for _, guess := range guesses {
		latestGuesses += fmt.Sprintf("%s - %s\n", guess.LhGuess, guess.GuessedBy)
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: latestGuesses,
		Flags:   discord.MessageFlagEphemeral,
	})
}

func (c *commands) onList(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guesses, err := c.Bot.Mongo.GetGuesses(c.Bot.Ctx, c.Bot.DBName, c.Bot.Collection)
	if err != nil {
		slog.Error("Error getting guesses", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Error getting guesses LUL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"LhGuess", "GuessedBy", "GuessedAt"})

	for _, guess := range guesses {
		var guessedAt string
		if !guess.GuessedAt.IsZero() {
			guessedAt = guess.GuessedAt.Format(time.RFC3339)
		} else {
			guessedAt = ""
		}
		writer.Write([]string{guess.LhGuess, guess.GuessedBy, guessedAt})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		slog.Error("Error writing CSV", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to generate CSV. LUL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	file := discord.NewFile("guesses.csv", "text/csv", &buf)
	return e.CreateMessage(discord.MessageCreate{
		Content: "Here is the list of guesses",
		Files:   []*discord.File{file},
	})
}

func (c *commands) onCount(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	count, err := c.Bot.Mongo.CountGuesses(c.Bot.Ctx, c.Bot.DBName, c.Bot.Collection)
	if err != nil {
		slog.Error("Error getting guess count", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Error getting guess count LUL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("There have been %d guesses", count),
	})
}

func (c *commands) onGuess(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var guess string
	for _, option := range data.Options {
		if option.Name == "guess" {
			guess = string(option.Value)
		}
	}
	if !checkGuess(guess) {
		now := time.Now()
		embed := discord.Embed{
			Title:       "That is not a valid guess 🚨",
			Description: "KEKL",
			Color:       0xFF0000,
			Timestamp:   &now,
		}
		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
			Flags:  discord.MessageFlagEphemeral,
		})
	}

	guess = cleanGuess(guess)

	guessObj := database.Guess{
		LhGuess:   guess,
		GuessedBy: e.Member().Member.User.Username,
		GuessedAt: time.Now(),
	}

	if err := c.Bot.Mongo.AddGuess(c.Bot.Ctx, c.Bot.DBName, c.Bot.Collection, guessObj); err != nil {
		slog.Error("Error adding guess", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Error adding guess LUL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Guess Submitted!\n\nYour guess: `%s`\n\nGuessed by: %s", guess, e.Member().Member.User.Mention()),
	})
}

func cleanGuess(guess string) string {
	s := strings.ToLower(guess)
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.TrimSpace(s)
	return s
}

func checkGuess(guess string) bool {
	re := regexp.MustCompile(`\b[lL]\w*`)
	return re.MatchString(guess) && !goaway.IsProfane(guess)
}
