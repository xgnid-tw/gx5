# TBL-002: Personal Transaction Database

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-002 |
| Table Name | Personal Transaction Database |
| Notion DB ID | Per-member (referenced by `notion_id` in TBL-001) |
| Version | 2.0 |
| Status | Draft |
| Date | 2026/03/18 |
| Author | — |

---

## 1. Overview

Each member has their own personal transaction database in Notion. It records individual purchase transactions and their payment status. The system queries this database to calculate the member's total unpaid amount.

**Parent page:** GX小精靈Ver.4

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `品項` | Title | Yes | Name/description of the purchased item |
| `台幣` | Number | Conditional | Amount in TWD (used when member's currency is `TWD`) |
| `日幣` | Number | Conditional | Amount in JPY (used when member's currency is `JPY`) |
| `付款狀況` | Select | Yes | Payment status of the transaction |
| `物品狀況` | Select | No | Item delivery/fulfillment status |
| `購買途徑` | Select | No | Store or platform where the item was purchased |
| `連結` | URL | No | Link to the product page or order |
| `備註` | Rich Text | No | Free-text notes |
| `預計到貨` | Date | No | Expected arrival date (single date or date range) |
| `建立時間` | Created Time | Auto | Record creation timestamp (auto-generated) |

---

## 3. Column Details

### `品項`

- **Type:** Title (primary column)
- **Format:** Free text
- **Note:** Identifies the purchased item; not used in query logic

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

### `付款狀況`

- **Type:** Select
- **Allowed Values:**

| Value | Color | Description |
|---|---|---|
| `已付款` | yellow | Paid |
| `尚未付款` | pink | Unpaid |

- **Used Filter Value:** `尚未付款` (unpaid)
- **Note:** The system filters records where this column equals `尚未付款` to calculate the unpaid total

### `物品狀況`

- **Type:** Select
- **Allowed Values:**

| Value | Color | Description |
|---|---|---|
| `已回台` | pink | Arrived in Taiwan |
| `未訂購` | yellow | Not yet ordered |
| `已下訂` | purple | Order placed |
| `已到貨` | default | Item received (in Japan) |
| `代付` | gray | Paid on behalf |
| `運費` | blue | Shipping fee |

- **Note:** Not used in query logic; for manual tracking only

### `購買途徑`

- **Type:** Select
- **Allowed Values:**

| Value | Color |
|---|---|
| `阿搜比` | purple |
| `HMV` | yellow |
| `aniplex+` | default |
| `任天堂官網` | brown |
| `アニメイトストア` | red |
| `Music Rain Store` | green |
| `店頭` | gray |
| `總之就是官網` | orange |
| `メルカリ` | blue |
| `夏音(?)` | pink |
| `ヤフオク` | green |
| `Booth` | gray |
| `駿河屋オンラインショップ` | pink |
| `Gift Online Shop` | blue |
| `あみあみオンラインショップ` | brown |
| `CyStore` | red |
| `列印` | pink |
| `おしながき` | yellow |

- **Note:** Not used in query logic; for manual tracking only

### `連結`

- **Type:** URL
- **Format:** Product page or order URL
- **Note:** Not used in query logic

### `備註`

- **Type:** Rich Text
- **Format:** Free text
- **Note:** Not used in query logic

### `預計到貨`

- **Type:** Date
- **Format:** ISO-8601 date; supports date ranges (start + end)
- **Note:** Not used in query logic; for manual tracking only

### `建立時間`

- **Type:** Created Time
- **Format:** ISO-8601 datetime (auto-generated)
- **Note:** Automatically set when the record is created

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
| 2.0 | 2026/03/18 | — | Add missing columns from Notion schema: `品項`, `物品狀況`, `購買途徑`, `連結`, `備註`, `預計到貨`, `建立時間`; add allowed values for `付款狀況`, `物品狀況`, `購買途徑` |
