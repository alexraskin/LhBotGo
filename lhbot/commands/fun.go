package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var funCommands = discord.SlashCommandCreate{
	Name:        "fun",
	Description: "Fun commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "cat",
			Description: "Get a random cat image",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "dog",
			Description: "Get a random dog image",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "meme",
			Description: "Get a random meme",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "trump",
			Description: "Get a random quote from Donald Trump",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "personalize",
					Description: "Get a personalized Trump Quote",
					Required:    false,
				},
			},
		},
	},
}

type CatData struct {
	ID        string   `json:"id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	URL       string   `json:"url"`
}

func (c *commands) onCat(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	resp, err := c.HTTPClient.Get("https://cataas.com/cat?json=true")
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	var catData CatData
	if err := json.Unmarshal(body, &catData); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	embed := discord.Embed{
		Image: &discord.EmbedResource{
			URL: catData.URL,
		},
		Footer: &discord.EmbedFooter{
			Text: "Powered by cataas.com",
		},
		Color: 0x5865F2,
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}

type DogData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (c *commands) onDog(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	resp, err := c.HTTPClient.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	var dogData DogData
	if err := json.Unmarshal(body, &dogData); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	embed := discord.Embed{
		Image: &discord.EmbedResource{
			URL: dogData.Message,
		},
		Footer: &discord.EmbedFooter{
			Text: "Powered by dog.ceo",
		},
		Color: 0x5865F2,
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}

type MemeData struct {
	PostLink  string   `json:"postLink"`
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler   bool     `json:"spoiler"`
	Author    string   `json:"author"`
	Ups       int      `json:"ups"`
	Preview   []string `json:"preview"`
}

func (c *commands) onMeme(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	resp, err := c.HTTPClient.Get("https://meme-api.com/gimme")
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	var memeData MemeData
	if err := json.Unmarshal(body, &memeData); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	embed := discord.Embed{
		Image: &discord.EmbedResource{
			URL: memeData.URL,
		},
		Color: 0x5865F2,
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}

type TrumpQuote struct {
	Message string `json:"message"`
}

func (c *commands) onTrump(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	personalize := data.String("personalize")
	var apiURL string
	if len(personalize) > 0 {
		apiURL = fmt.Sprintf("https://api.whatdoestrumpthink.com/api/v1/quotes/personalized?q=%s", url.QueryEscape(personalize))
	} else {
		apiURL = "https://api.whatdoestrumpthink.com/api/v1/quotes/random"
	}

	resp, err := c.HTTPClient.Get(apiURL)
	if err != nil {
		slog.Error("Error fetching Trump quote", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
		})
	}

	var trumpQuote TrumpQuote
	if err := json.Unmarshal(body, &trumpQuote); err != nil {
		slog.Error("Error unmarshalling response body", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("%s - Donald Trump", trumpQuote.Message),
	})
}
