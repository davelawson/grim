# SHIT THAT NEEDS DOING

## Authentication

- No authentication required:
  - create user
  - login
- Authentication required:
  - everything else
- Verifying bearer token
  - Some kinda controller that handles it automatically?
  - Should probably be route by route

## Permissions

- What would permissions handle?
  - lobby powers
  - terminating game
  - game participation

## Cleanup Logging

- fix calls to fmt.Println()
  - remove the {}s, those come from another language somewhere

## Database

- sqlite3
  - how do we handle migration?
    - on hold for now
  - how do we dump the database?
  - connecting MySql workbench (or some other GUI app)

## Error Handling

- All errors need to result in appropriate http status
- Ultimately, all errors should only be logged once

## Testing

- Scripts for creating test data
  - No longer in raw SQL since the addition of the hashes
  - Need to actually have scripts with curl commands to exercise the end points
- unit testing
- integration testing

## Between Game Functionality

- Endpoints
  - Create User
  - Login
  - Create Lobby
  - Get Lobbies
  - Join Lobby
  - Leave Lobby
  - Update Game Settings
  - Player Ready
  - Launch Game
  - Generic Lobby following login screen?

## During Game Functionality

## End of Turn Functionality
