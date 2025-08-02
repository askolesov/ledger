# Open Ledger Format — Personal Finance (v2.0)

> *A minimal, human‑readable schema for everyday finances that keeps every cent accounted for.*

---

## 1. Overview

The **Open Ledger Format (OLF)** stores personal‑finance data in plain‑text files—**YAML, JSON, or TOML**.

**Design principles**

1. **Human‑oriented** — anyone can maintain the file with a plain text editor.

2. **Consistency** — validation rules ensure strict, deterministic balances across the log.

3. **High‑level focus** — summary budgeting (≈ \$100–\$1 000 chunks), not coffee‑by‑coffee logging.

---

## 2 Data Model

```
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

* **amount** (*int*) — signed value; positive = income, negative = expense.
* **internal** (*bool, optional*) — `true` if the entry transfers money between two accounts in the same ledger. Defaults to `false`.
* **note** (*string*) — non‑empty description.
* **date** (*string, optional*) — ISO‑8601 date (`YYYY‑MM‑DD`).
* **tag** (*string, optional*) — category label.

---

## 3 Validation Rules (invariants)

### Year‑level

1. **Y‑1** — Consecutive years must chain totals: `prev.closing_balance = next.opening_balance`.
2. **Y‑2** — A year’s `opening_balance` equals the first month’s `opening_balance`, and its `closing_balance` equals the last month’s `closing_balance`.

### Month‑level

1. **M‑1** — Consecutive months (including across years) must chain totals: `prev.closing_balance = next.opening_balance`.
2. **M‑2** — A month’s `opening_balance` equals the sum of all account `opening_balance` values; its `closing_balance` equals the sum of all account `closing_balance` values.
3. **M‑3** — Within each month, Σ(`entry.amount` where `internal = true`) **must equal 0** (double‑entry constraint).

### Account‑level

1. **A‑1** — For every account: `opening_balance + Σ(entry.amount) = closing_balance`.
2. **A‑2** — If an account exists in consecutive months, `prev.closing_balance = next.opening_balance`.
3. **A‑3** — A new account must start with `opening_balance = 0`. An account may be omitted in later months only if its last `closing_balance = 0`.

### Entry‑level

1. **E‑1** — If an entry has a `date`, that date **must lie within the year and month of its parent `Month` object**.

*A file that violates any invariant is non‑conforming.*

---

## 4. Example (YAML)

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
