package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

var tcpAddress string

func copyWorker(dst io.Writer, src io.Reader, doneCh chan<- bool) {
	io.Copy(dst, src)
	doneCh <- true
}

func relayHandler(ws *websocket.Conn) {
	conn, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		log.Printf("[ERROR] %v \n", err)
		return
	}

	doneCh := make(chan bool)

	go copyWorker(conn, ws, doneCh)
	go copyWorker(ws, conn, doneCh)

	<-doneCh
	conn.Close()
	ws.Close()
	<-doneCh
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <tcpTargetAddress>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var port uint
	var certFile string
	var keyFile string

	flag.UintVar(&port, "p", 4223, "The port to listen on")
	flag.UintVar(&port, "port", 4223, "The port to listen on")
	flag.StringVar(&certFile, "tlscert", "", "TLS cert file path")
	flag.StringVar(&keyFile, "tlskey", "", "TLS key file path")
	flag.Usage = usage
	flag.Parse()

	tcpAddress = flag.Arg(0)
	if tcpAddress == "" {
		fmt.Fprintln(os.Stderr, "No address specified")
		os.Exit(1)
	}

	portString := fmt.Sprintf(":%d", port)

	log.Printf("[INFO] Listening on %s\n", portString)

	http.Handle("/", websocket.Handler(relayHandler))

	var err error
	if certFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(portString, certFile, keyFile, nil)
	} else {
		err = http.ListenAndServe(portString, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
