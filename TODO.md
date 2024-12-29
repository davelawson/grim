# SHIT THAT NEEDS DOING

## CURRENT

### Strictly segregate concerns between Lobby Controller and Lobby Service

### Move to using path params for Lobby controller stuffs

## Cleaning up Lobby-Related Endpoints

### Controller responsibility

Controllers should only be responsible for the http-related concerns of the web server.  More precisely:

- Extract path args
- Receive and convert Request bodies
- Map underlying errors to Http codes
- Generate http response
- Invoke the responsible service

### Service responsibility

- Transaction management
  - We need to achieve this at the top level of the service, so we don't have multiple transactions when services invoke one another
  - All the rest of the business logic

## DB SHIT

- Atomicity of Operation
  - We need a way to have transactions that span an entire controller operation, rather than a simple repo query.

## Updating lists vs Add and remove

When we have lists, should we use a simple update on the base item to manage them, or should we include add/remove endpoints?

- lobby users

## Robust Endpoint Logging

- We need a generic way of logging requests, including token, request user, URL and method, and request body.

## OOP Stuffs

- Try making a package that more closely acts as an 'object'.  Will we have far too much exposed?  Big messy namespace?  We'll see...

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
