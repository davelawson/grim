# Grimoire

Powerful wizards compete with one another to be the first to Ascend in this turn-based game.

## Components

### API

The Go application is a simple RESTful api the can receive requests from any sort of front end such as a Discord Bot, a Web Page, or a CLI.

### Shared Go API definitions

The `api` package contains the HTTP request and response types, including nested
match and catalogue types. It depends only on the Go standard library. Backend
route registration lives in `internal/server`.

Each operation exposes an `Endpoint` value with `Method` and `Route` fields,
such as `api.GetMatch.Method` and `api.GetMatch.Route`. Routes with parameters also
expose a concrete path helper:

```go
req, err := http.NewRequestWithContext(ctx, api.GetMatch.Method,
    serverURL+api.GetMatchPath(matchID), nil)
```

Use the `Route` field directly for endpoints without parameters, such as
`api.Login.Route`. Parameterized routes use server-side `:id` syntax;
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
    req, err := http.NewRequestWithContext(ctx, api.GetMatch.Method,
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
1. Create a fresh database with the reset script.
  `sqlite3 -bail "$GRIM_DB" < ./sql/create-database.sql`
1. Install the shared card catalogue.
  `sqlite3 -bail "$GRIM_DB" < ./sql/install-card-designs.sql`
1. Populate the database with some test data.
  `sqlite3 "$GRIM_DB" ".cd sql" ".read populate-test-data.sql"`

`new-database.sh` runs all three steps, including the optional test accounts.
The creation script resets the database and must only be used for fresh setup.
The server does not automatically install, migrate, seed, or reset data.

To upgrade an existing database, stop the backend and run only:

```sh
sqlite3 -bail "$GRIM_DB" < ./sql/install-card-designs.sql
```

This additive, transactional script creates `card_designs` and seeds all 98
original designs without changing accounts, lobbies, matches, or receipts.
Re-running it leaves existing catalogue rows untouched. Existing schema-3
`playable-v1` matches continue using the same definitions and design IDs;
unsupported older snapshots still require an explicit compatibility solution.

### Shared card catalogue

`card_designs` contains the entire catalogue, including designs absent from any
player's deck. One row represents a design within a catalogue version, keyed by
`(catalogue_version, design_id)`. Individual card instances retain their
`designid`, ownership, location, and changing state in match snapshots; those
properties do not belong in this shared table.

Names, text, timing, duration, persistence, disposal, and other ordinary properties
have columns. Stats have individual columns, for example `cost_wis`,
`play_requirement_fire`, `expertise_might`, and `income_ephemeral_wis`.
`NULL` means an absent stat; `0` is an explicit zero. Categories, effect IDs,
abilities, and Challenges use JSON arrays. Challenge reward order and repeated
rewards are preserved. `sort_order` runs consecutively from zero and preserves
catalogue order used by reproducible Research and rumour draws.

Manage catalogue changes with SQL scripts. A version becomes immutable when its
first match is saved: database triggers reject row additions, changes, and
removals for versions referenced by any match, including ended or soft-deleted
matches. Keep historical versions available. To change content, insert a complete
new version in a transaction, keeping the existing design IDs where appropriate,
then edit its columns before any match uses it. A stat change in an unused version
can be made with:

```sql
UPDATE card_designs SET cost_wis = 2
WHERE catalogue_version = 'edition-two' AND design_id = 'study';
```

Select that installed version for new matches and restart the backend:

```sh
export GRIM_CATALOGUE_VERSION=edition-two
```

The default is `playable-v1`. Catalogue versions are independent of rules versions:
a data-only edition still uses `rules_version='playable-v1'`. New game behavior
requires supported Go handlers and an engine release; adding a row is not a way
to introduce arbitrary executable effects. Existing matches and draft receipts
continue resolving their own pinned version. Public catalogue/draft responses
and gameplay endpoints are unchanged, so the CLI needs no protocol changes.

Startup validates the configured version and refuses to serve if it is missing
or invalid. Each match operation loads its pinned definitions through the current
database transaction; there is no process-global catalogue or built-in fallback.
The engine receives immutable definitions and preserves them through state
cloning without embedding them in saved snapshots. Accepted gameplay receipts
replay before current catalogue or state validation.

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

Create a public lobby with `POST /lobby` and `{"name":"Friday game"}`. It
returns `200 {"id":"…"}` and makes the authenticated creator its owner and
first member, initially unready. All lobbies are public and hold at most four
members; each user may belong to only one nondeleted open lobby.

Browse available lobbies using authenticated `GET /lobbies/public`:

```json
{"lobbies":[{"id":"…","name":"Friday game","owner":"…","membercount":1}]}
```

Discovery includes only open, nondeleted lobbies with fewer than four members,
ordered by ID, with no pagination. Empty results are `{"lobbies":[]}`. Users
who already belong to an open lobby can still browse. Full lobby details,
including roster and readiness, remain member-only.

Join as the authenticated user with no request body:

```sh
curl -k -X POST 'https://localhost:8080/lobby/<lobby-id>/join' \
  -H 'Authorization: <member-token>'
```

Joining returns `200 {"lobby":{...}}` with the current roster and readiness.
An identical join to an open lobby already joined returns its current view,
even when full, without clearing readiness. A new member joining clears all
members' readiness. Owners can also add users through the existing membership
endpoint, subject to the same capacity and open-membership limits.

Creation, joining, and owner additions enforce the one-open-lobby rule.
Conflicts return `409` with `lobby_full` or `already_in_open_lobby` in the
standard JSON error envelope. Joining a closed lobby returns `409 lobby_closed`;
a missing/deleted lobby returns `404 not_found`. Requests require the existing
raw login token in `Authorization`. Membership checks and updates are atomic,
including competing joins. Leave your current open lobby before joining or
creating another; owners must transfer ownership or delete it before leaving.
Launch or deletion releases members to enter another lobby while preserving
historical memberships. Existing conflicting legacy memberships are not
rewritten automatically. No database migration is required for this feature.

Once at least two members have joined, each member sets their own readiness:

```sh
curl -k -X PUT 'https://localhost:8080/lobby/<lobby-id>/ready' \
  -H 'Authorization: <member-token>' -H 'Content-Type: application/json' \
  -d '{"ready":true}'
```

Changing the roster clears all readiness. Set `ready` to false to withdraw your
own agreement. An ordinary member can leave with
`DELETE /lobby/<lobby-id>/leave` with no body. The token identifies the departing
user. Success returns 200 with no body; errors use the standard JSON envelope.
The owner must transfer ownership to another member through the existing lobby
update or delete the open lobby before leaving (`409 owner_must_transfer`).
Closed lobbies return `409 lobby_closed`; missing/deleted lobbies return
`404 not_found`. The existing `DELETE /lobby/<lobby-id>/user/<user-id>` remains
available for owner-managed removal and explicit self-departure with its existing
authorization. Both routes retain the lobby ID explicitly.

The sibling `grimcli` repository supports the full flow through `public-lobbies`,
`lobby`, `lobby-create`, `lobby-join`, `lobby-ready`, `lobby-launch`, `lobby-leave`,
`lobby-remove-user`, and `lobby-delete`. Its interactive Bash wrapper keeps lobby
and match selections separately. Launch commands use the lobby's UUID as
`requestid`, so retrying requires no separately saved request ID.

Once everyone is ready, the owner launches with a client-generated UUID request
ID. Keep that ID to retry the same request if its response is lost:

```sh
curl -k -X POST 'https://localhost:8080/lobby/<lobby-id>/launch' \
  -H 'Authorization: <owner-token>' -H 'Content-Type: application/json' \
  -d '{"requestid":"72e865f9-2d58-465e-b3bb-708411adfb22"}'
```

Launch returns 201 with `{"match":{...}}`, containing `id`, `name`, `lobbyid`,
`status`, `revision`, `participants`, `ring`, `outcome`, `rulesversion`,
`catalogueversion`, `wizards`, `turn`, and the owner's `you` view. Match responses
do not contain draft options.
The match starts in setup at revision 0 with an empty Ring, null turn and
outcome, `playable-v1` rules, and the configured catalogue version. The lobby closes permanently
and points to the match. Identical retries
replay the original result; a new request ID on the closed lobby returns 409.
Once the match is deleted, retries using its accepted request ID return 404.
Authenticated participants fetch current state with `GET /match/<match-id>`.

### Draft starting cards

Participants fetch their pinned card definitions with
`GET /match/<match-id>/catalogue`, which returns `{"catalogue":{...}}`. Cards
have stable design IDs and structured costs, requirements, Expertise, Affinity,
income, persistence, disposal, resource values, and typed effect identifiers.
The initial `playable-v1` catalogue contains 98 definitions: the original 31 designs, two ephemeral
resources, 32 Spells, 20 rumours, and 13 reward designs. Research Artifact and
the four unfinished victory paths are explicit inert placeholders. The
catalogue includes structured Challenge requirements, upkeep, timing,
durations, Research pools, and abilities; fetching it does not execute effects. Definitions remain available before and during drafting.

Begin or resume your draft with your own authentication token:

```sh
curl -k 'https://localhost:8080/match/<match-id>/draft' \
  -H 'Authorization: <participant-token>'
```

This read-only request returns `200 {"draft":{"step":"chantry"},"cards":[...],"revision":0,"status":"setup"}`
for a new draft. `cards` contains full card definitions for the current step
only: first Chantries, then victory-path cards, then the chosen Chantry's
Starting Minions. `draft` also contains your locked selections when resuming.
Fetching options does not change the match revision or create a receipt.

Each participant makes three irreversible choices independently, in order:
`choose_chantry`, `choose_victory_path`, then `choose_minion`. Submit each choice
to the same draft path with the participant's own authentication token:

```sh
curl -k -X POST 'https://localhost:8080/match/<match-id>/draft' \
  -H 'Authorization: <participant-token>' -H 'Content-Type: application/json' \
  -d '{"requestid":"9493d7de-393b-49e2-a13b-bc41dc33a3de","expectedrevision":0,"type":"choose_chantry","cardid":"chantry_fire"}'
```

Use a new UUID and the latest match revision for each new choice. For example,
after choosing `chantry_fire`, select `subjugate_rival` for Domination, then a Fire offer such as `bastian_redhand`.
Domination is the only offered victory path for every Chantry. It does not restrict which offered Minion can be selected. Drafting has no payments.

Selections use `api.DraftPickRequest` and return `200` with `api.DraftResponse`:
`{"draft":{...},"cards":[...],"revision":1,"status":"setup"}`. Each accepted
choice advances your draft and returns the next step's options. Ordinary gameplay uses `/match/<match-id>/actions/<operation>`, required choices use
`/match/<match-id>/choices/<choice-id>/<operation>`, and Priority responses use
`/match/<match-id>/reactions/<operation>`. Draft choices remain
on the dedicated draft endpoint.

The required nonnegative integer
`expectedrevision` refers to the whole match, so another participant's choice
may cause `409 stale_revision`; refresh before submitting at the new revision.
Retry a lost response using exactly the same request ID and input: it replays
the original draft progress and options even if the match has since advanced. Changed
input under an accepted request ID returns `409 request_id_conflict`.
Malformed or unknown command fields return `400 invalid_request`; a card not
offered for that step returns `400 invalid_choice`; attempting the wrong step
or changing a locked choice returns `409 draft_sequence`. New choices after
game start return `409 draft_closed`. Rejected requests do not change state or
create receipts. UUID casing, whitespace, and JSON key order do not affect retries.

The `you.draft` object shows only the caller's step and selected design IDs.
Other wizards expose only `draftcomplete` while setup is pending. A completed
player's draft endpoint returns `step:"complete"` and `cards:[]`, with status
`setup` while waiting for others or `playing` once the game starts. Choices have
unlimited copies, no deadline, and no automatic picks. The final Minion choice
atomically creates a random Ring, gives each wizard 3 Integrity and a Chantry
in play, independently shuffles their ten-card deck, and deals five cards.
It enters `playing` and advances Ring's first wizard through automatic Start
steps to the first input boundary. First-turn Income and the normal rumour
draw occur in the final draft transaction; a response window or optional
Start effect may therefore be pending before Actions.

After start, `wizards` contains public Integrity, Affinities, victory progress,
in-play cards, ordered discards, and hand/draw-pile counts. The optional
`you.userid` identifies the authenticated participant, including when their hand
is empty or they are eliminated. Compare it with `turn.activewizard`,
`exchange.priority`, and `choice.actor` to identify the next actor. This is a
projection-only addition; snapshots and rules versions do not change. Historical
accepted responses may lack identity, so fetch a fresh match before composing
new input. `you.hand` shows the caller's cards; `you.drawpile` lists their remaining cards sorted by instance
ID, never in draw order. Other draft selections remain private. Card instances
identify their design, owner, location, and exhaustion state separately from
the shared definition. A deleted match returns 404 for lookup, catalogue,
and both draft methods, including accepted-choice retries.

This gameplay release needs no database reset or table migration. New matches
use snapshot schema 3, `playable-v1` rules, and the configured catalogue; opening-v1 snapshots
are unsupported and must be replaced. Stored accepted responses retain the
existing retry rules; they do not upgrade or resume unsupported matches.

### Play through turns and victory

Poll `GET /match/<match-id>` with the participant's token. The server advances
only during accepted submissions, stopping for input. Check `status` first, then:

1. When `choice` is present, only `choice.actor` may answer at
   `/match/<match-id>/choices/<choice-id>/<operation>`.
2. Otherwise, when `exchange` is present, only `exchange.priority` may respond at
   `/match/<match-id>/reactions/<operation>`.
3. Otherwise, `turn.activewizard` may act during Actions at
   `/match/<match-id>/actions/<operation>`.

Every route below uses `POST`, accepts `requestid` (a canonical UUID) and
`expectedrevision` (including zero), and returns `200 {"match":{...}}` projected
for the authenticated actor. Operation-specific bodies reject unknown fields,
trailing JSON, and the obsolete `type`, `mode`, and body-level `choiceid` fields.
Draft choices continue using `/draft` with their existing request format.

| Action suffix | Fields beyond `requestid` and `expectedrevision` |
| --- | --- |
| `play-card` | `cardid`, optional `targets`, optional `support`, optional `payment`, optional `discard` |
| `attack` | `cardid`, `minions`, optional `support`, optional `payment` |
| `attempt-challenge` | `cardid`, `challengeid`, optional `minions`, optional `support`, optional `payment` |
| `activate-ability` | `cardid`, `abilityid`, optional `support`, optional `payment` |
| `research-element` | `cardid`, `element` |
| `research-spell` | `cardid`, optional `elements` |
| `heal-with-focus` | `cardid`, optional `support`, optional `payment` |
| `end-turn` | None |

| Choice suffix | Fields beyond `requestid` and `expectedrevision` |
| --- | --- |
| `maintenance` | `keep`, optional `payment` |
| `kneel` | `kneel` |
| `resolve-effect` | `effectid` |
| `skip-effect` | `effectid` |
| `undertake-hunt` | `effectid`, optional `minions`, optional `support`, optional `payment` |
| `research` | `designid` |
| `end-discard` | `cards` |
| `hand-discard` | `cardid`, `disposition` |
| `asset-loss` | `cardid` |
| `scry` | `take`, `order` |

| Reaction suffix | Fields beyond `requestid` and `expectedrevision` |
| --- | --- |
| `play-card` | `cardid`, optional `targets`, optional `support`, optional `payment` |
| `defend` | optional `minions`, optional `support` |
| `pass` | None |

For example, play a Minion at `POST /match/<match-id>/actions/play-card`:

```json
{"requestid":"9493d7de-393b-49e2-a13b-bc41dc33a3de","expectedrevision":12,"cardid":"<minion-instance>","payment":["<wis-instance>"]}
```

Pass at `POST /match/<match-id>/reactions/pass`:

```json
{"requestid":"9493d7de-393b-49e2-a13b-bc41dc33a3df","expectedrevision":13}
```

A required choice pauses resolution and Priority until it finishes. Choice routes
must match its kind. `asset-loss` serves both Asset-discard and Attack-loss prompts.
`undertake-hunt` alone accepts Full Moon's hunt commitments; `resolve-effect`
cannot undertake it. `skip-effect` takes the unprefixed effect ID. `kneel` requires
an explicit boolean, including valid `false`. `hand-discard` requires explicit
`disposition:"discard"` or `"destroy"`. Empty `keep`, `cards`, and `order` arrays
are valid when permitted by the pending choice.

Retry a lost response with the exact operation, choice ID, revision, UUID, and
input. Receipt lookup precedes current choices, Priority, and terminal-state
checks, so accepted inputs replay their recorded view even after restart or
victory. Fetch current state before composing new input. Reusing an accepted UUID
with changed input or operation returns `409 request_id_conflict`.

This is a coordinated backend/client replacement: `/commands` and the broad
`/reactions` route have been removed. Historical receipts remain untouched and
reserve their UUIDs; new routes never reinterpret them. Match snapshots,
catalogue/rules versions, response projections, and match storage are
unchanged by the route replacement. The catalogue table is installed separately
as described above. Deploy the backend and CLI together and restart wrapper sessions.

Resource cards in `payment` are spent immediately. Leftover resources remain
available for the active turn, and costs stay spent if an effect is countered.
Playing a resource as an ordinary action opens its normal response exchange.
Income creates ephemeral cards; unused resources and ephemeral cards expire at
the active turn's end, including a turn forfeited by kneeling.

All card references are instance IDs, except Research offers, which select a
design ID. A support entry is `{"cardid":"<id>","domains":["earth"]}`.
Elements need no `domains`; Focus on a Challenge or ordinary action chooses
one domain, while Focus on Attack/Defence needs none. Keen Insight chooses one
domain and Flowing Form chooses two distinct domains. Other Support Spells
need no domains and can be committed only with a Challenge.

To Attack, use `actions/attack` with Subjugate Rival and `minions` containing ready Might Minions,
and optionally Focus in `support`. Only one Attack may be initiated each turn.
Newly played Minions can defend immediately but cannot Attack or undertake
Challenges before their owner's next turn. Reactions use `/reactions/play-card` during Priority; `targets` identify the printed card targets, or
a pending Spell entry's ID for Counterspell. Scalding Mist and Freezing Fog put the removal
or return target first and the optional different recovery-delay target second.

Research uses `actions/research-element` with an `element`, or
`actions/research-spell` with zero, one, or two supporting Element-card instance
IDs in `elements`. Its persisted `research` choice supplies the design offer.
Standalone Focus healing uses `actions/heal-with-focus`; another Focus in
`support` can contribute toward its 2 Wis cost. Generic `play-card` excludes
Attack, Research, and standalone Focus healing. Purifying Steam requires one
hand-card ID in `discard` as an additional cost.

Required choices expose their ID, actor, type, source, selection bounds, and
context. `effects` choices list structured `effects`. Maintenance uses `keep` for
Assets whose upkeep will be paid, plus optional resource-card `payment`;
unselected upkeep Assets are discarded. End Phase's `end-discard` uses `cards`
to discard, including an empty list when already within hand size. The server
refills the retained hand and begins the next turn. Scrying Pool privately exposes
its inspected cards; submit `take` and every remaining inspected ID in top-first
`order`. No other gameplay interleaves with these choices, and restarts preserve
them exactly.

Other players see choice identity and timing, but no private hand or scrying
options. Public state includes resources, immunity, rumours and reserved
Investigate draws, commitments, Stack entries, and public modifiers. Owner draw
piles remain sorted by instance ID in normal views, never by draw order.

New gameplay errors include `400 invalid_choice` and `409 invalid_timing`,
`wrong_priority`, `stale_choice`, or `match_finished`; the existing malformed
request, membership, stale revision, retry conflict, and deletion errors apply.
Invalid commands leave state, revision, receipts, and randomness unchanged.

A winning match has `status:"finished"`, `outcome:"victory"`, and `winners`
containing user IDs. A zero-survivor draw has `outcome:"draw"` and `winners:[]`.
Five-step Domination ends the game after the awarding effect or submission;
last-survivor victory waits until the exchange finishes. Finished state remains
readable, and new gameplay is rejected. Two-player matches cannot use kneeling.

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

### CLI gameplay and guided sessions

The sibling `grimcli` repository provides explicit action, choice, and reaction
commands matching these routes. `challenge`, `activate`, and `react` name Challenge,
ability, and reaction-card operations; `choose-research` answers a Research offer.
Every submission requires explicit match ID, expected revision (including zero),
and a canonical request UUID. Choice commands also require `--choiceid`. Repeatable
flags compose Minions, support/domain pairs, payments, targets, discards, and order;
the CLI imports only `github.com/davelawson/grim/api`.

Its Bash session wrapper guides complete games using numbered pickers and a
review before submission, including all required choices and private scrying.
It keeps tokens, pinned catalogues, and exact pending submissions in memory;
refresh is manual. Refresh or navigation retains uncertain submissions with
the original endpoint, revision, UUID, and arguments. After success, retry only
the authoritative GET. Account changes and exit clear session data. See
[`grimcli` usage](../grimcli/README.md) for flags and interactive controls.

To run the optional CLI integration suite against production test routes:

```sh
(cd ../grimcli && go build -o /tmp/grimcli-gameplay .)
GRIMCLI_TEST_BINARY=/tmp/grimcli-gameplay go test ./internal/server -run '^TestCLICompleteGames$' -count=1 -v
```

It creates temporary databases and plays complete two-, three-, and four-player
games through the external CLI. The normal backend suite skips these tests when
`GRIMCLI_TEST_BINARY` is unset. Regenerate local Swagger with `swag init --output docs`
after changing the public projection; generated Swagger files remain ignored by Git.

### List your matches and lobbies

Authenticated `GET /matches` and `GET /lobbies` list memberships for the user
identified by the raw token in the `Authorization` header. Responses are
`{"matches":[{"id":"…","status":"setup"}]}` and
`{"lobbies":[{"id":"…","status":"open"}]}`. Each entry contains only ID and
status; use the single-item endpoints for details. Empty results return arrays
(`[]`), never null.

Lists include every nondeleted state, are ordered by ID, and are not paginated.
Lobby results reflect current memberships; leaving removes a lobby from that
user's list. Match listings read metadata without decoding saved game state.
Administrators also receive only their own memberships. Missing or invalid
tokens return 401; storage failures use the standard generic 500 error.
