# TBL-002: Personal Transaction Database

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-002 |
| Table Name | Personal Transaction Database |
| Notion DB ID | Per-member (referenced by `notion_id` in TBL-001) |
| Version | 1.0 |
| Status | Draft |
| Date | 2026/02/23 |
| Author | — |

---

## 1. Overview

Each member has their own personal transaction database in Notion. It records individual purchase transactions and their payment status. The system queries this database to calculate the member's total unpaid amount.

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `付款狀況` | Select | Yes | Payment status of the transaction |
| `台幣` | Number | Conditional | Amount in TWD (used when member's currency is `TWD`) |
| `日幣` | Number | Conditional | Amount in JPY (used when member's currency is `JPY`) |

---

## 3. Column Details

### `付款狀況`

- **Type:** Select
- **Used Filter Value:** `尚未付款` (unpaid)
- **Note:** The system filters records where this column equals `尚未付款` to calculate the unpaid total

### `台幣`

- **Type:** Number
- **Unit:** TWD (New Taiwan Dollar)
- **Condition:** Read when the member's `currency` (TBL-001) is `TWD`
- **Example:** `1500`

### `日幣`

- **Type:** Number
- **Unit:** JPY (Japanese Yen)
- **Condition:** Read when the member's `currency` (TBL-001) is `JPY`
- **Example:** `3000`

---

## 4. Query Logic

1. Filter: `付款狀況` equals `尚未付款`
2. Sum: the amount column matching the member's currency (`台幣` for TWD, `日幣` for JPY)
3. Compare against per-currency notification threshold:

| Currency | Threshold |
|---|---|
| TWD | > 2,000 |
| JPY | > 8,000 |

---

## 5. Related Tables

| Table | Relationship |
|---|---|
| TBL-001: User Database | Each TBL-002 instance is referenced by a member's `notion_id` in TBL-001 |

---

## 6. Usage

- Read by `gateway/notion/user_repository.go` → `GetUnpaidAmount()`
- Currency-to-column mapping defined in `currencyColumnMap`

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/02/23 | — | Initial draft |
