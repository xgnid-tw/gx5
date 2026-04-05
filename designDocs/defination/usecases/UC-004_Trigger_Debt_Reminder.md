# UC-004: Trigger Debt Reminder

## Document Metadata

| Item | Value |
|---|---|
| Use Case ID | UC-004 |
| Use Case Name | Trigger Debt Reminder |
| Version | 1.0 |
| Status | Draft |
| Date | 2026/04/05 |
| Author | — |

---

## 1. Use Case Overview

### Purpose

Allow the bot operator to manually trigger a debt reminder via a Discord slash command, replacing the previous cron-based automatic reminder (UC-001). The command runs the reminder immediately and schedules a second production run after a configurable number of days.

### Summary

The bot operator executes `/debt-reminder` with optional `days` and `debug` parameters. The system immediately runs the unpaid notification logic — in debug mode (log channel only) or production mode (actual DMs) depending on the `debug` flag. It then schedules a one-shot delayed job to run the same logic in production mode after `days` days, at the same time of day.

### Scope

**In scope:**
- Parsing slash command parameters (`days`, `debug`)
- Immediate execution of the unpaid notification logic
- Scheduling a one-shot delayed execution in production mode
- Debug mode toggle per invocation (affects immediate run only)

**Out of scope:**
- Recurring/cron-based scheduling (UC-001 is deprecated by this use case)
- Cancelling or listing scheduled jobs
- Modifying notification thresholds or user lists

---

## 2. Actor Information

### Primary Actor

| Actor | Role |
|---|---|
| Bot Operator | Discord user with Administrator permission who triggers the reminder |

### Secondary Actor

None.

### System Actor

| System | Role |
|---|---|
| Notion API | Data source — provides user list and unpaid transaction records |
| Discord API | Delivery channel — sends DMs to users and logs to the guild channel |
| gocron Scheduler | Manages the delayed one-shot job |

---

## 3. Pre-conditions and Post-conditions

### Pre-conditions

- The Discord bot is authenticated and connected to the guild
- The Notion user database is accessible and contains at least one user record
- The `/debt-reminder` slash command is registered with the Discord application
- The invoking user has Administrator permission in the Discord guild (BR-020)

### Post-conditions

**On success:**
- The unpaid notification logic has executed immediately (debug or production mode per BR-017)
- A one-shot job is scheduled to execute in production mode after `days` days (BR-018)
- The operator receives a confirmation message with execution summary

**On failure:**
- If the immediate run fails → error is reported to the operator; the delayed job is still scheduled
- If scheduling the delayed job fails → error is reported to the operator; the immediate run has already completed
- Per-user DM failures are isolated (same as UC-001 BR-004)

---

## 4. Business Flows

### Summary Flow

1. Bot Operator executes `/debt-reminder` in a Discord channel (optionally with `days` and `debug`)
2. Discord enforces command visibility to administrators only (BR-020)
3. System sends a deferred interaction response (Discord shows "thinking..." indicator)
4. System executes the unpaid notification logic immediately:
   - If `debug=true` → send reminders to log channel only, skip DMs (BR-017)
   - If `debug=false` → send reminders as DMs and log to guild channel
5. System schedules a one-shot job to run `days` days from now at the same time, in production mode (BR-018)
6. System edits the deferred response confirming:
   - Immediate run result (debug or production, number of users notified)
   - Scheduled production run date/time

### Detailed Business Flows

At this time, no specific business usage calling this function has been identified; therefore, a detailed business flow definition is not provided.

---

## 5. Business Rules

| ID | Rule Name | Description | Exception |
|---|---|---|---|
| BR-017 | Debug Mode (Immediate Run) | When `debug=true`, the immediate run sends reminders to the log channel only (no DMs). The delayed run is always production mode regardless of this flag. | None |
| BR-018 | Delayed One-Shot Execution | The system schedules a one-shot job to run exactly `days × 24 hours` after the command is issued. The delayed run always uses production mode. | If the bot restarts before the delayed job fires, the job is lost (no persistence) |
| BR-019 | Default Parameter Values | `days` defaults to 15; `debug` defaults to false | None |
| BR-020 | Operator Authorization | Command visibility is restricted via Discord's `DefaultMemberPermissions` (Administrator). Only server administrators can see and execute this command. | Fine-tune per-user/per-role in Discord Server Settings → Integrations → Bot → Command Permissions |
| BR-021 | Notification Thresholds | Same as UC-001: personal DB users notified when unpaid exceeds per-currency threshold (TWD > 2,000, JPY > 8,000); others DB users notified when any unpaid amount exists | None |
| BR-022 | DM Failure Isolation | A failure to send a DM to one user does not stop the notification process for remaining users (same as UC-001 BR-004) | None |

---

## 6. Related Use Cases

| Use Case | Relationship |
|---|---|
| UC-001 Notify Unpaid Users | **Deprecated by** this use case. UC-004 replaces the cron-based trigger with an on-demand slash command. The core notification logic (threshold checks, DM sending) is reused. |

---

## 7. Supplementary Information

### Expected Usage Frequency

- On-demand, triggered by the bot operator
- Estimated: once or twice per month
- No peak hours expected

### Operations and Maintenance Requirements

- The `WORKER_CORNTAB` environment variable is no longer required and should be removed from config
- The day guard (1st/15th check) in `NotifyUnpaid` is removed — the operator controls when to run
- The `clock` dependency in `NotifyUnpaid` is removed — no longer needed without the day guard
- Debug mode is now per-invocation (slash command parameter), not a global env var for this use case

### Deprecation: UC-001

The following components are removed:
- Cron job setup in `main.go`
- `WORKER_CORNTAB` env var and its config validation
- Day guard (`day != 1 && day != 15`) in `usecase/notify_unpaid.go`
- Debug-mode clock faking in `main.go`
- `clock.Clock` field from `NotifyUnpaid` struct

### Other Notes

- The delayed job is not persisted — if the bot restarts, the scheduled job is lost. This is acceptable for the current use case since the operator can re-issue the command.
- The `Notifier` must support per-call debug toggling. The `debug` field can be passed as a parameter to `Execute` rather than baked into the notifier at construction.

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/04/05 | — | Initial draft |
