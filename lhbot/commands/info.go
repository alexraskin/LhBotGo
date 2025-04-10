package commands

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/dustin/go-humanize"
)

var statsStartTime = time.Now()

var infoCommand = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Get some info about the bot",
}

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

func (c *commands) onInfo(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "disgo: \t%s\n", disgo.Version)
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	fmt.Fprintf(w, "Uptime: \t%s\n", getDurationString(time.Since(statsStartTime)))
	fmt.Fprintf(w, "Memory used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Concurrent tasks: \t%s\n", humanize.Comma(int64(runtime.NumGoroutine())))
	fmt.Fprintf(w, "```\n")
	w.Flush()
	out := buf.String()

	end := "built with ❤️ by twizycat"
	return e.CreateMessage(discord.MessageCreate{
		Content: out + end,
	})
}
