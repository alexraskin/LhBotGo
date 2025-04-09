package commands

import (
	"encoding/json"
	"io"

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
	},
}

type CatData struct {
	ID        string   `json:"id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	URL       string   `json:"url"`
}

func (c *commands) onCat(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	apiURL := "https://cataas.com/cat?json=true"

	resp, err := c.HTTPClient.Get(apiURL)
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
	apiURL := "https://dog.ceo/api/breeds/image/random"

	resp, err := c.HTTPClient.Get(apiURL)
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
	PostLink string   `json:"postLink"`
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler  bool     `json:"spoiler"`
	Author    string   `json:"author"`
	Ups       int      `json:"ups"`
	Preview   []string `json:"preview"`
}


func (c *commands) onMeme(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	apiURL := "https://meme-api.com/gimme"

	resp, err := c.HTTPClient.Get(apiURL)
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
