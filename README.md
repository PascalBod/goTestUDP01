# UDP client / server in Go

The `client.go`, `gateway.go`and `server.go` programs perform following actions:

* `client.go` (the *client*) sends a 5-byte packet to `gateway.go` (the *gateway*) every 5 seconds
* last byte of the client packet contains a counter which starts from 0 and is incremented after each transmission
* the gateway forwards every received packet to `server.go` (the *server*)
* the server adds two bytes to every received packet, and echoes it back to the gateway
* the gateway forwards every packet received from the server to the client
* the client displays every packet received from the gateway

Following ports are used:

```
+--------+                    +---------+                    +--------+
|        |------------------->|         |------------------->|        |
| client |20001          20000| gateway |ephemeral      30000| server |
|        |<-------------------|         |<-------------------|        |
+--------+                    +---------+                    +--------+
```
