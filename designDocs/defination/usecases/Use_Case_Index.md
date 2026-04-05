# Use Case Index

## Actor Definitions

| Actor | Type | Description |
|---|---|---|
| Scheduler | System | Cron-based trigger that initiates periodic jobs |
| Bot Operator | Human | Discord user with Administrator permission who manages group purchases (enforced via `DefaultMemberPermissions`) |
| Guild Member | Human | Discord user who interacts with the bot via slash commands |
| Notion API | External System | Data source for user records and transaction data |
| Discord API | External System | Delivery channel for user notifications, logging, and thread management |

---

## Use Case List

| ID | Use Case Name | Trigger | Primary Actor | Description | Status |
|---|---|---|---|---|---|
| [UC-001](UC-001_Notify_Unpaid_Users.md) | Notify Unpaid Users | Cron schedule | Scheduler | **Deprecated by UC-004.** Checks each user's unpaid amount in Notion and sends a Discord DM reminder if the total exceeds 2,000 TWD | Deprecated |
| [UC-004](UC-004_Trigger_Debt_Reminder.md) | Trigger Debt Reminder | `/debt-reminder` slash command | Bot Operator | Replaces UC-001. Immediately runs unpaid notification (debug or production mode) and schedules a one-shot production run after N days | Draft |
| [UC-002](UC-002_Create_New_Order.md) | Create New Order | `/newOrder` slash command | Bot Operator | Creates a Discord thread for a group purchase order and inserts a tracking record into the Notion Order List database (TBL-004); restricted to authorized operator only | Draft |
| [UC-003](UC-003_Register_Buy_Record.md) | Register Buy Record | `/buy` slash command (reply) | Guild Member | Registers a purchase record into a member's personal transaction database (TBL-002) with JPY amount and auto-calculated TWD | Draft |

---

## ID Assignment Ranges

| Range | Category |
|---|---|
| UC-001 – UC-099 | Payment, Notification & Order Management |
| UC-900 – UC-999 | Infrastructure (Logging, Auth) |

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/22 | — | Initial version — UC-001 registered |
| 1.1 | 2026/03/18 | — | Add UC-002 (Create New Order), add Guild Member actor |
| 1.2 | 2026/03/18 | — | Add UC-003 (Register Buy Record) |
| 1.3 | 2026/04/05 | — | Add UC-004 (Trigger Debt Reminder), deprecate UC-001 |
