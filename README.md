# LhBot

[![LhBot](https://i.gyazo.com/632f0e60dc0535128971887acad98993.png)](https://twitter.com/PetraYle)

[![Twitch Status](https://img.shields.io/twitch/status/lhcloudy27?color=6441a5&logo=twitch&logoColor=white)](https://www.twitch.tv/lhcloudy27)
[![Go Version](https://img.shields.io/github/go-mod/go-version/alexraskin/LhBotGo)](https://golang.org/doc/devel/release.html)


## Why did I make LhBot?

If you’ve ever watched LhCloudy on Twitch, you know one thing for sure, he refuses to reveal what “Lh” actually stands for. So, I created LhBot: a fun Discord bot that keeps track of everyone’s wild (and often hilarious) guesses.

## Connect with LhCloudy

- [Twitch](https://www.twitch.tv/lhcloudy27)
- [Twitter](https://twitter.com/LhCloudy)
- [Discord](https://discord.com/invite/jd6CZSj8jb)
- [Youtube](https://www.youtube.com/channel/UC2CV-HWvIrMO4mUnYtNS-7A)

## Run The Bot Locally

You will need a mongodb instance to run the bot locally.

```bash
cp config.toml.example config.toml
```

```bash
go run main.go -config config.toml
```
