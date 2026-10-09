# Grimoire

Powerful wizards compete with one another to be the first to Ascend in this turn-based game.

## Components

### API

The Go application is a simple RESTful api the can receive requests from any sort of front end such as a Discord Bot, a Web Page, or a CLI.

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
an upgrade. This lifecycle release needs the updated schema. The server does
not run migrations or recreate existing databases. The schema includes nullable
`deleted_at` timestamps on matches and lobbies; existing databases require fresh
setup to use this release.

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
launch, concurrent requests, safe retries, restart recovery, and soft deletion.
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
`status`, `revision`, `participants`, `ring`, and `outcome`. The match starts in
setup at revision 0 with null outcome; drafting and gameplay come in a later
slice. The lobby closes permanently and points to the match. Identical retries
replay the original result; a new request ID on the closed lobby returns 409.
Once the match is deleted, retries using its accepted request ID return 404.
Authenticated participants fetch current state with `GET /match/<match-id>`.

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
