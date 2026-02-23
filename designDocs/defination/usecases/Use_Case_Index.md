# Use Case Index

## Actor Definitions

| Actor | Type | Description |
|---|---|---|
| Scheduler | System | Cron-based trigger that initiates periodic jobs |
| Notion API | External System | Data source for user records and transaction data |
| Discord API | External System | Delivery channel for user notifications and logging |

---

## Use Case List

| ID | Use Case Name | Trigger | Primary Actor | Description | Status |
|---|---|---|---|---|---|
| [UC-001](UC-001_Notify_Unpaid_Users.md) | Notify Unpaid Users | Cron schedule | Scheduler | Checks each user's unpaid amount in Notion and sends a Discord DM reminder if the total exceeds 2,000 TWD | Draft |

---

## ID Assignment Ranges

| Range | Category |
|---|---|
| UC-001 – UC-099 | Payment & Notification |
| UC-900 – UC-999 | Infrastructure (Logging, Error Handling, Auth) |

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/22 | — | Initial version — UC-001 registered |
