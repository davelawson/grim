# Grimoire

Powerful wizards compete with one another to be the first to Ascend in this turn-based game.

## Components

### API

The Go application is a simple RESTful api the can receive requests from any sort of front end such as a Discord Bot, a Web Page, or a CLI.

### Shared Go API definitions

The `api` package contains the HTTP request and response types, including nested
match and catalogue types. It depends only on the Go standard library. Backend
route registration lives in `internal/server`.

Each operation exposes a method and route constant, such as `api.GetMatchMethod`
and `api.GetMatchRoute`. Routes with parameters also expose a concrete path helper:

```go
req, err := http.NewRequestWithContext(ctx, api.GetMatchMethod,
    serverURL+api.GetMatchPath(matchID), nil)
```

Use route constants directly for endpoints without parameters, such as
`api.LoginRoute`. Parameterized route constants use server-side `:id` syntax;
clients use path helpers, which escape each ID as a URL path segment. The caller
supplies the server origin (without a trailing slash), authentication, and HTTP
transport configuration. Protected endpoints currently take the raw login token
in the `Authorization` header.

Separate Go clients import `github.com/davelawson/grim/api`. Grimoire continues
running as an HTTP server; the client owns HTTP transport and authentication.
The `api` package is the intended client-facing contract.

After the module rename is merged and the initial `v0.1.0` tag is published,
install it from the consuming repository:

```sh
go get github.com/davelawson/grim@v0.1.0
```

For local development before publication, point the consuming module at this
checkout instead:

```sh
go mod edit -replace=github.com/davelawson/grim=/absolute/path/to/grim
go get github.com/davelawson/grim/api
```

Remove the local replacement before switching to a published version:

```sh
go mod edit -dropreplace=github.com/davelawson/grim
go get github.com/davelawson/grim@v0.1.0
```

For a private repository, configure authenticated Git access and add
`github.com/davelawson/grim` to `GOPRIVATE` (preserving any existing entries).

For example, fetch a match using the caller's configured HTTP client and token:

```go
package example

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/davelawson/grim/api"
)

func GetMatch(ctx context.Context, client *http.Client, serverURL, token, matchID string) (*api.MatchView, error) {
    req, err := http.NewRequestWithContext(ctx, api.GetMatchMethod,
        serverURL+api.GetMatchPath(matchID), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", token)
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        var failure api.ErrorResponse
        if err := json.NewDecoder(resp.Body).Decode(&failure); err != nil || failure.Error == nil {
            return nil, fmt.Errorf("get match: HTTP %d", resp.StatusCode)
        }
        failure.Error.Status = resp.StatusCode // Status is excluded from JSON.
        return nil, failure.Error
    }
    var result api.MatchResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return result.Match, nil
}
```

For requests with bodies, marshal the corresponding API request type with
`json.Marshal`, pass the bytes through `bytes.NewReader`, and set
`Content-Type: application/json`. Obtain the token with `api.LoginRequest` and
`api.LoginResponse`. Use an HTTP client with a suitable timeout and TLS trust
configuration for the server.

### DB

State is stored in a simple sqlite3 database.

## Running

### Setup Database

1. Make sure that you have sqlite3 installed.
1. Set the GRIM_DB environment variable to point to where you want the sqlite3 db stored.
  Eg: `export GRIM_DB=$HOME/sqlite/test.db`
1. Create the database by executing the create-database.sql file.
  Eg: `sqlite3 $GRIM_DB < ./sql/create-database.sql`
1. Populate the database with some test data.
  `sqlite3 "$GRIM_DB" ".cd sql" ".read populate-test-data.sql"`

The creation script resets the database; it is for explicit fresh setup, not
an upgrade. This drafting release requires fresh setup with the updated schema,
including `playing` match status and command receipts. The server does not run
migrations or recreate existing databases. Previous match snapshot schemas and
unknown rules/catalogue versions return `409 unsupported_state`.

### Generate SSL key and cert

To listen to https instead of just http, we need to generate a cert and key and serve them up.
This needs to be done outside of the repo, since it obviously shouldn't be committed or shared, and also because it needs to be different for every machine that is running a server.

1. Set the GRIM_SSL environment variable to point to a folder on the server that will contain the crt and key files generated below.
  Eg: `export GRIM_SSL=$HOME/ssl`
1. Move into the GRIM_SSL directory and create the cert and key files.
    1. `cd $GRIM_SSL`
    1. `openssl genrsa -out grim.key 2048`
    1. `openssl ecparam -genkey -name secp384r1 -out grim.key`
    1. `openssl req -new -x509 -sha256 -key grim.key -out grim.crt -days 3650`
        - Just leave everything blank (default)
    1. verify that the folder contains 2 files:
        - `grim.key`
        - `grim.crt`

### Run Service

1. Run the app using go.
  Eg: `go run main.go`
1. Issue requests against the api found at `http://localhost:8080`.

### Testing

Run `go test ./...` and `go vet ./...`. Lifecycle tests use temporary SQLite
databases and production routes; they cover authorization, readiness, atomic
launch, private drafting, atomic game start, concurrent requests, safe retries,
restart recovery, rollback, and soft deletion. Engine tests cover all opening
offers, deck composition, and reproducible random seating and deals.
Run `go test -race ./...` to check concurrent access. Tests do not use `GRIM_DB`.

#### cURL

Since everything is behind a simple web RESTful web server, we can use cURL to exercise endpoints.  Many endpoints require an authentication token, which can be obtained by calling the /login endpoint.
The web server is currently using a self-signed certificate, so curl will refuse unless the `-k` argument is provided.
Here is an example of a command to query a user:
`
curl -X 'POST' \
  'https://localhost:8080/user/getbyemail' \
  -H 'accept: application/json' \
  -H 'Authorization: wymKUm3TEyDb+UBtuYS31fEP/B+fup6zK2KcQrVY3ls=' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "tim@aol.com"
}'`

All API JSON fields use lowercase names. Login returns `{"token":"…"}`.
All protected endpoints return authentication errors as
`{"error":{"code":"…","message":"…"}}`: missing or invalid tokens return
401 with `unauthorized`; unexpected authentication failures return 500 with
`internal_error` and a generic message.
Lobby requests use `name`, `owner`, and `userid` as applicable; lobby lookup
returns a `lobby` object with `id`, `name`, `owner`, `members`, `status`, nullable
`matchid`, and `readiness` keyed by member user ID.
Deleting an open lobby sets its `deleted_at` timestamp and preserves its
memberships. Deleted lobbies return 404 for lookup and further operations,
including repeated deletion. Only the owner can delete an open lobby.
User lookup returns
`{"user":{"id":"…","name":"…","email":"…"}}`. Registration at `POST /user`
uses `email`, `name`, and `password`; capitalized or mixed-case versions of
those request keys return 400. Other requests continue accepting field names
case-insensitively.

### Launch and end a match

The owner adds two to four members (including themselves) using the existing
lobby membership endpoints. Each member then sets their own readiness:

```sh
curl -k -X PUT 'https://localhost:8080/lobby/<lobby-id>/ready' \
  -H 'Authorization: <member-token>' -H 'Content-Type: application/json' \
  -d '{"ready":true}'
```

Changing the roster clears all readiness. Set `ready` to false to withdraw your
own agreement. An ordinary member can leave with
`DELETE /lobby/<lobby-id>/user/<their-user-id>`; the owner transfers ownership
to another member through the existing lobby update before leaving.

Once everyone is ready, the owner launches with a client-generated UUID request
ID. Keep that ID to retry the same request if its response is lost:

```sh
curl -k -X POST 'https://localhost:8080/lobby/<lobby-id>/launch' \
  -H 'Authorization: <owner-token>' -H 'Content-Type: application/json' \
  -d '{"requestid":"72e865f9-2d58-465e-b3bb-708411adfb22"}'
```

Launch returns 201 with `{"match":{...}}`, containing `id`, `name`, `lobbyid`,
`status`, `revision`, `participants`, `ring`, `outcome`, `rulesversion`,
`catalogueversion`, `wizards`, `turn`, `draftoffers`, and the owner's `you` view.
The match starts in setup at revision 0 with an empty Ring, null turn and
outcome, and pinned `opening-v1` rules/catalogue. The lobby closes permanently
and points to the match. Identical retries
replay the original result; a new request ID on the closed lobby returns 409.
Once the match is deleted, retries using its accepted request ID return 404.
Authenticated participants fetch current state with `GET /match/<match-id>`.

### Draft starting cards

Participants fetch the pinned opening definitions and every draft offer with
`GET /match/<match-id>/catalogue`, which returns `{"catalogue":{...}}`. Cards
have stable design IDs and structured costs, requirements, Expertise, Affinity,
income, persistence, disposal, resource values, and typed effect identifiers.
The catalogue contains 31 opening designs, including four clearly marked inert
victory placeholders; it does not include rumours or execute card effects.

Each participant makes three irreversible choices independently, in order:
`choose_chantry`, `choose_victory_path`, then `choose_minion`. Each command uses
the same endpoint and the participant's own authentication token:

```sh
curl -k -X POST 'https://localhost:8080/match/<match-id>/commands' \
  -H 'Authorization: <participant-token>' -H 'Content-Type: application/json' \
  -d '{"requestid":"9493d7de-393b-49e2-a13b-bc41dc33a3de","expectedrevision":0,"type":"choose_chantry","cardid":"chantry_fire"}'
```

Use a new UUID and the latest match revision for each new choice. For example,
after choosing `chantry_fire`, select `subjugate_rival` for Domination (or one
of the offered placeholders), then a Fire offer such as `bastian_redhand`.
All five victory paths are available with every Chantry. Choosing Domination
does not restrict which offered Minion can be selected. Drafting has no payments.

Commands return `200 {"match":{...}}`. The required nonnegative integer
`expectedrevision` refers to the whole match, so another participant's choice
may cause `409 stale_revision`; refresh before submitting at the new revision.
Retry a lost response using exactly the same request ID and input: it replays
the original participant view even if the match has since advanced. Changed
input under an accepted request ID returns `409 request_id_conflict`.
Malformed or unknown command fields return `400 invalid_request`; a card not
offered for that step returns `400 invalid_choice`; attempting the wrong step
or changing a locked choice returns `409 draft_sequence`. New choices after
game start return `409 draft_closed`. Rejected requests do not change state or
create receipts. UUID casing, whitespace, and JSON key order do not affect retries.

The `you.draft` object shows only the caller's step and selected design IDs.
Other wizards expose only `draftcomplete` while setup is pending. Choices have
unlimited copies, no deadline, and no automatic picks. The final Minion choice
atomically creates a random Ring, gives each wizard 3 Integrity and a Chantry
in play, independently shuffles their ten-card deck, and deals five cards.
It enters `playing` with `turn` identifying Ring's first wizard, number 1,
phase `start`, and step `recovery`. No turn step, income, or rumour has executed;
turn execution is a later slice.

After start, `wizards` contains public Integrity, Affinities, victory progress,
in-play cards, ordered discards, and hand/draw-pile counts. `you.hand` shows the
caller's cards; `you.drawpile` lists their remaining cards sorted by instance
ID, never in draw order. Other draft selections remain private. Card instances
identify their design, owner, location, and exhaustion state separately from
the shared definition. A deleted match returns 404 for lookup, catalogue,
and command submission, including accepted-command retries.

### Administrative termination

The operator grants administrator permission directly in the database:

```sh
sqlite3 "$GRIM_DB" "UPDATE users SET admin = 1 WHERE id = '<admin-user-id>';"
```

An authenticated administrator can terminate any existing match, including
finished matches and matches with unsupported saved state:

```sh
curl -k -X POST 'https://localhost:8080/admin/match/<match-id>/end' \
  -H 'Authorization: <admin-token>'
```

Termination returns 204 with no body and sets the match's `deleted_at` to the
current UTC timestamp. It preserves status, outcome, revision, snapshot, players,
and launch receipts. Match lookup then returns 404. Repeating termination returns
204 without changing the original timestamp; an unknown match returns 404.
The linked closed lobby remains visible with its match link and cannot be
reopened or changed. Admins need no match membership to terminate a match;
lookup remains participant-only. Revoking the stored flag blocks the next admin
request. Player abandonment voting is deferred. New endpoint errors use
`{"error":{"code":"…","message":"…"}}`.

### Swagger

End points are documented using a swagger interface.  This also allows manual testing of the end points.
Swaggo is the implementation we're using.  It derives meaning from comments in the code, which isn't great, but we'll stick with it a while longer.

To setup Swagger:

1. Install swagger: `go install github.com/swaggo/swag/cmd/swag@latest`
1. Make sure that `$(go env GOPATH)/bin` is on your path
1. Generate the swagger documentation: `swag init .`
1. Run the application (detailed in earlier section)
1. Open the web page: `http://localhost:8080/swagger/index.html#`
