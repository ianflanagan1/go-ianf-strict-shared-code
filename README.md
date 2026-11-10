# go-adapters

Adapters shared by every project on [go-baseline](../go-baseline), and the values they share. It holds two sections:

- **Driven adapters**, under `adapters/driven/`, which may use third-party modules such as pgx. Each package implements a port by doing I/O and makes no decisions, so it can be reviewed once and trusted everywhere. Import them only from a composition root (`cmd/`) or a project's own adapters; the shell declares the ports they satisfy.
- **Values**, under `values/`. `values/id` is the shared ID unit-of-code: the row ID (a UUIDv7) that every owned table's ID type wraps, the public ID (a UUIDv4) that every public ID type wraps, and the `Internal` wrapper, whose printed, logged and serialised forms are a placeholder. Its factory functions only check and hold an ID they're given; new IDs come from the `ids` adapter. A project's values section may import it.

| Package | Adapts | Satisfies a port like |
|---|---|---|
| `clock` | The system clock | `Now() time.Time` |
| `env` | The process environment | `Lookup(key string) (string, bool)` |
| `random` | `crypto/rand` (`Secure`) | `io.Reader` |
| `ids` | The system clock and `random.Secure` (`System`), laying out each new row ID and public ID; `Fake` hands out IDs a test gives it | `NewRowID() id.RowID`, `NewPublicID() id.PublicID` |
| `files` | A directory, through `os.Root`, so paths can't escape it | `ReadFile(name string) ([]byte, error)`, `WriteFile(name string, content []byte) error` |
| `httpsend` | HTTP, through `net/http` | `Send(ctx context.Context, req httpmsg.Request) (httpmsg.Response, error)` |
| `sqlquery` | Postgres, through `database/sql` and pgx | `Query(ctx, sqlmsg.Query) (sqlmsg.Rows, error)`, `Exec(ctx, sqlmsg.Query) (sqlmsg.Result, error)`, `InTx(ctx, func(sqlmsg.Tx) error) error` |

The request, response, query and row types are in [go-ports](https://github.com/ianflanagan1/go-ports), a library of only core, so a project's core can build and read them and its shell can declare ports in them.

`sqlquery` supports Postgres only; its package documentation lists how each Postgres type arrives. `InTx` runs a function inside a transaction: it commits when the function returns nil, and rolls back when it returns an error or panics.

## Testing

`devbox run -- task test` runs everything that needs nothing outside the process. `sqlquery`'s tests need Postgres:

```sh
docker run --rm -d -e POSTGRES_PASSWORD=test -p 55432:5432 postgres:17
GO_ADAPTERS_POSTGRES_URL=postgres://postgres:test@localhost:55432/postgres devbox run -- task integration
```

Add an adapter here when more than one project needs it. Keep it thin: when it needs a decision, the decision belongs in the caller's core.
