# dbmodgen

[![License](https://img.shields.io/github/license/kukymbr/dbmodgen.svg)](https://github.com/kukymbr/dbmodgen/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/dbmodgen.svg)](https://github.com/kukymbr/dbmodgen/releases/latest)

The `dbmodgen` generates row-based models from the existing database structure.

Based on the code of the [genna](https://github.com/dizzyfool/genna), 
edited to a more simplified version, generating models without the `go-pg` (nor other ORM) direct integration.
Just the `db` struct tag (or any other you like) and freedom of the further usage.

Current RDBMS support:

| PostgreSQL | Microsoft SQL Server | MySQL | Oracle Database | SQLite | IBM Db2 |
|------------|----------------------|-------|-----------------|--------|---------|
| ✅          | ➖                    | ➖     | ➖               | ➖      | ➖       |

Yep, only Postgres, sorry folks.

## Example

Assume we have this `users` table in the database:

```sql
CREATE TABLE public.users
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    email      VARCHAR(128)     NOT NULL,
    name       TEXT                      DEFAULT NULL,
    password   VARCHAR(128)     NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);
```

With the default config, the `dbmodgen` will generate the next model:

```go
package dbmodel

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      *string   `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

```
## Installation

The go 1.24 is a minimal requirement for the `dbmodgen`, so the `go tool` is a preferred way to install:

```shell
go get -tool github.com/kukymbr/dbmodgen/cmd/dbmodgen@latest
```

## Usage

The `dbmodgen --help` output:

```text
Generates structured SQL getters in go code.
See https://github.com/kukymbr/dbmodgen for info.

Usage:
  dbmodgen [flags]
```

## License

[MIT](LICENSE).