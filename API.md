# Demon's Souls Server API
Demon's Souls brought From Software's unique asynchronous multiplayer to the
gaming world for the first time in 2009. This emulated server replaces the
core part of the infrastructure that the game interacts with.

## Topology
There are 3 elements in the Demon's Souls multiplayer infrastructure:

1. Game client
2. Game server
3. Sony Playstation Network

The game server handles all state management and interaction to power features
such as:

* World Tendency
* Messaging
* Replays
* Bloodstains
* Wandering ghosts
* Summoning

The latter of these, summoning, is where Sony Playstation Network comes into
play. The game client can request or offer assistance to other players, via the
summon sign mechanic. The game server stores the state of the summons, but it
is the responsibility of Playstation Network to actually perform the match
making for valid summon signs. This interaction comes from the game client and
not the game server.

The game server itself comprises 2 components; technically _at least_ 2. The
first component is the "bootstrap" and _n_ instances of regional game servers.

## Bootstrap Server
The bootstrap server is responsible for configuring a game client that has
just booted up. It provides configuration information for the client to operate
its multiplayer functions and also advertises URLs for regional game servers.

### Transport
The server uses a standard HTTP/1.1 protocol.

### Endpoints

#### Root
**Path:** `/`

**Description:** Returns a base64 encoded XML object containing game client configuration.

**Request content-type:** N/A

**Request body:** N/A

**Response content-type:** `text/plain`

**Example response:**
```
PHNzPjA8L3NzPgo8bGFuZzE+PC9sYW5nMT4KPGxhbmcyPjwvbGFuZzI+CjxsYW5nMz48L2xhbmczPgo8bGFuZzQ+PC9sYW5nND4KPGxhbmc1PjwvbGFuZzU+CjxsYW5nNj48L2xhbmc2Pgo8bGFuZzc+PC9sYW5nNz4KPGxhbmc4PjwvbGFuZzg+CjxsYW5nMTE+PC9sYW5nMTE+CjxsYW5nMTI+PC9sYW5nMTI+CjxnYW1ldXJsMT5odHRwOi8vMTI3LjAuMC4xOjE4NjY2L2NnaS1iaW4vPC9nYW1ldXJsMT4KPGdhbWV1cmwyPmh0dHA6Ly8xMjcuMC4wLjE6MTg2NjcvY2dpLWJpbi88L2dhbWV1cmwyPgo8Z2FtZXVybDM+aHR0cDovLzEyNy4wLjAuMToxODY2OC9jZ2ktYmluLzwvZ2FtZXVybDM+CjxnYW1ldXJsND5odHRwOi8vMTI3LjAuMC4xOjE4NjY4L2NnaS1iaW4vPC9nYW1ldXJsND4KPGdhbWV1cmw1Pmh0dHA6Ly8xMjcuMC4wLjE6MTg2NjcvY2dpLWJpbi88L2dhbWV1cmw1Pgo8Z2FtZXVybDY+aHR0cDovLzEyNy4wLjAuMToxODY2Ny9jZ2ktYmluLzwvZ2FtZXVybDY+CjxnYW1ldXJsNz5odHRwOi8vMTI3LjAuMC4xOjE4NjY3L2NnaS1iaW4vPC9nYW1ldXJsNz4KPGdhbWV1cmw4Pmh0dHA6Ly8xMjcuMC4wLjE6MTg2NjcvY2dpLWJpbi88L2dhbWV1cmw4Pgo8Z2FtZXVybDExPmh0dHA6Ly8xMjcuMC4wLjE6MTg2NjgvY2dpLWJpbi88L2dhbWV1cmwxMT4KPGdhbWV1cmwxMj5odHRwOi8vMTI3LjAuMC4xOjE4NjY4L2NnaS1iaW4vPC9nYW1ldXJsMTI+Cjxicm93c2VydXJsMT48L2Jyb3dzZXJ1cmwxPgo8YnJvd3NlcnVybDI+PC9icm93c2VydXJsMj4KPGJyb3dzZXJ1cmwzPjwvYnJvd3NlcnVybDM+CjxpbnRlcnZhbDE+MTIwPC9pbnRlcnZhbDE+CjxpbnRlcnZhbDI+MTIwPC9pbnRlcnZhbDI+CjxpbnRlcnZhbDM+MTIwPC9pbnRlcnZhbDM+CjxpbnRlcnZhbDQ+MTIwPC9pbnRlcnZhbDQ+CjxpbnRlcnZhbDU+MTIwPC9pbnRlcnZhbDU+CjxpbnRlcnZhbDY+MTIwPC9pbnRlcnZhbDY+CjxpbnRlcnZhbDc+MTIwPC9pbnRlcnZhbDc+CjxpbnRlcnZhbDg+MTIwPC9pbnRlcnZhbDg+CjxpbnRlcnZhbDExPjEyMDwvaW50ZXJ2YWwxMT4KPGludGVydmFsMTI+MTIwPC9pbnRlcnZhbDEyPgo8Z2V0V2FuZGVyaW5nR2hvc3RJbnRlcnZhbD4yMDwvZ2V0V2FuZGVyaW5nR2hvc3RJbnRlcnZhbD4KPHNldFdhbmRlcmluZ0dob3N0SW50ZXJ2YWw+MjA8L3NldFdhbmRlcmluZ0dob3N0SW50ZXJ2YWw+CjxnZXRCbG9vZE1lc3NhZ2VOdW0+ODA8L2dldEJsb29kTWVzc2FnZU51bT4KPGdldFJlcGxheUxpc3ROdW0+ODA8L2dldFJlcGxheUxpc3ROdW0+CjxlbmFibGVXYW5kZXJpbmdHaG9zdD4xPC9lbmFibGVXYW5kZXJpbmdHaG9zdD4=
```

**Response fields:**

| Field | Description |
| --- | ----------- |
| `ss` | Unknown. Observed to be `0` |
| `lang1` | Unknown |
| `lang2` | Unknown |
| `lang3` | Unknown |
| `lang4` | Unknown |
| `lang5` | Unknown |
| `lang6` | Unknown |
| `lang7` | Unknown |
| `lang8` | Unknown |
| `lang11` | Unknown |
| `lang12` | Unknown |
| `gameurlN` | Repeated field, where `N` is a zero based index. URL to an instance of the game server. |
| `browserurl1` | Unknown |
| `browserurl2` | Unknown |
| `browserurl3` | Unknown |
| `interval1` | Unknown, observed to be `120` |
| `interval2` | Unknown, observed to be `120` |
| `interval3` | Unknown, observed to be `120` |
| `interval4` | Unknown, observed to be `120` |
| `interval5` | Unknown, observed to be `120` |
| `interval6` | Unknown, observed to be `120` |
| `interval7` | Unknown, observed to be `120` |
| `interval8` | Unknown, observed to be `120` |
| `interval11` | Unknown, observed to be `120` |
| `interval12` | Unknown, observed to be `120` |
| `getWanderingGhostInterval` | Interval, in seconds, at which to get wandering ghost data, observed at `20` |
| `setWanderingGhostInterval` | Interval, in seconds, at which to set wandering ghost data, observed at `20` |
| `getBloodMessageNum` | Number of blood messages to retrieve, observed at `80` |
| `getReplayListNum` | Number of replays to retrieve, observed at `80` |
| `enableWanderingGhost` | Pseudo-boolean (1 or 0) indicating if wandering ghosts are enabled |