# SHIT THAT NEEDS DOING

## Authentication

- how do we get https working?
- authentication endpoint
  - compares passHash
  - updates sessions table with new authentication token
  - returns authentication token
- all other endpoints expect an authentication bearer token
- How are we gonna handle permissions?
  - For now, don't worry about deleting users and games.  That can be done in the console directly.

## Database

- sqlite3
  - how do we handle migration?
  - how do we dump the database?
  - connecting MySql workbench (or some other GUI app)
  - a more elegant way of handling errors arising for queries gone awry

## Testing

- Scripts for creating test data
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
