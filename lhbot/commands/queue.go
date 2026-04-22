package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var queueCommands = discord.SlashCommandCreate{
	Name:        "q",
	Description: "Viewer game queue commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "join",
			Description: "Join the Viewer Games Queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "leave",
			Description: "Leave the Viewer Games Queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List current Viewer Games Queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "played",
			Description: "Mark a user as played and remove from the queue (admin only)",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "User to mark as played",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "clear",
			Description: "Clear the Viewer Games Queue (admin only)",
		},
	},
}

func guildOnlyMessage() discord.MessageCreate {
	return discord.MessageCreate{Content: "The Viewer Games Queue can only be used in a server.", Flags: discord.MessageFlagEphemeral}
}

func (c *commands) onQueueJoin(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guildID := e.GuildID()
	if guildID == nil {
		return e.CreateMessage(guildOnlyMessage())
	}
	userID := e.User().ID.String()
	c.queueMu.Lock()
	defer c.queueMu.Unlock()
	if slices.Contains(c.queues[*guildID], userID) {
		return e.CreateMessage(discord.MessageCreate{Content: "You're already in the Viewer Games Queue."})
	}
	c.queues[*guildID] = append(c.queues[*guildID], userID)
	return e.CreateMessage(discord.MessageCreate{Content: "You've joined the Viewer Games Queue!"})
}

func (c *commands) onQueueLeave(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guildID := e.GuildID()
	if guildID == nil {
		return e.CreateMessage(guildOnlyMessage())
	}
	userID := e.User().ID.String()
	c.queueMu.Lock()
	defer c.queueMu.Unlock()
	queue := c.queues[*guildID]
	for i, id := range queue {
		if id == userID {
			c.queues[*guildID] = slices.Delete(queue, i, i+1)
			return e.CreateMessage(discord.MessageCreate{Content: "You have left the queue."})
		}
	}
	return e.CreateMessage(discord.MessageCreate{Content: "You are not in the Viewer Games Queue."})
}

func (c *commands) onQueueList(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guildID := e.GuildID()
	if guildID == nil {
		return e.CreateMessage(guildOnlyMessage())
	}
	c.queueMu.Lock()
	defer c.queueMu.Unlock()
	queue := c.queues[*guildID]
	if len(queue) == 0 {
		return e.CreateMessage(discord.MessageCreate{Content: "Viewer Games Queue is Empty."})
	}
	msg := "Current Viewer Games Queue:\n"
	for i, id := range queue {
		msg += fmt.Sprintf("%d. <@%s>\n", i+1, id)
	}
	return e.CreateMessage(discord.MessageCreate{Content: msg})
}

func (c *commands) onQueuePlayed(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guildID := e.GuildID()
	if guildID == nil {
		return e.CreateMessage(guildOnlyMessage())
	}
	if !e.Member().Permissions.Has(discord.PermissionModerateMembers) {
		return e.CreateMessage(discord.MessageCreate{Content: "You don't have permission to do that.", Flags: discord.MessageFlagEphemeral})
	}
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

	user, ok := data.Resolved.Users[userID]
	if !ok {
		slog.Error("Could not resolve user from Discord API", "userID", userID)
		return e.CreateMessage(discord.MessageCreate{
			Content: "An error occurred - KEKL",
		})
	}
	c.queueMu.Lock()
	defer c.queueMu.Unlock()
	queue := c.queues[*guildID]
	for i, id := range queue {
		if id == user.ID.String() {
			c.queues[*guildID] = slices.Delete(queue, i, i+1)
			return e.CreateMessage(discord.MessageCreate{Content: fmt.Sprintf("<@%s> has been marked as played and removed from the queue.", user.ID)})
		}
	}
	return e.CreateMessage(discord.MessageCreate{Content: "User is not in the Viewer Games Queue."})
}

func (c *commands) onQueueClear(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	guildID := e.GuildID()
	if guildID == nil {
		return e.CreateMessage(guildOnlyMessage())
	}
	if !e.Member().Permissions.Has(discord.PermissionModerateMembers) {
		return e.CreateMessage(discord.MessageCreate{Content: "You don't have permission to do that.", Flags: discord.MessageFlagEphemeral})
	}
	c.queueMu.Lock()
	delete(c.queues, *guildID)
	c.queueMu.Unlock()
	return e.CreateMessage(discord.MessageCreate{Content: "Viewer Games Queue has been cleared."})
}

func (c *commands) onQueueHelp(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	embed := discord.Embed{
		Title:       "Queue Help",
		Description: "Help for the Viewer Games Queue command.",
		Fields: []discord.EmbedField{
			{
				Name:  "Join",
				Value: "Join the Viewer Games Queue.",
			},
			{
				Name:  "Leave",
				Value: "Leave the Viewer Games Queue.",
			},
			{
				Name:  "List",
				Value: "List the Viewer Games Queue.",
			},
			{
				Name:  "Played",
				Value: "Mark a user as played and remove them from the Viewer Games Queue (admin only).",
			},
			{
				Name:  "Clear",
				Value: "Clear the Viewer Games Queue (admin only).",
			},
		},
	}
	return e.CreateMessage(discord.MessageCreate{Embeds: []discord.Embed{embed}})
}
