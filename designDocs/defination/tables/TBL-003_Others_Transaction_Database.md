# TBL-003: Others Transaction Database (其他)

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-003 |
| Table Name | Others Transaction Database |
| Notion DB ID | Configured via `NOTION_OTHERS_DB_ID` env var |
| Version | 1.2 |
| Status | Draft |
| Date | 2026/02/23 |
| Author | — |

---

## 1. Overview

A shared transaction database for infrequent buyers who do not have their own personal transaction database. Records are attributed to users via the `購買人` (buyer name) column. Users whose `notion_id` (TBL-001) equals `NOTION_OTHERS_DB_ID` have their unpaid amount calculated exclusively from this table.

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `品項` | Title | Yes | Name of the purchased item |
| `購買人` | Select | Yes | Name of the buyer (matched against `name` in TBL-001) |
| `台幣` | Number | Conditional | Amount in TWD |
| `日幣` | Number | Conditional | Amount in JPY |
| `付款狀況` | Select | Yes | Payment status of the transaction |

---

## 3. Column Details

### `品項`

- **Type:** Title (primary column)
- **Note:** Item description; not used in query logic

### `購買人`

- **Type:** Select
- **Note:** Must match the `name` column value in TBL-001 for the user to be associated with this record

### `台幣`

- **Type:** Number
- **Unit:** TWD (New Taiwan Dollar)
- **Condition:** Read when the member's `currency` (TBL-001) is `TWD`

### `日幣`

- **Type:** Number
- **Unit:** JPY (Japanese Yen)
- **Condition:** Read when the member's `currency` (TBL-001) is `JPY`

### `付款狀況`

- **Type:** Select
- **Used Filter Value:** `尚未付款` (unpaid)
- **Note:** The system filters records where this column equals `尚未付款`

---

## 4. Query Logic

1. Filter: `購買人` equals the user's `name` (from TBL-001) **AND** `付款狀況` equals `尚未付款`
2. Sum: the amount column matching the member's currency (`台幣` for TWD, `日幣` for JPY)
3. Notify if any unpaid amount exists (BR-006)

---

## 5. Related Tables

| Table | Relationship |
|---|---|
| TBL-001: User Database | `購買人` is matched against `name` in TBL-001; routing determined by `notion_id` |

---

## 6. Usage

- Read by `gateway/notion/user_repository.go` → `GetOthersUnpaidAmount()`
- Uses `notionapi.AndCompoundFilter` to combine `購買人` and `付款狀況` filters
- Currency-to-column mapping shared with TBL-002 via `currencyColumnMap`

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/23 | — | Initial draft |
| 1.1 | 2026/02/23 | — | Fix query logic: TBL-003 is queried exclusively (not summed with TBL-002); correct routing description |
| 1.2 | 2026/02/23 | — | Fix query logic step 3: reference BR-006 (amount > 0), not BR-001 (per-currency threshold) |
