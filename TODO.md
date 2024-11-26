# SHIT THAT NEEDS DOING

## Authentication

- User log in
- Some kinda token on Endpoints
- Access to particular games

## Database

- sqlite3
  - how do we handle migration?
  - how do we dump the database?
  - connecting MySql workbench (or some other GUI app)
  - a more elegant way of handling errors arising for queries gone awry

## Passing Data to the Server

- Use:
  - path args (internal requests only, such as getting a player by id)
  - request bodies (all external requests)

- Avoid:
  - query params (inconsistencies with encoding and marshalling scare me)

## Testing

- Scripts for creating test data
- unit testing
- integration testing

## Cleanup

- Get rid of all the album crap once we have basic user routes working

## Between Game Functionality

- Endpoints
  - Create Lobby
  - Join Lobby
  - Leave Lobby
  - Update Game Settings
  - Player Ready
  - Launch Game

## During Game Functionality

## End of Turn Functionality
