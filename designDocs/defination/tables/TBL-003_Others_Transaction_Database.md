# TBL-003: Others Transaction Database (其他)

## Document Metadata

| Item | Value |
|---|---|
| Table ID | TBL-003 |
| Table Name | Others Transaction Database |
| Notion DB ID | Configured via `NOTION_OTHERS_DB_ID` env var |
| Version | 2.0 |
| Status | Draft |
| Date | 2026/03/18 |
| Author | — |

---

## 1. Overview

A shared transaction database for infrequent buyers who do not have their own personal transaction database. Records are attributed to users via the `購買人` (buyer name) column. Users whose `notion_id` (TBL-001) equals `NOTION_OTHERS_DB_ID` have their unpaid amount calculated exclusively from this table.

**Parent page:** GX小精靈Ver.4

---

## 2. Column Definition

| Column | Notion Type | Required | Description |
|---|---|---|---|
| `品項` | Title | Yes | Name/description of the purchased item |
| `購買人` | Select | Yes | Name of the buyer (matched against `name` in TBL-001) |
| `台幣` | Number | Conditional | Amount in TWD |
| `日幣` | Number | Conditional | Amount in JPY |
| `付款狀況` | Select | Yes | Payment status of the transaction |
| `物品狀況` | Select | No | Item delivery/fulfillment status |
| `購買途徑` | Select | No | Store or platform where the item was purchased |
| `連結` | URL | No | Link to the product page or order |
| `備註` | Rich Text | No | Free-text notes |
| `預計到貨` | Date | No | Expected arrival date (single date or date range) |
| `建立時間` | Created Time | Auto | Record creation timestamp (auto-generated) |
| `建立時間 (1)` | Created Time | Auto | Duplicate creation timestamp (auto-generated, legacy) |

---

## 3. Column Details

### `品項`

- **Type:** Title (primary column)
- **Format:** Free text
- **Note:** Item description; not used in query logic

### `購買人`

- **Type:** Select
- **Allowed Values:**

| Value | Color |
|---|---|
| `歐魯` | orange |
| `坂木` | gray |
| `魚肉` | default |
| `Gahe` | green |
| `蛋頭` | pink |
| `貓熊` | blue |
| `robo` | orange |
| `冷氣` | red |
| `海葵` | yellow |
| `烏龜` | red |
| `貴貴` | brown |
| `琉璃` | purple |
| `彥彥` | yellow |
| `睡狼` | green |
| `SAA` | blue |
| `蛋糕` | pink |
| `貢丸` | orange |
| `褚喵` | blue |
| `阿秋` | pink |
| `小響` | red |
| `蝶芙` | default |
| `茶茶` | yellow |
| `眼爸` | pink |
| `アルバ` | green |
| `阿星` | gray |
| `芙蕾` | green |
| `小張` | orange |
| `西蒙` | pink |

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
- **Allowed Values:**

| Value | Color | Description |
|---|---|---|
| `已付款` | yellow | Paid |
| `尚未付款` | pink | Unpaid |

- **Used Filter Value:** `尚未付款` (unpaid)
- **Note:** The system filters records where this column equals `尚未付款`

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
| `轉寄轉運` | orange | Forwarding/transshipment |

- **Note:** Not used in query logic; for manual tracking only. Has one additional value (`轉寄轉運`) compared to TBL-002.

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
| `總之就是官網` | orange |
| `店頭` | gray |
| `メルカリ` | blue |
| `夏音(?)` | pink |
| `ヤフオク` | green |
| `Booth` | gray |
| `駿河屋オンラインショップ` | pink |
| `Gift Online Shop` | blue |
| `あみあみオンラインショップ` | brown |
| `CyStore` | red |
| `ゲーマーズオンラインショップ` | blue |
| `虎之穴` | pink |
| `7-11Net` | red |
| `楽天ショップ` | orange |
| `melonbooks` | pink |
| `列印` | blue |
| `樂天` | gray |
| `eplus` | gray |
| `キャラボム` | purple |
| `paypay free market` | brown |
| `GEO` | yellow |
| `楽天` | red |
| `a-on` | purple |
| `コトフギア` | purple |
| `ゲオ` | red |
| `amazon` | orange |
| `にじさんじ` | gray |
| `KT` | yellow |
| `passmarket` | gray |
| `しまむら` | gray |
| `vv` | gray |
| `Zozotown` | purple |
| `ラブライブ School idol Store` | default |
| `PB` | red |
| `ソフトマップ` | default |

- **Note:** Not used in query logic; for manual tracking only. Superset of TBL-002 values.

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

### `建立時間 (1)`

- **Type:** Created Time
- **Format:** ISO-8601 datetime (auto-generated)
- **Note:** Duplicate of `建立時間`; likely a legacy column

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
| 2.0 | 2026/03/18 | — | Add missing columns from Notion schema: `物品狀況`, `購買途徑`, `連結`, `備註`, `預計到貨`, `建立時間`, `建立時間 (1)`; add all allowed values for `付款狀況`, `物品狀況`, `購買人`, `購買途徑` |
