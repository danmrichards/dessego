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

**Request fields:** N/A

**Response content-type:** `text/plain`

**Response fields:**

| Field | Description |
| :--- | :--- |
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

## Game Server
The Demon's Souls game server is responsible for handling interactions with the
game client to power it's asynchronous multiplayer functionality. Persistent
storage is used to store data for player characters, replays and messages left
in the game world.

### Transport
The server uses a standard HTTP/1.1 protocol.

### Request
Demon's Souls makes API requests using standard HTTP verbs (e.g. POST) but with
encrypted contents. The body of an API request is an AES encrypted representation
of an HTTP form POST. Requests are encrypted using the AES key `11111111222222223333333344444444`

As an example, first imagine a POST body like so:

```
ver=100&characterID=foobar&index=1
```

The game client would then need to AES encrypt this using the known key. This
encryption method works using blocks (16 bytes in size), as a result the raw
value may need to be padded to the next whole block (i.e. the length must be a
multiple of the block size). As a result, when decrypting the POST body into set
of HTTP url params, you may see a trailing ampersand.

### Response
Responses sent back from the server to the game client are not encrypted. Instead,
they are base64 encoded byte sequences. It should also be noted that the byte
sequence is expected to be followed by a new line character, otherwise the
game client cannot parse it.

The format of the byte sequence in a server response is as follows:

| 0            | 1 .. 4      | 5 .. n |
| :----------- | :---------- | :----- |
| Command flag | Data length | Data   |

The command flag is a single byte that indicates to the game client the type
of response being returned.

### Endpoints

All response fields here are contained in the `Data` section of the response
byte sequence. Byte ordering is shown from `0` but in reality starts from `5`
in the real response.

#### Login
**Path:** `/cgi-bin/login.spd`

**Description:** Logs a new client into the server

**Request content-type:** N/A

**Request fields:** N/A

**Response content-type:** `text/plain`

**Response fields:**

| Field | Description |
| :--- | :--- |
| Status | Status of the server. Possible values:<br>0x00 - present EULA, create account<br>0x01 - present MOTD, can be multiple<br>0x02 - "Your account is currently suspended"<br>0x03 - "Your account has been banned."<br>0x05 - undergoing maintenance<br>0x06 - online service has been terminated<br>0x07 - network play cannot be used with this version |
| Data | Dependent on the status. In the message-of-the-data case, contains the encoded message string |

#### Time
**Path:** `/cgi-bin/getTimeMessage.spd`

**Description:** Get a time message from the server

**Request content-type:** N/A

**Request fields:** N/A

**Response content-type:** `text/plain`

**Response fields:**

Unknown

#### Initialise character
**Path:** `/cgi-bin/initializeCharacter.spd`

**Description:** Initialises a new character

**Request content-type:** N/A

**Request fields:**

| Field | Description |
| :--- | :--- |
| `characterID` | ID of the new character |
| `index` | Character index, allows for multiple characters on the same client |
| `ver` | Game client version |

**Response content-type:** `text/plain`

**Response fields:**

| Field | Description |
| :--- | :--- |
| Character ID | ID of the created character |
| Terminator | Termination byte (0x00) |