# Claude Code Guidelines

## 1. Core Principles

**Role:** Act as a calm and objective technical expert.
**Tone:** Use concise and clear language. Do not use emotional modifiers, empathetic expressions, or excessive honorifics.
**Language:** Use English in this project (gx5). Use Japanese by default in all other projects.

## 2. Strict Prohibitions

### Prohibited in Conversation and Expression

The following "unnecessary preambles" and "emotional fillers" are strictly prohibited. Responses must begin immediately with the conclusion, solution, or confirmation items.

Prohibited phrases:
- "I hope this helps"
- "Thank you for your question"
- "That's a great question"
- "I'm sorry, but"
- "Understood"
- "Yes, let me explain about..."

### Prohibited in Information Processing (Anti-Hallucination Measures)

- **No guessing:** Do not fill in gaps when context, code, or specifications are unclear.
- **No fabrication:** Never create nonexistent libraries, methods, syntax, or facts.
- **Eliminate uncertainty:** When information is insufficient, do not attempt to answer. Specifically identify the missing information and ask the user to provide it.

## 3. Code Review

All code will be reviewed by Codex and Copilot. Write clean, reviewable code.

## 4. Output Format

- Limit explanations to technical facts; avoid verbose commentary.
- Use code blocks and bullet points heavily to improve readability.
- Maintain logical structure: **Conclusion → Reasoning → Examples/Code**.

---

# GX5 - Project Context for Claude

## Project Purpose

A personal Discord bot written in Go that automates payment tracking for group purchases. It:
- Pulls transaction data from a Notion database
- Calculates how much each person owes
- Sends reminders via Discord DMs when unpaid exceeds per-currency threshold (TWD > 2000, JPY > 8000) or any unpaid record is older than 3 months

Not intended for public use — private utility for a specific Discord guild ("GX小精靈倉庫").

---

## Tech Stack

- **Language:** Go 1.25.3
- **Module:** `github.com/xgnid-tw/gx5`
- **Key Libraries:**
  - `github.com/bwmarrin/discordgo` — Discord API
  - `github.com/go-co-op/gocron/v2` — Cron job scheduling (Asia/Tokyo timezone)
  - `github.com/jomei/notionapi` — Notion API
  - `github.com/joho/godotenv` — `.env` loading
  - `github.com/benbjohnson/clock` — Testable time abstraction

---

## Architecture

Clean Architecture with 4 layers. Dependencies point inward only.

```
main.go (wiring)
  ├── config/          ← env loading & validation
  ├── domain/          ← entities (no external deps)
  │   ├── user.go
  │   ├── order.go
  │   └── transaction.go
  ├── port/            ← interfaces between layers
  │   ├── repository.go      (UserRepository)
  │   ├── notifier.go        (Notifier)
  │   ├── order_repository.go (OrderRepository)
  │   ├── order_creator.go   (OrderCreator)
  │   ├── transaction_repository.go (TransactionRepository)
  │   └── buy_record.go      (BuyRecordRegisterer)
  ├── usecase/         ← business logic
  │   ├── notify_unpaid.go
  │   ├── create_order.go
  │   └── register_buy_record.go
  └── gateway/         ← external service adapters
      ├── notion/        (implements port.UserRepository, OrderRepository, TransactionRepository)
      │   ├── user_repository.go
      │   ├── order_repository.go
      │   └── transaction_repository.go
      └── discord/
          ├── notifier.go       (implements port.Notifier)
          ├── thread_creator.go (implements port.ThreadCreator)
          └── command/          (slash/message command handlers)
              ├── command_handler.go
              ├── respond.go
              ├── new_order_handler.go
              └── buy_command.go
```

Dependency rule: `gateway` → `domain`, `usecase` → `port` → `domain`. No layer imports `gateway` or `usecase` except `main.go`.

---

## Coding Conventions

### Discord Command Handlers (`gateway/discord/command/`)

- **Registration pattern:** All commands must use `RegisterXxxCommand(ch *Handler, uc port.Interface)`. The function creates the command definition, registers it with the handler, and wires the interaction callback internally. Do not expose `NewXxxCommand()` + `HandleXxx()` as separate functions.
- **Port interfaces:** Command handlers must depend on port interfaces (e.g., `port.OrderCreator`), never on concrete usecase types. Define the interface in `port/` with an `Execute` method matching the usecase signature.
- **Response helpers:** Use `respondError(s, i, msg)` for ephemeral error messages and `respondSuccess(s, i, msg)` for visible success messages. Both are defined in `gateway/discord/command/respond.go`. Do not create per-file response helpers.
- **User-facing text:** All Discord-visible text (command descriptions, option descriptions, modal titles/labels, error messages, success messages) must be in Chinese. Command `Name` fields and option `Name` fields must remain ASCII (Discord requirement).
- **Named constants:** Command names, modal prefixes, and input custom IDs must be `const` at the top of the file. No magic strings in handler logic.

### Business Logic

- **No hardcoded business values:** Configurable values (exchange rates, thresholds, IDs) must come from environment variables loaded in `config/`. Only use constants for truly fixed values (e.g., valid tag names defined in domain).

### Test Conventions (`gateway/notion/`)

- **Hand-rolled mocks** for `notionapi` services (`mockDatabaseService`, `mockPageService`) live in `gateway/notion/mock_test.go`. Do not define mocks in individual test files.

---

## Key Design Decisions

- **Scheduling:** gocron with Asia/Tokyo timezone; `/debt-reminder` slash command triggers immediate run + one-shot delayed run
- **Notification threshold:** Per-currency (TWD > 2000, JPY > 8000) defined as `notificationAmountLimit` in `usecase/notify_unpaid.go`; also notifies if any unpaid record is older than 3 months and total > 0
- **Debug mode:** `/debt-reminder debug:true` sends reminders to log channel only (no DMs); the delayed run always uses production mode
- **No channel passing:** use case calls notifier directly; one failure does not block other users (logged, not fatal)

---

## Required Environment Variables

| Variable | Purpose |
|---|---|
| `DISCORD_TOKEN` | Discord bot token |
| `DISCORD_APP_ID` | Discord application ID (for slash command registration) |
| `DISCORD_GUILD_ID` | Discord guild (server) ID |
| `DISCORD_GUILD` | Guild name (informational) |
| `DISCORD_GUILD_LOG_CHANNEL_ID` | Channel ID for logging |
| `NOTION_TOKEN` | Notion API token |
| `NOTION_USER_DB_ID` | Notion user database ID |
| `NOTION_ORDER_DB_ID` | Notion Order List database ID (TBL-004) |
| `EXCHANGE_RATE_JPY_TWD` | JPY to TWD exchange rate (e.g. `0.24`) |
| `TAG_ROLE_MAP` | Comma-separated tag=roleID pairs for Discord role mentions |

---

## Build & Run

```bash
go mod download   # download dependencies
go run main.go    # run directly (requires .env)
go build          # compile binary
```

## Linting & Formatting

```bash
./bin/golangci-lint run   # lint
gofumpt                   # format
```

---

## CI/CD (CircleCI)

- **build:** downloads deps, creates `.env` from CircleCI env vars
- **lint:** runs golangci-lint
- **deploy:** rsync to `ssh.xgnid.space` → `gx5/` directory (main branch only)
- SSH fingerprint: `SHA256:NPj4IcXxqQEKGXOghi/QbG2sohoNfvZ30JwCcdSSNM0`

---

## Notion Database Schema (relevant fields)

- User DB: has `DiscordID`, `Name`, `NotionID`, `Currency` (TWD or JPY)
- Transaction records: filtered by unpaid status `"尚未付款"`, amount column determined by user's currency (`台幣` for TWD, `日幣` for JPY)

---

## DM Message Format

```
[欠費提醒] https://www.notion.so/{notionID} (如果有漏登聯絡一下XG)
```

---

## Current Branch

Main branch is `main`.
