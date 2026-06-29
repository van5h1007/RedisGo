# GoRedis

A small Redis server clone written in Go.

It implements a subset of the Redis wire protocol and supports a handful of
commands (`PING`, `SET`, `GET`, `DEL`, `EXISTS`, `HELLO`, `COMMAND`, `CLIENT`)
backed by an in-memory key-value store. It's compatible with `redis-cli` and
the official `go-redis` client.


## Architecture

```
TCP client
    |
    v
acceptLoop()  -- accepts connections, spawns one goroutine per peer
    |
    v
Peer.readLoop() -- parses RESP values off the wire, turns them into
    |               typed Commands, sends them as Messages
    v
msgCh (channel) ---------------------+
addPeerCh / delPeerCh (channels) ----+--> loop() -- the ONLY goroutine
                                      |     that touches shared state
                                      v
                              handleMessage()
                                      |
                                      v
                              KV store (map + RWMutex)
```

Only one goroutine (`loop()`) ever reads or writes the `peers` map and the
`KV` store. Every connection's goroutine talks to it exclusively through
channels — this is why there's no manual locking needed around `peers`, and
why a `select` (rather than a simple channel range) is required once more
than one channel needs watching at once.

## Project layout

| File | Responsibility |
|---|---|
| `main.go` | Parses flags, boots the `Server` |
| `server.go` | `Server`/`Config` structs, `acceptLoop`, `loop` (the event loop), `handleMessage` |
| `peer.go` | `Peer` type wrapping a connection |
| `message.go` | `Message` — couples a parsed `Command` with the `Peer` it came from |
| `command.go` | `Command` types (`SetCommand`, `GetCommand`, etc.) and RESP-array-to-command parsing |
| `kv.go` | `KV` — the in-memory store (`map[string]string` + `RWMutex`) |
| `resp.go` | RESP **reader**: parses bulk strings and arrays off a `bufio.Reader` |
| `resp_writer.go` | RESP **writer**: simple strings, errors, bulk strings, integers, RESP3 maps |
| `server_test.go` | Integration test using the real `go-redis` client |
| `peer_test.go` | Raw-socket smoke test for the `HELLO` handshake |

## Supported commands

| Command | Behavior |
|---|---|
| `PING` | Replies `PONG` |
| `SET key val` | Stores a value |
| `GET key` | Returns the value, or a null bulk string if missing |
| `DEL key` | Deletes a key, returns `1`/`0` for found/not-found |
| `EXISTS key` | Returns `1`/`0` |
| `HELLO` | Returns a minimal RESP3 map (handshake used by real Redis clients) |
| `COMMAND` | Returns an empty array (satisfies clients that probe this on connect) |
| `CLIENT ...` | Acknowledges with `OK` (subcommands like `SETNAME` aren't actually implemented) |

## Running it

```bash
go run . -listenAddr :5001
```

`-listenAddr` defaults to `:5001` if omitted.

Talk to it with `redis-cli`:

```bash
redis-cli -p 5001
> SET foo bar
OK
> GET foo
"bar"
> DEL foo
(integer) 1
```

Or raw RESP over `nc`:

```bash
printf '*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n' | nc localhost 5001
```

## Running the tests

```bash
go test -v ./...
```

Note: `peer_test.go`'s `TestHelloCommand` expects a server already running on
`:5001` — start one in a separate terminal first. `server_test.go`'s
`TestOfficialRedisClient` starts its own server, so it doesn't need this.


## Possible next steps

- `EXPIRE`/`TTL` support (requires a background goroutine sweeping expired
  keys)
- `INCR`/`DECR`
- Multi-key `DEL`
- Pub/sub, using the now-tracked `peers` map to broadcast messages
- Persistence to disk (e.g. periodic snapshot via `encoding/gob`)
- A `Dockerfile` (multi-stage build) and `docker-compose.yml` for running the
  server and a test client as separate containers
