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
- Sends reminders via Discord DMs (personal DB: when unpaid exceeds per-currency threshold TWD > 2000, JPY > 8000; others DB: when any unpaid amount exists)

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
  │   └── user.go
  ├── port/            ← interfaces between layers
  │   ├── repository.go  (UserRepository)
  │   └── notifier.go    (Notifier)
  ├── usecase/         ← business logic
  │   └── notify_unpaid.go
  └── gateway/         ← external service adapters
      ├── notion/        (implements port.UserRepository)
      │   └── user_repository.go
      └── discord/       (implements port.Notifier)
          └── notifier.go
```

Dependency rule: `gateway` → `domain`, `usecase` → `port` → `domain`. No layer imports `gateway` or `usecase` except `main.go`.

---

## Key Design Decisions

- **Scheduling:** gocron with Asia/Tokyo timezone; cron expression from `WORKER_CORNTAB` env var
- **Notification threshold:** Personal DB: per-currency (TWD > 2000, JPY > 8000) defined as `notificationAmountLimit` in `usecase/notify_unpaid.go`; Others DB: any amount > 0
- **Notification schedule:** Users are notified on the 1st and 15th of each month
- **Debug mode:** `DEBUG=1` disables actual Discord DMs (logs only) and fakes time to the 1st via `clock.NewMock()`
- **No channel passing:** use case calls notifier directly; one failure does not block other users (logged, not fatal)

---

## Required Environment Variables

| Variable | Purpose |
|---|---|
| `DISCORD_TOKEN` | Discord bot token |
| `DISCORD_GUILD` | Guild name (informational) |
| `DISCORD_GUILD_LOG_CHANNEL_ID` | Channel ID for logging |
| `NOTION_TOKEN` | Notion API token |
| `NOTION_USER_DB_ID` | Notion user database ID |
| `NOTION_OTHERS_DB_ID` | Notion shared "其他" database ID |
| `WORKER_CORNTAB` | Cron schedule (e.g. `*/1 * * * *`) |
| `DEBUG` | Set to `1` to enable debug mode |

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

`add_debug_mode` — adds debug mode feature (prevents real DMs when `DEBUG=1`)

Main branch is `main`.
