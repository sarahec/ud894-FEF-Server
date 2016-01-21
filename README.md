# Back-End Server for Front-End Frameworks course

This server provides support for the Front-End Frameworks coding labs. It both
serves the app files and provides a REST interface for reading and storing data.

## Building and running

1. [Download](https://golang.org/dl/) and install the Go programming language
2. Download or clone this project
3. Compile the server: `go build -o main`
4. Run the `main` program to start the server . Use the `--www=` flag to point
to your front-end code (e.g. on Mac OS X and Linux
`./main --www=../FEF-UdaciMeals-Backbone`)
5. Use the `--logrest` flag to see all of the incoming and outgoing traffic
from the server.

## Server details

All of the data is stored in JSON format in the `_data` directory. `menu.json`
is the storage file.

The server implements a REST API at `/api/items`:
* `GET /api/items` (no trailing slash) returns a JSON array of menu items
* `PUT /api/items` is disallowed (you cannot put the whole array at once)
* `GET /api/items/[id]` (e.g. `GET /api/items/strawberry-shortcake`) gets the
menu item with the specified ID and returns it as JSON
* `PUT /api/items/[:id]` takes a menu item (JSON format, in the body) and
updates the existing item if the ID exists or appends a new one if the ID
doesn't exist yet

The server also serves web files, mapping the directory specified by the `--www`
  flag to `/` (e.g. `./main --www=../web-files` serves `../web-files/index.html`
  as `/index.html`)
