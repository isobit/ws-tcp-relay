# ws-tcp-relay
A relay between Websocket and TCP.

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
