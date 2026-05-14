# DB Migration Pattern

### `migration-db` — Database Schema Migration

Migration DB is a separate repo for database schema files. It stores ordered SQL migrations by environment.

**Key property**: Migration schema must match `internal/domain` table names and `gorm:"column:..."` tags. Domain code and migration SQL must describe the same table.

#### Package structure

```
<service>-migration-db/
  <env>/
    schema/
      {version}_{action}_{table_name}.up.sql
      {version}_{action}_{table_name}.down.sql
    data/
      {version}_{action}_{table_name}.up.sql
      {version}_{action}_{table_name}.down.sql
```

Example:
```
transfer-migration-db/
  sit/
    schema/
      1_create_table_verify_transfer_transactions.up.sql
      1_create_table_verify_transfer_transactions.down.sql
      2_create_table_transfer_transactions.up.sql
      2_create_table_transfer_transactions.down.sql
```

- `schema/` → table/schema manipulation, e.g. create table, alter table, create index
- `data/` → data manipulation, e.g. insert seed data, update reference data
- `{version}` → migration order/version, e.g. `1`, `2`, `3`
- `{action}` → action to do, e.g. `create_table`, `add_column`, `create_index`, `insert`
- `{table_name}` → target table name, snake_case, matches `TableName()`
- `.up.sql` / `.down.sql` → run direction

Rule:
- Use one version for one action
- Every `up.sql` must have a matching `down.sql`
- `down.sql` reverses only its matching `up.sql`

Examples:
```text
1_create_table_verify_transfer_transactions.up.sql
1_create_table_verify_transfer_transactions.down.sql
2_insert_transfer_type.up.sql
2_insert_transfer_type.down.sql
```

---

#### Example: create table

Domain model:
```go
type VerifyTransferTransactions struct {
    Id          string                   `gorm:"column:id"`
    Type        transfertype.Type        `gorm:"column:type"`
    FromAccount string                   `gorm:"column:from_account"`
    ToAccount   string                   `gorm:"column:to_account"`
    Amount      float64                  `gorm:"column:amount"`
    Status      transactionstatus.Status `gorm:"column:status"`
}

func (t VerifyTransferTransactions) TableName() string {
    return "verify_transfer_transactions"
}
```

Migration:
```sql
-- 1_create_table_verify_transfer_transactions.up.sql
CREATE TABLE verify_transfer_transactions (
    id              VARCHAR(36)     PRIMARY KEY,
    type            VARCHAR(10)     NOT NULL,
    from_account    VARCHAR(20)     NOT NULL,
    to_account      VARCHAR(20)     NOT NULL,
    amount          DECIMAL(17,2)   NOT NULL,
    status          VARCHAR(10)     NOT NULL
);
```

Rollback:
```sql
-- 1_create_table_verify_transfer_transactions.down.sql
DROP TABLE verify_transfer_transactions;
```

---

#### Type mapping

Common mapping from domain to SQL:

```
string id / UUID string        → VARCHAR(36)
account number string          → VARCHAR(20)
enum string                    → VARCHAR(10) or size that fits values
money amount                   → DECIMAL(17,2)
optional timestamp string/time → TIMESTAMP WITHOUT TIME ZONE nullable
```

Example with optional timestamp:
```sql
CREATE TABLE transfer_transactions (
    id                  VARCHAR(36)     PRIMARY KEY,
    transaction_ref_id  VARCHAR(36)     NOT NULL,
    type                VARCHAR(10)     NOT NULL,
    from_account        VARCHAR(20)     NOT NULL,
    to_account          VARCHAR(20)     NOT NULL,
    amount              DECIMAL(17,2)   NOT NULL,
    status              VARCHAR(10)     NOT NULL,
    transaction_time    TIMESTAMP WITHOUT TIME ZONE
);
```

**Key points:**
- Table name must match domain `TableName()`
- Column names must match `gorm:"column:..."`
- Enums are stored as `VARCHAR`, not database enum types
- `down.sql` should reverse only what the matching `up.sql` created
- Keep migration order explicit with number prefix
