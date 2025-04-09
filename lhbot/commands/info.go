package commands

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/shirou/gopsutil/v3/mem"
)

var infoCommand = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Get some info about the bot",
}

type infoMessageData struct {
	Version          string
	GoVersion        string
	Uptime           string
	MemoryUsed       float64
	MemoryTotal      float64
	GarbageCollected float64
	ConcurrentTasks  int
}

func (c *commands) onInfo(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	uptime := time.Since(c.Bot.StartTime).Round(time.Second)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	vMem, _ := mem.VirtualMemory()

	gcStats := debug.GCStats{}
	debug.ReadGCStats(&gcStats)

	usedMemGB := float64(memStats.Alloc) / 1e9
	totalMemGB := float64(vMem.Total) / 1e9
	garbageCollectedGB := gcStats.PauseTotal.Seconds() * float64(memStats.TotalAlloc) / 1e9

	info := infoMessageData{
		Version:          c.Bot.Version,
		GoVersion:        c.Bot.GoVersion,
		Uptime:           uptime.String(),
		MemoryUsed:       usedMemGB,
		MemoryTotal:      totalMemGB,
		GarbageCollected: garbageCollectedGB,
		ConcurrentTasks:  runtime.NumGoroutine(),
	}

	stats := fmt.Sprintf(
		"```\nVersion: %s\nGo Version: %s\nUptime: %s\nMemory Used: %.1f GB / %.1f GB (%.0f GB garbage collected)\nConcurrent Tasks: %d\n```\nMade with ❤️ by <@297398689415168000>",
		info.Version,
		info.GoVersion,
		info.Uptime,
		info.MemoryUsed,
		info.MemoryTotal,
		info.GarbageCollected,
		info.ConcurrentTasks,
	)

	return e.CreateMessage(discord.MessageCreate{
		Content: stats,
		Flags:   discord.MessageFlagEphemeral,
	})
}
