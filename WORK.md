# Current Work Slice

This file holds the design for one active slice of work. Before implementation
begins, fill in the goal and design below. Keep unresolved choices here while
the slice is in progress. `TODO.md` remains the backlog.

As each part is implemented and verified, move its lasting design details to
[ARCHITECTURE.md](ARCHITECTURE.md) for application and technical decisions or
[GAME.md](GAME.md) for game rules. Remove those completed details from this file
so it continues to describe only work in progress. When the slice is complete,
reset this file for the next slice.

For a design-only slice, move its accepted provisional rules to the appropriate
permanent document when the design is complete, even if implementation follows
in a later slice. Label those rules as provisional until they are tested.

## Active slice

No drafting design work remains. The user confirmed the complete ruleset.
Its accepted provisional rules, starting offers, and Subjugate Rival are
recorded in [GAME.md](GAME.md#drafting); terminology is in [CONTEXT.md](CONTEXT.md).
The planned change from launch-time seating to seating after drafting is
recorded in [ARCHITECTURE.md](ARCHITECTURE.md#lobby-and-match-lifecycle).
Drafting implementation and the other four victory paths remain backlog work
in [TODO.md](TODO.md). Backend commands and implementation were outside this
design-only slice.

### Completed lifecycle slice

Lobby-to-match creation, participant lookup, and administrator termination are
implemented and verified. Their lasting contract is recorded in
[ARCHITECTURE.md](ARCHITECTURE.md#implemented-lobby-and-match-lifecycle) and usage
in [README.md](README.md#launch-and-end-a-match). Drafting implementation,
gameplay, player votes, and migrations remain deferred.

The pending registration design is retained below; it is not part of the
completed lifecycle slice.

## Previous slice: user creation (retained)

The prior registration design below remains pending and is retained rather
than folded into the lifecycle slice.

User creation: registration design recorded; broader implementation remains
pending. The separately approved lowercase user-JSON subset is implemented
and verified; its lasting behavior is recorded in [ARCHITECTURE.md](ARCHITECTURE.md).
The previously confirmed backend architecture remains recorded in
[ARCHITECTURE.md](ARCHITECTURE.md#planned-match-backend).

## Goal

Create an agreed implementation plan for user creation, grounded in the
existing account-creation and authentication code. Database migrations are
deferred to a later slice.

## Design

Account creation is self-registration only. Administrator provisioning is
outside this slice. The existing unauthenticated `POST /user` provides the
starting point.

Registration requires an email address, display name, and password. Email is
the unique login identity; display names need not be unique. A generated user
ID distinguishes accounts independently of their display names.

This slice targets local/private testing with known players. The user accepted
that initial audience; registration open to the public is outside the current
target.

New accounts are immediately usable; email verification is deferred.

Emails are trimmed of surrounding whitespace, lowercased, and checked for valid
address syntax. Dots and `+suffixes` are preserved. Login uses the same email
normalization. Display names are trimmed, contain 1–80 Unicode characters,
allow internal spaces, and reject control characters.

Passwords have no length or quality validation in this slice. Passwords are
preserved exactly. The password field must be present and contain a string;
an empty string is accepted, while missing or null is rejected. The earlier
length-policy and common-password-blocklist recommendation was not accepted.

Successful registration returns `201 Created` with the new public user ID,
normalized email, and display name. Registration does not issue a token;
the user subsequently calls `/login`.

Invalid requests return 400, duplicate email returns 409, and unexpected
failures return 500. Errors are consistent JSON without internal database
details. A duplicate attempt does not change the existing account.

Registration JSON uses lowercase keys. The request fields are `email`, `name`,
and `password`; request keys must use exactly those lowercase spellings.
Capitalized or mixed-case keys are rejected with 400. Success returns
`{"user":{"id":"…","email":"…","name":"…"}}`. Errors use
`{"error":{"code":"…","message":"…"}}`, with an optional `fields` map
inside `error` for field-validation failures. Use `invalid_request`,
`email_already_registered`, and `internal_error` as stable error codes.
Use registration-specific response types so other endpoints retain their
current JSON response shapes.

Accept exactly one JSON object with all three required string fields. Reject
unknown fields, repeated fields (including differently cased versions of the
same field), nulls, incorrect types, and additional JSON values with 400.
Whitespace around the JSON object is allowed. Empty passwords remain valid.

Retain the existing scrypt password-storage format and compatible login
verification in this slice. A hashing-format upgrade is deferred. Excluding
credentials, hashes, and tokens from logs is part of this slice.

This repository contains the API and Swagger documentation; clients are
separate applications in the established architecture.

Inspection found missing input validation, generic server errors for duplicate
emails, and plaintext password logging in account creation and login. These
are implementation gaps to address in the plan. The existing architecture
already requires excluding credentials and tokens from logs; registration,
login, password-hash, and token logging must be brought into line with it.

## Open decisions

None in the registration contract. The complete plan awaits the user's
confirmation of shared understanding before implementation.

Read-only aggregate inspection of the configured database found no emails with
ASCII uppercase or surrounding ASCII whitespace, and no collision groups under
SQLite `lower(trim(email))`. This is not a comprehensive Unicode audit.
No automatic rewriting of existing accounts is planned. Existing canonical
accounts must continue to log in; unusual noncanonical legacy records would
require separate diagnosis rather than a schema migration in this slice.

The existing backend architecture and game-content backlog remain in effect.

## Implementation plan

1. Build on the implemented lowercase user models and registration casing
   check; define registration-specific response models and error codes. Complete
   request parsing locally to this endpoint, tracking field presence and duplicate
   keys; do not change Gin's decoding policy globally. Enforce one JSON object,
   exact lowercase allowed fields and string types. Keep
   password presence separate from length so an empty password is accepted.
2. Put email normalization and account-field validation behind the controller
   so service callers cannot bypass them. Trim and lowercase emails without
   provider-specific rewriting. Validate a plain email address without a
   display-name wrapper; do not perform DNS or mailbox checks. Trim names,
   count Unicode code points for the 1–80 limit, and reject control characters.
   Leave passwords unchanged. Reuse email normalization in login.
3. Extend registration's service/facade/repository return values to supply the
   created public user, including the generated UUID. Use the existing users
   table and transaction convention; retain scrypt parameters and email salt.
   Let the database email-uniqueness constraint arbitrate competing inserts,
   classify that specific constraint failure as duplicate email, and leave
   other database failures unexpected. Return success only after commit.
4. Map registration outcomes to the agreed JSON bodies and 201/400/409/500
   statuses. Return immediately after writing an error. Keep hash/token fields
   out of public models and prohibit account overwrite on duplicate requests.
5. Align the login email path and remove secret-bearing logging from user,
   auth, hashing, and token helpers. Preserve current token format and response
   conventions. Propagate token-persistence failures so login never returns
   success with a token that was not saved. Keep broader session changes out.
6. Add focused Go tests with `httptest` and isolated SQLite fixtures. Exercise
   parsing, normalization, validation, persistence outcomes, and registration
   followed by login. Create test credentials directly instead of relying on
   the undocumented shared binary fixture hash. Use controlled dependency
   failures for commit and token-write cases where appropriate.
7. Update registration Swagger annotations and README examples, regenerate
   local Swagger files through the existing Swaggo workflow, then run
   `go test ./...` and `go vet ./...`. Resolve relevant failures before declaring
   the implementation complete. Generated Swagger files are currently ignored.

## Acceptance scenarios

- A valid unauthenticated request returns 201 and one public user with a UUID,
  normalized email, and trimmed name; one matching row is committed.
- Users with the same display name and different emails both register.
- Email case and surrounding-whitespace variants identify the same account.
  Re-registering that identity returns 409 without changing its ID, name,
  password hash, or token. Competing registrations never create two accounts
  for that identity; the SQLite uniqueness-conflict path maps to 409.
- Missing or null required fields, wrong types, non-object bodies, unknown or
  duplicate fields, malformed JSON, and trailing JSON return 400 with no user
  created. Capitalized or mixed-case request keys return 400; only lowercase
  `email`, `name`, and `password` are accepted.
- Invalid email syntax, empty trimmed names, names longer than 80 code points,
  and names containing control characters return 400. Unicode names and
  names with internal spaces are accepted within the limit.
- An explicitly empty password registers and authenticates. Short, long,
  Unicode, and whitespace-containing passwords are preserved without length,
  quality, composition, or blocklist checks. Missing and null passwords fail.
- Registration returns no login token. Subsequent login accepts the normalized
  email and exact password; an incorrect password fails using current login
  response conventions. Existing normally formatted accounts remain compatible
  with the retained hash format.
- Database insert/commit failures produce registration 500 without internal
  details or a falsely successful response. Token-write failure produces login
  failure rather than an unusable success token.
- Responses and application logs do not expose passwords, hashes, or tokens,
  except for the existing intentional token returned by successful login.
  Registration errors never include raw SQL or submitted password values.

## Deferred work

Database migrations, public registration, email verification/delivery,
administrator provisioning, password length/quality policy, password-hash
format upgrades, password reset/change, session redesign, and client UI are
outside this slice. No schema changes or existing-account rewrites are planned.

## Progress

Existing user creation and authentication have been inspected. Self-registration
only, email identity, shared display names, required fields, private testing,
immediate activation, email/display-name rules, password-field presence with
empty-string acceptance, separate login, response statuses, and retained scrypt
format are confirmed. Password length and quality checks are excluded at the
user's request. Lowercase registration JSON and strict parsing are confirmed.
Strictly lowercase registration request keys are confirmed. The User glossary
entry and complete implementation plan are recorded. The user approved a
narrower implementation plan: lowercase tags across user endpoint models and
rejection of nonlowercase registration field aliases, preserving existing
status codes, lookup decoding, and other API formats. That subset is implemented;
focused HTTP/SQLite tests pass via `go test ./...`. Local Swagger was regenerated
and its user field names verified; formatting checks passed. The broader
registration plan has not been implemented or confirmed as a whole.
