# ws-tcp-relay
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
A relay between Websocket and TCP. All messages will be copied from all 
Websocket connections to the target TCP server, and vice-versa.

## Installation
```go get -u github.com/joshglendenning/ws-tcp-relay```

## Usage
```
Usage: ws-tcp-relay <tcpTargetAddress>
  -p int
        Port to listen on. (default 1337)
  -port int
        Port to listen on. (default 1337)
  -tlscert string
        TLS cert file path
  -tlskey string
        TLS key file path
```

## WSS Support
To use secure websockets simply specify both the `tlscert` and `tlskey` flags.

## Building
`go build ws-tcp-relay.go`
