# SHIT THAT NEEDS DOING

## Authentication

- Make sure emails are always lower case.  Controllers will be responsible for making sure of this.
- No authentication required:
  - create user
  - login
- Authentication required:
  - everything else
- Verifying bearer token
  - Clean this up to be more easily re-used

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
  - connecting MySql workbench (or some other GUI app)
- implement some kinda caching

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
