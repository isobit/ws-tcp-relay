package main

import (
	"fmt"
	"log"
	"flag"
	"io"
	"os"
	"net"
	"net/http"
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
	var port int
	var certFile string
	var keyFile string

	flag.IntVar(&port, "p", 4223, "Port to listen on.")
	flag.IntVar(&port, "port", 4223, "Port to listen on.")
	flag.StringVar(&certFile, "tlscert", "", "TLS cert file path")
	flag.StringVar(&keyFile, "tlskey", "", "TLS key file path")
	flag.Usage = usage;
	flag.Parse();

	tcpAddress = flag.Arg(0)
	if tcpAddress == "" {
		fmt.Fprintln(os.Stderr, "no address specified")
		os.Exit(1)
	}

	portString := fmt.Sprintf(":%d", port)

	log.Printf("[INFO] starting server on port %d\n", port)

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
