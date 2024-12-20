# SHIT THAT NEEDS DOING

## Endpoint URL Construction

- How do we extract path params?
  - Implement GET /lobby to use a UUID path param
- How do we extract query args?
  - Implement GET on /user to use a query arg for get by email

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
    - on hold for now
  - how do we dump the database?
  - transaction management?
- implement some kinda caching

## Testing

- create a new script for populating the database
- unit testing
- integration testing

## Between Game Functionality

- Endpoints
  - Implemented
    - Login
    - User
      - Create
      - GetByEmail
    - Lobby
      - Create
        - owner_id should probably be removed in favour of the join table having a permissions field
  - Upcoming
    - Lobby
      - Join
        - lobby-player join table
          - might end up being where permissions are stored
          - store ready status here as well?
      - Delete
      - Leave
      - Get
        - My Lobbies
        - Search by lobby name
        - Search by participant
        - All lobbies

## During Game Functionality

## End of Turn Functionality
