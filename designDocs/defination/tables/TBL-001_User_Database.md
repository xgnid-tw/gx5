# TBL-001: User Database

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-001 |
| Table Name | User Database |
| Notion DB ID | `NOTION_USER_DB_ID` |
| Version | 1.1 |
| Status | Draft |
| Date | 2026/02/23 |
| Author | — |

---

## 1. Overview

Stores registered member information. Each row represents a member of the Discord guild who participates in group purchase tracking.

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `discord_id` | Title | Yes | Discord user ID used for sending DMs |
| `name` | Rich Text | Yes | Display name of the member |
| `notion_id` | Rich Text | Yes | Database ID of the member's personal transaction database (TBL-002) |
| `currency` | Select | Yes | Currency code for the member's transactions |

---

## 3. Column Details

### `discord_id`

- **Type:** Title (primary column)
- **Format:** Discord snowflake ID (numeric string)
- **Example:** `"123456789012345678"`

### `name`

- **Type:** Rich Text
- **Format:** Free text
- **Example:** `"Alice"`

### `notion_id`

- **Type:** Rich Text
- **Format:** Notion database UUID (without hyphens)
- **Example:** `"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"`
- **Note:** References a personal transaction database (TBL-002) unique to each member

### `currency`

- **Type:** Select
- **Allowed Values:**

| Value | Description | Maps to Transaction Column |
|---|---|---|
| `TWD` | New Taiwan Dollar | `台幣` |
| `JPY` | Japanese Yen | `日幣` |

- **Note:** Determines which amount column is read from the member's transaction database

---

## 4. Related Tables

| Table | Relationship |
|---|---|
| TBL-002: Personal Transaction Database | `notion_id` references a TBL-002 instance per member |

---

## 5. Usage

- Read by `gateway/notion/user_repository.go` → `GetUsers()`
- Maps to `domain.User` struct

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/23 | — | Initial draft |
| 1.1 | 2026/02/23 | — | Fix `currency` column type: Rich Text → Select |
