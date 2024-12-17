# SHIT THAT NEEDS DOING

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
    - Game
      - Launch
  - Update Game Settings
  - Launch Game
  - Generic Lobby following login screen?

## During Game Functionality

## End of Turn Functionality
