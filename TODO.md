# SHIT THAT NEEDS DOING

## Ongoing

- Standardized Logging

## CURRENT

- Refactor user package to match with our standards

## Updating lists vs Add and remove

When we have lists, should we use a simple update on the base item to manage them, or should we include add/remove endpoints?

## Robust Endpoint Logging

- Log requester identity, URL, method, and diagnostic request details while
  excluding credentials and authentication tokens.

### UUIDs in Endpoints

| Use Case | Operation | Endpoint |
| --- | --- | --- |
| Lookup User by UUID | GET | /user/12341234-1234-12341234-1234-12341234 |
| Update Lobby | PUT | /lobby/3455342-2354-35245234-3245-23455234 |
| Delete a Match | DELETE | /match/23452345-2345-23452345-2345-23452345 |
| Create a new User | POST | /user |

- Request that we have UUIDs for each model
- Doesn't graceful accommodate querying by any other fields.

| User Case | Operation | Endpoint |
| --- | --- | --- |
| Get User by Email | POST | /user/getByEmail |

### UUID-free Endpoints

| Use Case | Operation | Endpoint |
| --- | --- | --- |
| Lookup User by UUID | POST | /user/getbyid |
| Update Lobby | PUT | /lobby|
| Delete a Match | DELETE | /match |
| Create a new User | POST | /user |

- Having to use POST when performing a lookup sucks.

## Permissions

- What would permissions handle?
  - lobby powers
  - terminating game
  - game participation

## Database

- sqlite3
  - how do we handle migration?
    - introduce ordered incremental migrations and saved-document schema
      upgrades as part of the backend expansion in [WORK.md](WORK.md)
  - how do we dump the database?
  - transaction management?
- implement some kinda caching

## Testing

- create a new script for populating the database
- unit testing
- integration testing

## Between Game Functionality

Implement the confirmed [backend expansion](ARCHITECTURE.md#planned-match-backend)
in its documented sequence. Lobby readiness/departure/ownership transfer,
saved match setup, atomic one-time launch, and administrative termination are
implemented, including UpdateLobby authorization. Incremental database and
saved-document upgrades remain deferred; the current creation script sets up
the complete lifecycle schema directly.

- Endpoints
  - Implemented
    - Login
    - User
      - Create
      - GetByEmail
    - Lobby
      - Create
        - owner_id should probably be removed in favour of the join table having a permissions field
      - Get (members only)
      - Update and ownership transfer (owner only)
      - Add/remove members and voluntary departure
      - Delete an open lobby
      - Set own readiness
      - Atomic one-time launch
    - Match
      - Participant lookup
      - Administrative termination
  - Upcoming
    - Lobby
      - Self-service join policy (currently owner-managed)
      - Get
        - My Lobbies
        - Search by lobby name
        - Search by participant
        - All lobbies

## During Game Functionality

## End of Turn Functionality

## Game Design

- Complete the small first-playable catalogue: working victory challenges
  and ranks beyond the four drafting placeholders, Spell and Artifact creation
  pools, and remaining Research rules. These are prerequisites for the
  complete-game backend milestone defined in [ARCHITECTURE.md](ARCHITECTURE.md).
- Implement the agreed [drafting rules](GAME.md#drafting), including moving
  random seating from launch to draft completion.
- Playtest the agreed drafting, card anatomy, turn phases, response exchanges, and
  kneeling flow in [GAME.md](GAME.md).
- Define immunity duration precisely under continuous turns and changing
  neighbours.
- Decide whether separate defensive preparation exists alongside standing
  Might and response commitments.
- Resume the [parked Spell design discussion](docs/spell-design-notes.md) when
  choosing the next design slice.
