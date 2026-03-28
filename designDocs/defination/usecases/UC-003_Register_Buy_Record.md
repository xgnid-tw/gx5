# UC-003: Register Buy Record

## Document Metadata

| Item | Value |
|---|---|
| Use Case ID | UC-003 |
| Use Case Name | Register Buy Record |
| Version | 1.0 |
| Status | Draft |
| Date | 2026/03/18 |
| Author | — |

---

## 1. Use Case Overview

### Purpose

Allow the operator to quickly register a purchase record for a guild member by replying to that member's message with `/buy`, entering the JPY amount, and having the system automatically insert a transaction record into the member's personal Notion database.

### Summary

A guild member executes the `/buy` slash command as a reply to another member's message within a thread. The bot prompts for the total amount in JPY. Upon receiving the amount, the bot looks up the replied-to member's personal transaction database (TBL-002) via TBL-001, inserts a new unpaid record with the thread title as the item name, and confirms completion.

### Scope

**In scope:**
- Identifying the target member from the replied-to message author
- Prompting for and receiving the JPY amount
- Looking up the target member's personal Notion database (TBL-002) via TBL-001
- Calculating the TWD equivalent using a fixed exchange rate
- Inserting a new transaction record into the target member's TBL-002
- Replying with confirmation message

**Out of scope:**
- Editing or deleting existing transaction records
- Handling members who only exist in TBL-003 (其他 database)
- Dynamic exchange rate fetching (uses a constant)
- Payment status updates

---

## 2. Actor Information

### Primary Actor

| Actor | Role |
|---|---|
| Guild Member | Discord user who initiates the `/buy` command by replying to a message |

### Secondary Actor

None.

### System Actor

| System | Role |
|---|---|
| Discord API | Provides the replied-to message author, thread title, and interactive message flow |
| Notion API | Looks up the target member in TBL-001, inserts record into the member's TBL-002 |

---

## 3. Pre-conditions and Post-conditions

### Pre-conditions

- The Discord bot is authenticated and connected to the guild
- The `/buy` command is registered with the Discord application
- The command is executed as a reply to another member's message within a thread
- The replied-to member exists in TBL-001 with a valid `notion_id` pointing to a TBL-002 instance
- The Notion user database (TBL-001) and target personal transaction database (TBL-002) are accessible

### Post-conditions

**On success:**
- A new record exists in the target member's TBL-002 with:
  - `品項` = thread title
  - `日幣` = user-input JPY amount
  - `台幣` = JPY amount × 0.24
  - `付款狀況` = `尚未付款`
- The bot has replied "登記完畢" in the thread

**On failure:**
- If the replied-to member is not found in TBL-001 → error response to user; no record created
- If Notion record insertion fails → error is logged; user is informed
- If the command is not used as a reply → error response to user

---

## 4. Business Flows

### Summary Flow

1. Guild member replies to a target member's message with `/buy` in a thread
2. System extracts the replied-to message author's Discord ID
3. System prompts the guild member: "請輸入金額 (JPY)" (or equivalent prompt for the total amount)
4. Guild member inputs the JPY amount
5. System looks up the target member in TBL-001 by matching `discord_id`
6. System retrieves the target member's `notion_id` (TBL-002 database ID)
7. System retrieves the current thread title from Discord
8. System calculates TWD amount = JPY amount × 0.24 (BR-011)
9. System inserts a new record into the target member's TBL-002 (BR-012, BR-013)
10. System replies "登記完畢" in the thread

### Detailed Business Flows

At this time, no specific business usage calling this function has been identified; therefore, a detailed business flow definition is not provided.

---

## 5. Business Rules

| ID | Rule Name | Description | Exception |
|---|---|---|---|
| BR-011 | Fixed Exchange Rate | TWD amount is calculated as `JPY amount × 0.24`; the exchange rate is a constant defined in code | None |
| BR-012 | Item Name from Thread Title | The `品項` column is populated with the Discord thread title where the `/buy` command was executed | None |
| BR-013 | Default Payment Status | New records are always created with `付款狀況` = `尚未付款` (unpaid) | None |
| BR-014 | Target Member Lookup | The target member is identified by the Discord ID of the replied-to message author; this ID is matched against `discord_id` in TBL-001 to resolve the member's `notion_id` (TBL-002 database ID) | If the replied-to user is not found in TBL-001, the operation fails with an error |

---

## 6. Related Use Cases

None.

---

## 7. Supplementary Information

### Expected Usage Frequency

- On-demand, triggered by guild members during group purchase periods
- Estimated: several times per week
- No peak hours expected

### Operations and Maintenance Requirements

- The exchange rate constant (0.24) is hardcoded; changes require code modification and redeployment
- TBL-001 must be kept up to date with correct `discord_id` → `notion_id` mappings
- Adding new members requires creating a new TBL-002 database in Notion and registering the member in TBL-001

### Other Notes

- The command must be used as a reply to a message; using it without replying is an error
- The command must be used within a thread to resolve the thread title for `品項`
- The interactive prompt (asking for amount) implies a multi-step Discord interaction (e.g., modal or follow-up message collector)
- Only the `日幣` and `台幣` amount columns are populated; the currency is always JPY input with an auto-calculated TWD equivalent

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/03/18 | — | Initial draft |
