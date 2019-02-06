[![Build Status](https://travis-ci.com/stvhwrd/SENG-468-HTTP-Server.svg?token=pkeb5Zss7eZT4vjxYMRQ&branch=master)](https://travis-ci.com/stvhwrd/SENG-468-HTTP-Server)

# SENG 468 HTTP Server

The HTTP server component of the day trading system designed and built for SENG 468 at UVic

## Installation

> You must have [Go](https://golang.org/) and [Git](https://git-scm.com/) installed.

1. Clone this repository.

-   With HTTPS:

`git clone https://github.com/stvhwrd/SENG-468-HTTP-Server.git $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

-   or with SSH:

`git clone git@github.com:stvhwrd/SENG-468-HTTP-Server.git $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

2. Install the dependencies.

`go get -v -u github.com/joho/godotenv`

3. Build the Go files.

`go build $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

## Configuration

In its production environment, the HTTP server will communicate with a Transaction Server, a User Account Database Server, and a Logging Server.

The ports on which to communicate **must** be specified at runtime, using a `.env` file.

The file should be located in the same directory as `main.go` (as it is by default in this repository).

## Usage

The HTTP server can be started up using the `go run` command after building.

eg. `cd $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server && go build && go run .`

## Updating

If you've already installed this repository as per [installation](#installation), you may update in-place.

`go get -u $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

## See Also

Components and documentation of the day trading system:

- [Common Library](https://github.com/kurtd5105/SENG-468-Common-Lib)
- [Database Server](https://github.com/sterlinglaird/SENG-468-Database-Server)
- [HTTP Server](https://github.com/stvhwrd/SENG-468-HTTP-Server)
- [Logging Server](https://github.com/dukeng/SENG-468-Logging-Server)
- [Project Documentation](https://github.com/stvhwrd/SENG-468-Documentation)
- [Transaction Server](https://github.com/kurtd5105/SENG-468-Transaction-Server)
- [Web Client](https://github.com/dukeng/SENG-468-Web-Client)
- [Workload Generator](https://github.com/dukeng/SENG-468-Workload-Generator)
