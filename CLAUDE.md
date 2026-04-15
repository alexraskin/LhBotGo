# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A Discord bot for the LhCloudy Twitch community, built in Go using the disgo library. It tracks guesses about what "Lh" stands for, provides streamer info, and has fun commands. Data is stored in MongoDB.

## Commands

```bash
# Run locally (requires MongoDB)
go run main.go -config config.toml

# Build
go build ./...

# Run all tests
go test ./...

# Run a specific package's tests
go test ./lhbot/...
```

## Architecture

**Entry point:** `main.go` loads config, creates dependencies (Discord client, MongoDB, HTTP client), wires them into the `Bot` struct, registers event listeners, and starts the gateway.

**Key wiring pattern:**
- `lhbot/bot.go` — Bot struct holds all dependencies; `Start()` syncs slash commands globally then opens the Discord gateway
- `lhbot/commands/commands.go` — Collects all `ApplicationCommandCreate` definitions in `Commands` slice and builds a `handler.Router` with grouped routes (e.g. `/lh/*`, `/q/*`, `/lhcloudy/*`)
- `lhbot/handlers.go` — Non-slash event handlers (`!shatter` text command, ready event for bot presence)
- `lhbot/database/mongo.go` — `MongoClient` interface for guess CRUD operations

**Command groups:** `/lh/*` (guess tracking), `/lhcloudy/*` (streamer info), `/ow/*` (Overwatch), `/fun/*` (random images), `/q/*` (viewer game queue), `/stats`, `/help`

**Config:** Loaded from TOML file with environment variable overrides (env takes precedence). See `config.toml.example` and `.env.example` for available keys.

## Key Dependencies

- **disgo** — Discord bot framework (slash commands, gateway, event handling)
- **go-away** — Profanity filter for guess submissions
- **mongo-driver v2** — MongoDB persistence
- **BurntSushi/toml** — Config file parsing

## Deployment

Docker multi-platform build via GitHub Actions. CI runs tests and build on all pushes to main and PRs. Deploy workflow auto-increments semver patch tags and pushes to GHCR.
