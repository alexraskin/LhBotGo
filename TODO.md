# TODO

## Bugs

- [x] **Database query field mismatch** — `database/mongo.go:85` queries `bson.M{"guess": guess}` but the BSON tag is `"lhguess"`. `GetGuess()` will never find a match.
- [x] **Deprecated rand pattern** — `handlers.go:44` and `overwatch.go:73` call `rand.New(rand.NewSource(...))` but discard the result and use the global `rand` anyway. Remove the unused call and use `rand.IntN()` directly.

## Missing Error Handling

- [x] **No HTTP status code checks** — `fun.go:38,85,139` don't check `resp.StatusCode` before unmarshalling. A 404/500 response will silently fail.
- [x] **Discord message send errors ignored** — handler functions (e.g. `handlers.go`) call `CreateMessage` without checking the returned error.

## Hardcoded Values

- [x] **Hardcoded Discord user ID** — `handlers.go:46` had `"127122091139923968"` inline. Moved to config as `owner_id` / `BOT_OWNER_ID`.
- [x] **Magic color constant** — `0x5865F2` (Discord blurple) was repeated 5 times. Extracted to `embedColor` constant in `commands.go`.

## Code Quality

- [x] **Global mutable queue state** — `queue.go:15-17` uses package-level `var queue []string` with a mutex. Moved onto the `commands` struct.
- [x] **Stats uptime captured at import time** — `stats.go:18` sets `statsStartTime = time.Now()` at package load, not bot start. Moved to `commands` struct, initialized in `New()`.
- [x] **Verbose var declaration** — `guess.go:54` uses `var hint []string = []string{...}` instead of `hint := []string{...}`.

## Testing

- [ ] Add tests for database operations
- [ ] Add tests for command handlers
- [ ] Add tests for HTTP API calls (cat/dog/meme)
