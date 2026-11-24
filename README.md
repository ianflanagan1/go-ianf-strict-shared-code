# go-ianf-strict-shared-code

The code shared by every project on [go-ianf-strict-baseline](../go-ianf-strict-baseline): adapters and value code. It holds three sections:

- **Driven adapters**, under `adapters/driven/`, which may use third-party modules such as pgx. Each package implements a port by doing I/O and makes no decisions, so it can be reviewed once and trusted everywhere. Import them only from a composition root (`cmd/`) or a project's own adapters; the shell declares the ports they satisfy.
- **Driving adapters**, under `adapters/driving/`. `adapters/driving/timer`, the system timer, is the only one: it calls the function it is given each time a duration passes, so a program can run a use case on a schedule. Start it only from a composition root, with a function that calls an orchestrator.
- **Values**, under `values/`. `values/id` is the shared ID unit-of-code: the row ID (a UUIDv7) that every owned table's ID type wraps, the public ID (a UUIDv4) that every public ID type wraps, and the `Internal` wrapper, whose printed, logged and serialised forms are a placeholder. Its factory functions only check and hold an ID they're given; new IDs come from the `ids` adapter. `values/email` is the shared email address unit-of-code, described under [Email addresses](#email-addresses). A project's values section may import either.

| Package | Adapts | Satisfies a port like |
|---|---|---|
| `clock` | The system clock | `Now() time.Time` |
| `stopwatch` | The system's monotonic clock, measuring how long has passed and never what time it is; `Fake` reports durations a test gives it | `Start() func() time.Duration` |
| `timezones` | The machine's local zone and the time zone database; `Fake` offers zones a test gives it | `Local() *time.Location`, `Load(name string) (*time.Location, error)` |
| `deadlines` | The system's timers, through `context.WithTimeout`; `Fake` hands out contexts whose deadlines have passed or never pass, as a test says | `After(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc)` |
| `env` | The process environment | `Lookup(key string) (string, bool)` |
| `random` | `crypto/rand` (`Secure`) | `io.Reader` |
| `ids` | The system clock and `random.Secure` (`System`), laying out each new row ID and public ID; `Fake` hands out IDs a test gives it | `NewRowID() id.RowID`, `NewPublicID() id.PublicID` |
| `files` | A directory, through `os.Root`, so paths can't escape it | `ReadFile(name string) ([]byte, error)`, `WriteFile(name string, content []byte) error` |
| `httpsend` | HTTP, through `net/http` | `Send(ctx context.Context, req httpmsg.Request) (httpmsg.Response, error)` |
| `sqlquery` | Postgres, through `database/sql` and pgx | `Query(ctx, sqlmsg.Query) (sqlmsg.Rows, error)`, `Exec(ctx, sqlmsg.Query) (sqlmsg.Result, error)`, `InTx(ctx, func(sqlmsg.Tx) error) error` |
| `timer` (driving) | The system's timers, calling a function each time a duration passes until its context is cancelled | Called from a composition root as `Run(ctx context.Context, every time.Duration, tick func(ctx context.Context))` |

The request, response, query and row types are in [go-ianf-strict-ports](https://github.com/ianflanagan1/go-ianf-strict-ports), a library of only core, so a project's core can build and read them and its shell can declare ports in them.

`sqlquery` supports Postgres only; its package documentation lists how each Postgres type arrives. `InTx` runs a function inside a transaction: it commits when the function returns nil, and rolls back when it returns an error or panics.

## Email addresses

`email.ParseAddress` turns text into an address, or says why it holds none:

```go
import "github.com/ianflanagan1/go-ianf-strict-shared-code/values/email"

address, err := email.ParseAddress("Ian.Flanagan@Example.COM")
var failure email.AddressError
if errors.As(err, &failure) {
	// email.ErrMissingAt, email.ErrLocalPartInvalid, email.ErrDomainNotFullyQualified, ...
}

address.LocalPart // "Ian.Flanagan"
address.Domain    // "example.com"
address.Text()    // "Ian.Flanagan@example.com"
```

It accepts the addresses people type into a form, in any script:

- The local part is atoms joined by single dots, of at most 64 bytes. An atom holds letters, digits and ``!#$%&'*+/=?^_`{|}~-``, and any character beyond ASCII but those that show nothing, show as a space or change the direction of the text around them: `用户` is a local part, and one holding a zero-width space is not.
- The domain is a host name of two labels or more, of letters, digits and hyphens, whose last label isn't all digits. A label in another script is accepted and held in its ASCII form: `münchen.de` is held as `xn--mnchen-3ya.de`.
- The address is at most 254 bytes.

It rejects what the standards allow and people don't use: a quoted local part (`"ian"@example.com`), a comment, an IP address for a domain (`ian@[192.168.0.1]`), a domain of one label (`ian@localhost`) and a dot that ends the domain. It rejects text with anything around the address too, a space or a newline included, so a caller that reads a form field trims it first.

`ParseAddress` normalises, so that two spellings of one address are equal values: the local part to NFC, and the domain to lower case and to its ASCII form. It keeps the local part's case, which the standards leave to the mail server: `Ian@example.com` and `ian@example.com` are two values. Whether they are one mailbox, and what a `+tag` in a local part means, are questions for a project's decision code, each a function with its own name.

The nine failures are the constants of `email.AddressError`, an error type: find it in the error with `errors.As`, and branch on it with a `switch` that names every constant, so that the `exhaustive` linter reports a failure added later. Its message says why and holds nothing of the text.

An address with a local part beyond ASCII reaches only a mail server that supports SMTPUTF8.

The domain's conversion (UTS #46, with only letters, digits and hyphens in a label) comes from `golang.org/x/net/idna`, and the local part's NFC from `golang.org/x/text/unicode/norm`. Both hold their Unicode tables in their own code, so the tables are the ones `go.mod` and the Go version fix (Unicode 17 from Go 1.27) and not those of the machine. The type is a port of php-ianf-strict-shared-code's `EmailAddress`, and parses every text that class's tests name to the same local part and domain, or fails it for the same reason.

## Testing

`devbox run -- task test` runs everything that needs nothing outside the process. `sqlquery`'s tests need Postgres:

```sh
docker run --rm -d -e POSTGRES_PASSWORD=test -p 55432:5432 postgres:17
GO_ADAPTERS_POSTGRES_URL=postgres://postgres:test@localhost:55432/postgres devbox run -- task integration
```

Add an adapter here when more than one project needs it. Keep it thin: when it needs a decision, the decision belongs in the caller's core.
