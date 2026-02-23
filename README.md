# gx5

A private Discord bot written in Go for automating payment tracking within a group. It reads unpaid transaction records from Notion and sends monthly Discord DM reminders to members whose outstanding balance exceeds the threshold.

> Not intended for public use — private utility for a specific Discord guild.

---

## How It Works

On the 1st and 15th of each month, the bot:

1. Fetches all registered members from the Notion user database
2. For each member, queries their personal Notion database for unpaid records
3. Sums the unpaid amounts (using the column matching the user's currency)
4. Sends a Discord DM if the total exceeds the per-currency threshold (TWD > 2,000, JPY > 8,000)
5. Logs every sent reminder to a designated guild log channel

---

## Architecture

Built with Clean Architecture (Ports and Adapters):

```
main.go          ← wiring only
config/          ← env loading and validation
domain/          ← User entity
port/            ← interfaces
usecase/         ← business logic
gateway/
  notion/        ← implements Repository via Notion API
  discord/       ← implements Notifier via Discord API
```

---

## Requirements

- Go 1.25+
- A Discord bot token with DM and guild permissions
- A Notion integration token with access to the user database

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values:

| Variable                       | Description                                               |
| ------------------------------ | --------------------------------------------------------- |
| `DISCORD_TOKEN`                | Discord bot token                                         |
| `DISCORD_GUILD_LOG_CHANNEL_ID` | Channel ID for logging sent reminders                     |
| `NOTION_TOKEN`                 | Notion integration token                                  |
| `NOTION_USER_DB_ID`            | Notion database ID for the user list                      |
| `WORKER_CORNTAB`               | Cron expression for the scheduler (e.g. `0 9 1 * *`)      |
| `DEBUG`                        | Set to any non-empty value to suppress actual Discord DMs |

---

## Notion Database Schema

### User Database (`NOTION_USER_DB_ID`)

| Column       | Type      | Description                                      |
| ------------ | --------- | ------------------------------------------------ |
| `discord_id` | Title     | Discord user ID                                  |
| `name`       | Rich Text | Member name                                      |
| `notion_id`  | Rich Text | ID of the member's personal transaction database |
| `currency`   | Rich Text | Currency code (`TWD` or `JPY`)                   |

### Personal Transaction Database (per member)

| Column     | Type   | Description                               |
| ---------- | ------ | ----------------------------------------- |
| `付款狀況` | Select | Payment status — filter value: `尚未付款` |
| `台幣`     | Number | Amount in TWD (for TWD users)             |
| `日幣`     | Number | Amount in JPY (for JPY users)             |

---

## Local Development

```bash
# Install dependencies
go mod download

# Run
go run main.go

# Run with debug mode (no DMs sent)
DEBUG=1 go run main.go
```

### Linting

```bash
./bin/golangci-lint run
```

### Testing

```bash
go test ./...
```

Mocks are generated from `port/` interfaces using [mockery](https://github.com/vektra/mockery):

```bash
go tool mockery
```

---

## CI/CD

Managed with CircleCI:

| Job      | Trigger      | Action                       |
| -------- | ------------ | ---------------------------- |
| `build`  | All branches | Download deps, create `.env` |
| `lint`   | All branches | Run `golangci-lint`          |
| `deploy` | `main` only  | rsync to `ssh.xgnid.space`   |

---

## Design Documents

Located in [`designDocs/defination/`](designDocs/defination/):

```
guideline/     ← document authoring guidelines
usecases/      ← UC-XXX use case definitions
flow/          ← BF-XXX business flow definitions
functions/     ← functional requirements index
```
