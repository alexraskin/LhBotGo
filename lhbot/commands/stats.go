package commands

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/dustin/go-humanize"
)

var statsStartTime = time.Now()

var statsCommand = discord.SlashCommandCreate{
	Name:        "stats",
	Description: "Get some stats about the bot",
}

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

func (c *commands) onStats(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	guesses, err := c.Bot.Mongo.CountGuesses(c.Bot.Ctx, c.Bot.DBName, c.Bot.Collection)
	if err != nil {
		slog.Error("Error getting guesses", "error", err)
		return e.CreateMessage(discord.MessageCreate{
			Content: "Error getting guesses",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error getting hostname", "error", err)
		hostname = "unknown"
	}
	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "Go: \t%s\n", c.Bot.Version.GoVersion)
	fmt.Fprintf(w, "disgo: \t%s\n", disgo.Version)
	fmt.Fprintf(w, "LhBotGo: \t%s\n", c.Bot.Version.Version)
	fmt.Fprintf(w, "Host: \t%s\n", hostname)
	fmt.Fprintf(w, "Uptime: \t%s\n", getDurationString(time.Since(statsStartTime)))
	fmt.Fprintf(w, "Memory used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Concurrent tasks: \t%s\n", humanize.Comma(int64(runtime.NumGoroutine())))
	fmt.Fprintf(w, "Number of Guesses: \t%d\n", guesses)
	fmt.Fprintf(w, "```\n")
	w.Flush()
	out := buf.String()

	end := "built with ❤️ by twizycat"
	return e.CreateMessage(discord.MessageCreate{
		Content: out + end,
	})
}
