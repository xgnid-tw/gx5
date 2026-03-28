# UC-002: Create New Order

## Document Metadata

| Item | Value |
|---|---|
| Use Case ID | UC-002 |
| Use Case Name | Create New Order |
| Version | 1.2 |
| Status | Draft |
| Date | 2026/03/18 |
| Author | — |

---

## 1. Use Case Overview

### Purpose

Allow the bot operator to create a new group purchase order via a Discord slash command, automatically setting up a discussion thread and tracking the order in Notion.

### Summary

The bot operator executes the `/newOrder` slash command with order details. The system verifies the caller is the authorized operator (BR-015), creates a Discord thread titled with the order name, posts shop URL / tag mentions / deadline as the first message, and inserts a corresponding record into the Notion Order List database (TBL-004).

### Scope

**In scope:**
- Parsing slash command parameters (`orderTitle`, `deadline`, `shopURL`, `tags`) — all required
- Creating a Discord thread in the channel where the command was issued
- Sending the formatted first message in the created thread
- Inserting a new record into the Notion Order List database (TBL-004)

**Out of scope:**
- Editing or deleting existing orders
- Closing or archiving threads
- Validating the shop URL content
- Managing order fulfillment or payment status

---

## 2. Actor Information

### Primary Actor

| Actor | Role |
|---|---|
| Bot Operator | The authorized Discord user (ID: `374867612519366657`) who manages group purchases |

### Secondary Actor

None.

### System Actor

| System | Role |
|---|---|
| Discord API | Creates the thread and posts the first message |
| Notion API | Persists the order record to the Order List database (TBL-004) |

---

## 3. Pre-conditions and Post-conditions

### Pre-conditions

- The Discord bot is authenticated and connected to the guild
- The Notion Order List database (TBL-004, configured via `NOTION_ORDER_DB_ID` env var) is accessible
- The `/newOrder` slash command is registered with the Discord application
- The invoking user's Discord ID matches the authorized operator (BR-015)

### Post-conditions

**On success:**
- A new Discord thread exists in the invoking channel with the title `orderTitle`
- The thread's first message contains the shop URL, tag mentions, and deadline (BR-008)
- A new record exists in the Notion Order List database with `threadName`, `deadline`, and `tags` populated

**On failure:**
- If the invoking user is not the authorized operator → command is rejected with an error; no action taken
- If thread creation fails → error response returned to the user; no Notion record is created
- If Notion record creation fails → the thread already exists but the order is not tracked in Notion; error is logged

---

## 4. Business Flows

### Summary Flow

1. User executes `/newOrder orderTitle deadline shopURL tags` in a Discord channel (all parameters required)
2. System verifies the invoking user is the authorized operator (BR-015); if not, reject with an error
3. System validates required parameter `orderTitle` is present (BR-007)
4. System creates a new Discord thread in the current channel with title `orderTitle`
5. System sends the first message in the thread with the format defined in BR-008
6. System inserts a new record into the Notion Order List database (TBL-004) with `threadName` = `orderTitle`, `deadline` = `deadline`, `tags` = `tags` (BR-009)
7. System responds to the slash command interaction confirming success

### Detailed Business Flows

At this time, no specific business usage calling this function has been identified; therefore, a detailed business flow definition is not provided.

---

## 5. Business Rules

| ID | Rule Name | Description | Exception |
|---|---|---|---|
| BR-007 | Required Parameters | All parameters (`orderTitle`, `deadline`, `shopURL`, `tags`) are required; Discord enforces this at the command level | None |
| BR-008 | Thread First Message Format | The first message in the created thread follows the format: line 1 = `shopURL`, line 2 = tag mentions (each prefixed with `@`), line 3 = deadline display (`截止時間: {deadline}`) | None (all fields are always present) |
| BR-009 | Notion Record Mapping | The Notion record maps as follows: `threadName` ← `orderTitle` (Title), `deadline` ← `deadline` (Date, ISO-8601), `tags` ← `tags` (Select, single value) | `shopURL` is not stored in Notion (TBL-004 has no such column) |
| BR-010 | Tag Values | Tag must correspond to a valid select option defined in TBL-004: `315pro`, `学マス`, `283pro`, `346pro`, `765pro` (single value only) | Unknown tag is passed as-is; Notion API will reject invalid values |
| BR-015 | Operator Authorization | Only the authorized operator (configured via `DISCORD_OWNER_ID` env var) may execute this command; all other users are rejected | None |

---

## 6. Related Use Cases

None.

---

## 7. Supplementary Information

### Expected Usage Frequency

- On-demand, triggered by the bot operator only
- Estimated: a few times per week during active group purchase periods
- No peak hours expected

### Operations and Maintenance Requirements

- Notion DB column name changes in TBL-004 (`threadName`, `deadline`, `tags`) require corresponding code updates in the order repository gateway
- Adding new tag options requires updating the Notion database multi-select configuration
- The slash command must be registered with Discord during bot startup

### Other Notes

- The `shopURL` parameter is only used in the Discord thread message and is not persisted to Notion
- Thread creation and Notion record insertion are sequential — if thread creation succeeds but Notion insertion fails, the thread will exist without a tracking record
- The `deadline` parameter accepts ISO-8601 date format; the Discord message displays it in a human-readable form
- The `tags` parameter accepts comma-separated tag names that map to both Discord mentions and Notion multi-select values

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/03/18 | — | Initial draft |
| 1.1 | 2026/03/18 | — | Restrict command to authorized operator only (BR-015, configured via `DISCORD_OWNER_ID` env var); update actor, pre-conditions, and flow |
| 1.2 | 2026/03/28 | — | All parameters (`orderTitle`, `deadline`, `shopURL`, `tags`) are now required; `tags` uses Discord Choices dropdown (BR-010) |
