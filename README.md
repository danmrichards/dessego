# DeSSE Go [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/danmrichards/go8080/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/danmrichards/dessego)](https://goreportcard.com/report/github.com/danmrichards/dessego)
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

## Connecting from Demon's Souls
### Native PS3
To start with you'll need some sort of DNS proxy where you can configure the following URLs to route to your dessego server:

```
ds-eu-c.scej-online.jp
ds-eu-g.scej-online.jp
c.demons-souls.com
g.demons-souls.com
cmnap.scej-online.jp
demons-souls.scej-online.jp
```

Then you can follow these steps:
1. On your PS3 Navigate to the following menu: `Settings > Network > Internet Connection Settings > Custom > Enter Manually -> Scroll to DNS Section -> Manual`
2. Set Primary DNS to the host/port of your DNS proxy

### RPCS3
First ensure you have [RPCS3][3] installed and Demon's Souls is working (figure that one out yourselves!)

1. Open RPCS3 & Create a custom configuration of your game, proceed into the Network settings and set the following options as said:
2. Set Network Status to Connected
3. Set PSN Status to RPCN
4. Set DNS to `8.8.8.8`
5. Set IP/Host Switch to `ds-eu-c.scej-online.jp=<DESSEGO_IP>&&ds-eu-g.scej-online.jp=<DESSEGO_IP>&&c.demons-souls.com=<DESSEGO_IP>&&g.demons-souls.com=<DESSEGO_IP>&&cmnap.scej-online.jp=<DESSEGO_IP>&&demons-souls.scej-online.jp=<DESSEGO_IP>` where `DESSEGO_IP` is the host/port where you're running dessego.
6. Save Configuration
7. Go to RPCS3 main menu, proceed to 'Configuration', then 'RPCN'
8. Set Host to `np.rpcs3.net`
9. Set NPID to your preferred username
10. Set Password to your preferred password for RPCN
11. Click Create Account
12. You will be asked to enter you email, to which you will receive a email with a token inside it
13. In the Token Field, enter the token you received in the email.
14. Save, and Start the game

[1]: https://github.com/ymgve/desse
[2]: https://go.dev
[3]: https://rpcs3.net
