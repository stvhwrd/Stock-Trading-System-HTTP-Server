[![Build Status](https://travis-ci.com/stvhwrd/SENG-468-HTTP-Server.svg?token=pkeb5Zss7eZT4vjxYMRQ&branch=master)](https://travis-ci.com/stvhwrd/SENG-468-HTTP-Server)

# SENG 468 HTTP Server

The HTTP server component of the day trading system designed and built for SENG 468 at UVic

## Installation

> You must have [Go](https://golang.org/) and [Git](https://git-scm.com/) installed.

1. Clone this repository.

  * With HTTPS:

  `git clone https://github.com/stvhwrd/SENG-468-HTTP-Server.git $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

  * or with SSH:

  `git clone git@github.com:stvhwrd/SENG-468-HTTP-Server.git $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

2. Build the Go files.

  `go build $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

## Updating

If you've already installed this repository as per [installation](#installation), you may update in-place.

  `go get -u $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server`

## Usage

In its production environment, the HTTP server will communicate with a Transaction Server and a Database Server.
The ports on which to communicate **must** be specified at runtime, using the following flags:

* `-dbport`

eg. `-dbport=4564`

* `-txport`

eg. `-txport=5675`

> The HTTP server runs on port 80 by default.

The HTTP server can be started up using the `go run` command.

eg. `go run $GOPATH/src/github.com/stvhwrd/SENG-468-HTTP-Server -dbport=4564 -txport=5675`


## See Also

Components and documentation of the day trading system:

* [Project Documentation](https://github.com/stvhwrd/SENG-468-Documentation)
* [Web Client](https://github.com/dukeng/)
* [HTTP Server](https://github.com/stvhwrd/SENG-468-HTTP-Server)
* [Common Library](https://github.com/kurtd5105/SENG-468-Common-Lib)
* [Transaction Server](https://github.com/kurtd5105/SENG-468-Transaction-Server)
* [Quote Management Server](https://github.com/sterlinglaird/)
* [Database Server](https://github.com/sterlinglaird/SENG-468-Database-Server)
* [Workload Generator](https://github.com/dukeng/SENG-468-Workload-Generator)
