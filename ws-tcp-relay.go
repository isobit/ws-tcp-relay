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
		//log.Printf("[ERROR] TCP connection to %s failed\n", tcpAddress)
		return
	}
	//log.Printf("[INFO] accepted connection with %s\n", conn.RemoteAddr())

	doneCh := make(chan bool)

	go copyWorker(conn, ws, doneCh)
	go copyWorker(ws, conn, doneCh)

	<-doneCh
	conn.Close()
	ws.Close()
	<-doneCh

	//log.Printf("[INFO] closed connection with %s\n", conn.RemoteAddr())
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <tcpTargetAddress>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)

	var port int
	var certFile string
	var keyFile string
	flag.IntVar(&port, "port", 1337, "Port to listen on.")
	flag.StringVar(&certFile, "cert", "", "TLS cert file path")
	flag.StringVar(&keyFile, "key", "", "TLS key file path")
	flag.Usage = usage;
	flag.Parse();
	tcpAddress = flag.Arg(0)
	if tcpAddress == "" {
		log.Fatal("no address specified")
	}

	log.Printf("[INFO] listening on port %d\n", port)
	http.Handle("/", websocket.Handler(relayHandler))
	var err error
	if certFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, nil)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
