# UC-003 Design: Optional Item Name & Environment Exchange Rate

## Summary

Two requirement changes to UC-003 (Register Buy Record):

1. **Optional item name field** — add a second TextInput to the buy modal, pre-filled with the thread title. If the user clears or changes it, use their input; if left as-is or empty, fall back to thread title.
2. **Exchange rate from env** — replace the hardcoded `JPYToTWDRate = 0.24` constant with `EXCHANGE_RATE_JPY_TWD` env var loaded at startup.

## Change 1: Optional Item Name in Modal

### Current Behavior

- Modal has 1 field: JPY Amount
- `品項` is always the thread title (`channel.Name`)
- Thread title is encoded in modal CustomID: `buy_modal:<discordID>:<threadTitle>`

### New Behavior

- Modal has 2 fields: JPY Amount + 品項 (item name)
- Item name TextInput:
  - `CustomID`: `"item_name"`
  - `Label`: `"品項"`
  - `Style`: `TextInputShort`
  - `Required`: `false`
  - `Value`: pre-filled with thread title
- In `handleBuyModal`: if item name input is empty, use thread title from CustomID; otherwise use user input

### Impact

- **Usecase layer**: No change — `Execute` already accepts `itemName string`
- **Gateway layer**: `buy_command.go` — add TextInput, extract value in modal handler
- **Domain layer**: No change

## Change 2: Exchange Rate from Environment

### Current Behavior

- `const JPYToTWDRate = 0.24` in `usecase/register_buy_record.go`
- Rate is hardcoded; changes require code modification and redeployment

### New Behavior

- `EXCHANGE_RATE_JPY_TWD` env var, parsed as `float64`
- Loaded in `config.Config` with validation (must be > 0)
- Passed to `NewRegisterBuyRecord` as a constructor parameter
- Hardcoded constant removed

### Impact

- **Config layer**: `config.go` — add `ExchangeRateJPYTWD float64` field, parse and validate
- **Usecase layer**: `register_buy_record.go` — accept rate as param, store in struct, remove const
- **Wiring**: `main.go` — pass `cfg.ExchangeRateJPYTWD` to `NewRegisterBuyRecord`
- **Tests**: `register_buy_record_test.go` — update constructor calls with explicit rate

## Files Changed

| File | Change |
|------|--------|
| `config/config.go` | Add `ExchangeRateJPYTWD` field, env loading, validation |
| `usecase/register_buy_record.go` | Accept rate as constructor param, remove `JPYToTWDRate` const |
| `usecase/register_buy_record_test.go` | Update constructor calls with rate param |
| `gateway/discord/command/buy_command.go` | Add item name TextInput to modal, handle default in modal handler |
| `main.go` | Pass `cfg.ExchangeRateJPYTWD` to usecase constructor |
| `designDocs/defination/usecases/UC-003_Register_Buy_Record.md` | Update BR-011, BR-012 |

## Business Rule Updates

| Rule | Before | After |
|------|--------|-------|
| BR-011 | `JPY × 0.24` (hardcoded constant) | `JPY × EXCHANGE_RATE_JPY_TWD` (from env var) |
| BR-012 | `品項` = thread title (always) | `品項` = user-provided item name, defaulting to thread title if empty |

## New Environment Variable

| Variable | Type | Required | Default | Purpose |
|----------|------|----------|---------|---------|
| `EXCHANGE_RATE_JPY_TWD` | float64 | Yes | — | JPY to TWD exchange rate for buy record calculation |
