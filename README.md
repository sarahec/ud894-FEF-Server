# Back-End Server for Udacity's Front-End Frameworks course

This server provides support for the Front-End Frameworks coding labs. It both
serves the app files and provides a REST interface for reading and storing data.

## Running the server

The [binaries](./binaries) directory contains builds for Windows, Mac OS X
("darwin"), and Linux. Run the appropriate `server` program to start the server,
supplying the `--www=` flag to point to your front-end code (e.g. on Mac OS X use
`./server_darwin_amd64 --www=../../FEF-UdaciMeals-Backbone`)

Use the `--log` flag to see all of the incoming and outgoing traffic
from the server.

## Server details

The server serves web files and provides a REST API. The `--www`flag
specifies the directory to be served as `/` (e.g. `./main --www=../web-files`
serves `../web-files/index.html` as `/index.html`)

The server implements a REST API at `/api/items`:
* `GET /api/items` (no trailing slash) returns an object with a JSON array of menu items
* `PUT /api/items` is disallowed (you cannot put the whole array at once)
* `GET /api/items/[id]` (e.g. `GET /api/items/strawberry-shortcake`) gets the
menu item with the specified ID and returns it as JSON
* `PUT /api/items/[:id]` takes a menu item (JSON format, in the body) and
updates the existing item if the ID exists or appends a new one if the ID
doesn't exist yet

All of the data is stored in JSON format in the `_data` directory. `menu.json`
is the storage file. (The file at [server/assets/menu.json](./server/assets/menu.json)
provides the starting values and is compiled into the server binary.)

## Building and running

If you want to modify the code or build for another platform:

1. [Download](https://golang.org/dl/) and install the Go programming language
2. Download or clone this project
3. Change into the `server` directory
3. Compile the server: `go build`
4. Run the `server` program to start the server.

You can also use the [compile.sh](./compile.sh) script to cross-compile for
Windows, Mac OS X, and Linux.

If you want to change the initial data file, install the `rice` tool from
[https://github.com/GeertJohan/go.rice](https://github.com/GeertJohan/go.rice),
edit [server/assets/menu.json](./server/assets/menu.json), and follow the
directions for using `embed-go` before re-building the binary.
