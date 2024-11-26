# Grimoire

Powerful wizards compete with one another to be the first to Ascend in this turn-based game.

## Components

### API

The Go application is a simple RESTful api the can receive requests from any sort of front end such as a Discord Bot, a Web Page, or a CLI.

### DB

State is stored in a simple sqlite3 database.

## Running

### Setup Database

1. Make sure that you have sqlite3 installed.
1. Set the GRIM_DB environment variable to point to where you want the sqlite3 db stored.
  Eg: `export GRIM_DB=$HOME/sqlite/test.db`
1. Create the database by executing the create-database.sql file.
  Eg: `sqlite3 $GRIM_DB < ./sql/create-database.sql`
1. Populate the database with some test data.
  `sqlite3 $GRIM_DB < ./sql/create-test-data.sql`

### Run Service

1. Run the app using go.
  Eg: `go run app.go`
1. Issue requests against the api found at `http://localhost:8080`.

### Testing

Currently, testing is extremely sparse.  There are a shell script files in the test folder that can be run, but they may require some hand holding, such as resetting the database before each run.

### Swagger

End points are documented using a swagger interface.  This also allows manual testing of the end points.
Swaggo is the implementation we're using.  It derives meaning from comments in the code, which isn't great, but we'll stick with it a while longer.

To setup Swagger:

1. Install swagger: `go install github.com/swaggo/swag/cmd/swag@latest`
1. Make sure that `$(go env GOPATH)/bin` is on your path
1. Generate the swagger documentation: `swag init .`
1. Run the application (detailed in earlier section)
1. Open the web page: `http://localhost:8080/swagger/index.html`
