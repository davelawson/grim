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

### Generate SSL key and cert

To listen to https instead of just http, we need to generate a cert and key and serve them up.
This needs to be done outside of the repo, since it obviously shouldn't be committed or shared, and also because it needs to be different for every machine that is running a server.

1. Set the GRIM_SSL environment variable to point to a folder on the server that will contain the crt and key files generated below.
  Eg: `export GRIM_SSL=$HOME/ssl`
1. Move into the GRIM_SSL directory and create the cert and key files.
    1. `cd $GRIM_SSL`
    1. `openssl genrsa -out grim.key 2048`
    1. `openssl ecparam -genkey -name secp384r1 -out grim.key`
    1. `openssl req -new -x509 -sha256 -key grim.key -out grim.crt -days 3650`
        - Just leave everything blank (default)
    1. verify that the folder contains 2 files:
        - `grim.key`
        - `grim.crt`

### Run Service

1. Run the app using go.
  Eg: `go run app.go`
1. Issue requests against the api found at `http://localhost:8080`.

### Testing

Currently, testing is extremely sparse.  There are a shell script files in the test folder that can be run, but they may require some hand holding, such as resetting the database before each run.

#### cURL

Since everything is behind a simple web RESTful web server, we can use cURL to exercise endpoints.  Many endpoints require an authentication token, which can be obtained by calling the /login endpoint.
The web server is currently using a self-signed certificate, so curl will refuse unless the `-k` argument is provided.
Here is an example of a command to query a user:
`
curl -X 'POST' \
  'https://localhost:8080/user/getbyemail' \
  -H 'accept: application/json' \
  -H 'Authorization: wymKUm3TEyDb+UBtuYS31fEP/B+fup6zK2KcQrVY3ls=' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "tim@aol.com"
}'`

### Swagger

End points are documented using a swagger interface.  This also allows manual testing of the end points.
Swaggo is the implementation we're using.  It derives meaning from comments in the code, which isn't great, but we'll stick with it a while longer.

To setup Swagger:

1. Install swagger: `go install github.com/swaggo/swag/cmd/swag@latest`
1. Make sure that `$(go env GOPATH)/bin` is on your path
1. Generate the swagger documentation: `swag init .`
1. Run the application (detailed in earlier section)
1. Open the web page: `http://localhost:8080/swagger/index.html`
