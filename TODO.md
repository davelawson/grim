# SHIT THAT NEEDS DOING

## Authentication

- User log in
- Some kinda token on Endpoints
- Access to particular games

## Database

  Try and do it with both SQLite and MySql, and see which one does a better job.
  Migration is a must!
  Dumping is a must!

- SQLite
    Is there a flyway equivalent?
    Is there a go driver for it?
- MySql
--- Creating a Game Lobby ---
- Endpoints
      Create Lobby
      Join Lobby
      Leave Lobby
      Update Game Settings
      Player Ready
      Launch Game
--- Ending a Game Turn ---
- Endpoint to a Player's Turn
- Processing the turn end
- Sending new game state to players
