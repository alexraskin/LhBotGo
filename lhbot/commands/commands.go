package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"

	"github.com/alexraskin/LhBotGo/lhbot"
)

var Commands = []discord.ApplicationCommandCreate{
	statsCommand,
	guessCommands,
	overwatchCommands,
	lhCloudyCommands,
	funCommands,
	helpCommand,
}

type commands struct {
	*lhbot.Bot
}

func New(b *lhbot.Bot) handler.Router {
	cmds := &commands{b}

	router := handler.New()
	router.Use(middleware.Go)
	router.SlashCommand("/stats", cmds.onStats)
	router.SlashCommand("/help", cmds.onHelp)
	router.Route("/lh", func(r handler.Router) {
		r.SlashCommand("/guess", cmds.onGuess)
		r.SlashCommand("/count", cmds.onCount)
		r.SlashCommand("/list", cmds.onList)
		r.SlashCommand("/hint", cmds.onHint)
		r.SlashCommand("/latest", cmds.onLatest)
	})
	router.Route("/ow", func(r handler.Router) {
		r.SlashCommand("/shatter", cmds.onShatter)
		r.SlashCommand("/reinquote", cmds.onQuote)
	})
	router.Route("/lhcloudy", func(r handler.Router) {
		r.SlashCommand("/birthday", cmds.onBirthday)
		r.SlashCommand("/from", cmds.onFrom)
		r.SlashCommand("/youtube", cmds.onYoutube)
		r.SlashCommand("/twitter", cmds.onTwitter)
		r.SlashCommand("/tips", cmds.onTips)
		r.SlashCommand("/code", cmds.onCode)
		r.SlashCommand("/lhfurry", cmds.onLhfurry)
		r.SlashCommand("/instagram", cmds.onInstagram)
		r.SlashCommand("/age", cmds.onAge)
		r.SlashCommand("/interview", cmds.onInterview)
		r.SlashCommand("/socials", cmds.onLinks)
	})
	router.Route("/fun", func(r handler.Router) {
		r.SlashCommand("/cat", cmds.onCat)
		r.SlashCommand("/dog", cmds.onDog)
		r.SlashCommand("/meme", cmds.onMeme)
	})
	return router
}
