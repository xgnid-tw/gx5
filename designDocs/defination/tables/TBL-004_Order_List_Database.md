# TBL-004: Order List Database (orderList)

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-004 |
| Table Name | Order List Database |
| Notion DB ID | Configured via `NOTION_ORDER_DB_ID` env var |
| Version | 1.2 |
| Status | Draft |
| Date | 2026/03/18 |
| Author | — |

---

## 1. Overview

Tracks group purchase order threads. Each row represents an order with its deadline and associated franchise/series tags. Used to manage order scheduling within the Discord guild.

**Parent page:** GX小精靈Ver.4

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `threadName` | Title | Yes | Name of the order thread |
| `deadline` | Date | No | Order deadline (single date or date range) |
| `tags` | Select | No | Category tag indicating the franchise/series (single value) |

---

## 3. Column Details

### `threadName`

- **Type:** Title (primary column)
- **Format:** Free text
- **Note:** Identifies the order thread

### `deadline`

- **Type:** Date
- **Format:** ISO-8601 date or datetime; supports date ranges (start + end)
- **Note:** Represents the cutoff date for the order

### `tags`

- **Type:** Select
- **Allowed Values:**

| Value | Color | Description |
|---|---|---|
| `315pro` | default | 315 Production |
| `学マス` | blue | 学園アイドルマスター |
| `283pro` | brown | 283 Production |
| `346pro` | purple | 346 Production |
| `765pro` | gray | 765 Production |

---

## 4. Related Tables

| Table | Relationship |
|---|---|
| — | No direct relationships to other tables |

---

## 5. Usage

- Not currently referenced in application code
- Serves as a standalone order tracking database in Notion

---

**Revision History**

| Version | Date | Author | Description |
|---|---|---|---|
| 1.0 | 2026/03/18 | — | Initial draft |
| 1.1 | 2026/03/18 | — | Replace hardcoded Notion DB ID with `NOTION_ORDER_DB_ID` env var |
| 1.2 | 2026/03/18 | — | Fix `tags` column type: Multi Select → Select (single value) per actual Notion schema |
