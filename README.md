# ws-tcp-relay
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)

An extremely simple relay/bridge/proxy between WebSocket clients and TCP servers. Data
received from WebSocket clients is simply forwarded to the specified TCP
server, and vice-versa. In other words, it's
[websocketd](https://github.com/joewalnes/websocketd), but for TCP connections
instead of `STDIN` and `STDOUT`.

## Usage
```
Usage: ws-tcp-relay <tcpTargetAddress>
  -b	Use binary frames instead of text frames
  -binary
    	Use binary frames instead of text frames
  -p uint
    	The port to listen on (default 4223)
  -port uint
    	The port to listen on (default 4223)
  -tlscert string
    	TLS cert file path
  -tlskey string
    	TLS key file path
```

### Binary Data
By default, `golang.org/x/net/websocket` uses text frames to deliver payload
data. To use binary frames instead, use either the `b` or `binary` flags.

### WSS Support
To use secure WebSockets simply specify both the `tlscert` and `tlskey` flags.

## Installation
```
go get -u github.com/isobit/ws-tcp-relay
```

Binaries are also available on the [release page](https://github.com/isobit/ws-tcp-relay/releases/).
