# DeSSE Go
A Demon's Souls server emulator implemented in Go and using SQLite for persistent
data storage.

## Acknowledgements
This is heavily based upon the work done by ymgve in the [desse][1] project. This
was the first working version of a Demon's Souls server after the official
shutdown. We're in the debt of the packet analysis and efforts from ymgve.

## Features
This server implements the core feature set for Demon's Souls:

* Login
* Character creation
* World tendency
* Blood messages
* Wandering ghosts
* Blood stains
* Summons

It should be noted that the full summon/multiplayer flow is not handled by this
server and relies on the Sony Playstation Network matchmaking system. There
is every chance they'll drop support for PS3 Demon's Souls at some point.

## Requirements
* [Go][2] 1.13+

## Installation
```bash
$ go get -u github.com/danmrichards/dessego/cmd/server/...
```

## Building From Source
Clone this repo and build the binary:

```bash
$ make build
```

## Usage
```bash
Usage of ./bin/dessego-linux-amd64:
  -seed
        Seed database tables with legacy data
```

[1]: https://github.com/ymgve/desse
[2]: https://go.dev