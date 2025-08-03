# Open Ledger Format — Personal Finance (v2.0)

*A minimal, human‑readable schema for everyday finances that keeps every cent accounted for.*

---

## Table of contents

1. [Overview](#1-overview)
2. [Data Model](#2-data-model)
3. [Validation Rules (invariants)](#3-validation-rules-invariants)
4. [Example (YAML)](#4-example-yaml)
5. [Versioning](#5-versioning)
6. [License](#6-license)

---

## 1. Overview

The **Open Ledger Format (OLF)** stores personal‑finance data in plain‑text files—**YAML, JSON**.

### Design principles

1. **Human‑oriented** — anyone can maintain the file with a plain‑text editor.
2. **Consistency** — validation rules ensure strict, deterministic balances across the log.
3. **High‑level focus** — summary budgeting (≈ \$100–\$1 000 chunks), not coffee‑by‑coffee logging.

---

## 2 Data Model

```text
Ledger (root)
└─ Years{int}          → Year
     └─ Months{int}    (1–12) → Month
          └─ Accounts{string} → Account
               └─ Entries[]    → Entry
```

All monetary amounts are **signed integers**. A ledger adopts one currency unit for the whole file—typically *whole units* such as dollars for a cleaner high‑level view, but cents or another smallest unit are also acceptable if used consistently. Positive numbers add funds; negative numbers spend them.

### 2.1 `Ledger`

* **years** (*map\[int]Year*) — dictionary keyed by calendar year.

### 2.2 `Year`

* **opening\_balance** (*int*) — balance at 00 : 00 on 1 Jan.
* **closing\_balance** (*int*) — balance at 23 : 59 on 31 Dec.
* **months** (*map\[int]Month*) — twelve slots indexed 1–12 (missing keys allowed).

### 2.3 `Month`

* **opening\_balance** (*int*) — balance at month start.
* **closing\_balance** (*int*) — balance at month end.
* **accounts** (*map\[string]Account*) — dictionary keyed by account name.

### 2.4 `Account`

* **opening\_balance** (*int*) — balance carried into the month.
* **closing\_balance** (*int*) — balance at month end.
* **entries** (*\[]Entry*) — list of transaction records (may be empty).

### 2.5 `Entry`

* **amount** (*int*) — signed value; positive = income, negative = expense. **Required**.
* **internal** (*bool, optional*) — `true` if the entry transfers money between two accounts in the same ledger. Defaults to `false`.
* **note** (*string*) — non‑empty description. **Required**.
* **date** (*string, optional*) — ISO‑8601 date (`YYYY‑MM‑DD`).
* **tag** (*string, optional*) — category label.

---

## 3 Validation Rules (invariants)

### Year‑level

0. **Y‑0** — Year key (`yearNum`) **must be a positive integer** (`yearNum > 0`).
1. **Y‑1** — Consecutive years must chain totals: `prev.closing_balance = next.opening_balance`.
2. **Y‑2** — A year’s `opening_balance` equals the first month’s `opening_balance`.
3. **Y‑3** — A year’s `closing_balance` equals the last month’s `closing_balance`.
4. **Y‑4** — A `Year` **must contain at least one `Month`** entry.

### Month‑level

0. **M‑0** — Month key (`monthNum`) **must be between 1 and 12** (inclusive).
1. **M‑1** — Consecutive months (including across years) must chain totals: `prev.closing_balance = next.opening_balance`.
2. **M‑2** — A month’s `opening_balance` equals the sum of all account `opening_balance` values.
3. **M‑3** — A month’s `closing_balance` equals the sum of all account `closing_balance` values.
4. **M‑4** — Within each month, Σ(`entry.amount` where `internal = true`) **must equal 0** (double‑entry constraint).
5. **M‑5** — A `Month` **must contain at least one `Account`** entry.

### Account‑level

1. **A‑1** — For every account: `opening_balance + Σ(entry.amount) = closing_balance`.
2. **A‑2** — If an account exists in consecutive months, `prev.closing_balance = next.opening_balance`.
3. **A‑3** — A new account must start with `opening_balance = 0`.
4. **A‑4** — An account may be omitted in later months only if its last `closing_balance = 0`.

### Entry‑level

1. **E‑1** — If an entry has a `date`, that date **must lie within the year and month of its parent `Month` object**.
2. **E‑2** — Every `Entry` **must include both `amount` and non‑empty `note` fields**.
3. **E‑3** — If `date` is present, **it must strictly follow the ISO‑8601 `YYYY‑MM‑DD` format**.

*A file that violates any invariant is non‑conforming.*

---

## 4 Example (YAML)

```yaml
years:
  2025:
    opening_balance: 1000      # USD, high‑level dollars
    closing_balance: 1050      # matches last month
    months:
      1:
        opening_balance: 1000
        closing_balance: 1050  # 620 + 430
        accounts:
          Checking:
            opening_balance: 600
            closing_balance: 620
            entries:
              - {amount:  50, internal: false, note: "Salary (January)", date: "2025-01-28", tag: Income}
              - {amount: -30, internal: true, note: "Transfer to Savings", date: "2025-01-30", tag: Transfer}
          Savings:
            opening_balance: 400
            closing_balance: 430
            entries:
              - {amount: 30, internal: true, note: "Transfer from Checking", date: "2025-01-30", tag: Transfer}
```

---

## 5. Versioning

* **Patch** — wording tweaks, no structural changes.
* **Minor** — additive, backward‑compatible fields or invariants.
* **Major** — changes that can invalidate previously valid files.

Add a `spec_version` field at the root when the first stable release (v1.0) is published.

---

## 6. License

Released under **CC‑BY‑4.0** — use, fork, extend; credit the authors.

---

*End of document.*
