# UC-001: Notify Unpaid Users

## Document Metadata

| Item | Value |
|---|---|
| Use Case ID | UC-001 |
| Use Case Name | Notify Unpaid Users |
| Version | 1.2 |
| Status | Draft |
| Date | 2026/02/23 |
| Author | — |

---

## 1. Use Case Overview

### Purpose

Automatically remind users with outstanding payment obligations to settle their balances, eliminating the need for manual tracking and follow-up within the group.

### Summary

The system periodically checks each registered user's unpaid amount in Notion (from either the user's personal transaction database or the shared "其他" database, depending on routing) and sends a Discord DM reminder to users with qualifying unpaid balances.

### Scope

**In scope:**
- Fetching all registered users from the Notion user database
- Calculating each user's unpaid amount from either their personal Notion database (TBL-002) or the shared "其他" database (TBL-003), based on routing (BR-005)
- Sending payment reminder DMs via Discord
- Logging sent reminders to the guild log channel

**Out of scope:**
- Marking records as paid
- Recording new purchase transactions
- Handling user replies to the reminder DM
- Modifying Notion data

---

## 2. Actor Information

### Primary Actor

| Actor | Role |
|---|---|
| Scheduler | Cron-based system trigger that initiates the use case periodically |

### Secondary Actor

None.

### System Actor

| System | Role |
|---|---|
| Notion API | Data source — provides user list and unpaid transaction records |
| Discord API | Delivery channel — sends DMs to users and logs to the guild channel |

---

## 3. Pre-conditions and Post-conditions

### Pre-conditions

- The scheduler is running with a valid `WORKER_CORNTAB` cron expression
- The Discord bot is authenticated and connected to the guild
- The Notion user database is accessible and contains at least one user record
- Each user record contains valid values for `discord_id`, `notion_id`, `name`, and `currency`

### Post-conditions

**On success:**
- Qualifying users (personal DB exceeding threshold per BR-001, or others DB with any unpaid amount per BR-006) have received a Discord DM
- Each sent reminder is logged to the guild log channel

**On failure:**
- Errors are logged to the application log
- Per-user DM failure does not affect other users (isolated)
- Repository-level failures (Notion API error) terminate the current run; the next scheduled execution retries

---

## 4. Business Flows

### Summary Flow

1. Scheduler triggers the use case execution
2. If today is not the 1st or 15th of the month → exit immediately (BR-002)
3. Fetch all users from the Notion user database
4. For each user:
   1. If `notion_id` equals `NOTION_OTHERS_DB_ID` → query the shared "其他" database for unpaid records matching the user's `name` (BR-005)
   2. Else → query the user's personal Notion database for unpaid records (BR-003)
   3. If personal DB user: notify when total exceeds the per-currency threshold (BR-001)
   4. If others DB user: notify when any unpaid amount exists (BR-006)
   5. Send Discord DM and log to guild channel
5. If DM fails for one user → log error, continue to next user (BR-004)

### Detailed Business Flows

Refer to:
- [BF-001-1: Notify Unpaid Users — Normal Flow](../flow/BF-001-1_Notify_Unpaid_Users_Normal.md)
- [BF-001-2: Notify Unpaid Users — Notion API Error](../flow/BF-001-2_Notify_Unpaid_Users_Notion_Error.md)
- [BF-001-3: Notify Unpaid Users — Discord DM Failure](../flow/BF-001-3_Notify_Unpaid_Users_DM_Failure.md)

---

## 5. Business Rules

| ID | Rule Name | Description | Exception |
|---|---|---|---|
| BR-001 | Personal DB Notification Threshold | For personal DB users, a reminder is sent only when the unpaid amount exceeds the per-currency threshold (TWD > 2,000, JPY > 8,000) | None |
| BR-002 | Bi-monthly Reminder Frequency | All users are evaluated on the 1st and 15th of each month | None |
| BR-003 | Unpaid Status Filter | Only records with `付款狀況 = 尚未付款` are included in the amount calculation | None |
| BR-004 | DM Failure Isolation | A failure to send a DM to one user does not stop the notification process for remaining users | None |
| BR-005 | Others Table Routing | Users whose `notion_id` equals `NOTION_OTHERS_DB_ID` have their unpaid amount calculated from the shared "其他" database (TBL-003) by matching `購買人` to the user's `name` | None |
| BR-006 | Others DB Notification Threshold | For others DB users, a reminder is sent when any unpaid amount exists (amount > 0) | None |

---

## 6. Related Use Cases

None.

---

## 7. Supplementary Information

### Expected Usage Frequency

- Determined by `WORKER_CORNTAB` environment variable
- Development: every minute (`*/1 * * * *`)
- Production: recommended monthly or daily depending on operational needs
- Peak: no peak hours expected; execution is lightweight

### Operations and Maintenance Requirements

- Notion DB column name changes (`付款狀況`, `台幣`, `日幣`, `購買人`, `discord_id`, `notion_id`, `name`, `currency`) require corresponding code updates in `gateway/notion/user_repository.go`
- Discord bot token rotation requires updating `DISCORD_TOKEN` in `.env` and redeploying
- Debug mode (`DEBUG=1`) suppresses actual Discord DMs — must be disabled in production

### Other Notes

- The personal DB notification thresholds (TWD > 2,000, JPY > 8,000) are defined as `notificationAmountLimit` in `usecase/notify_unpaid.go` — change requires code modification, not config
- Others DB users are notified for any unpaid amount (no threshold)
- The day-of-month guard uses injectable `clock.Clock` for testability
- The "其他" database ID is configured via `NOTION_OTHERS_DB_ID` environment variable

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/22 | — | Initial draft |
| 1.1 | 2026/02/23 | — | Fix BR-001 (per-currency threshold), BR-002 (1st and 15th), BR-003 (remove nonexistent age filter, correct to status filter), add BR-005 (其他 table), update references |
| 1.2 | 2026/02/23 | — | Split threshold rules: BR-001 scoped to personal DB, add BR-006 (others DB notifies on any amount > 0), fix summary/scope to clarify exclusive routing |
