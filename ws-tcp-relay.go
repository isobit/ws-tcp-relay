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

var (
	tcpAddress string
	certFile   string
	keyFile    string
	port       uint
	binaryMode bool
	fs         flag.FlagSet
)

func init() {
	fs = *flag.NewFlagSet("Commands", flag.ExitOnError)
	fs.UintVar(&port, "p", 4223, "The port to listen on")
	fs.UintVar(&port, "port", 4223, "The port to listen on")
	fs.StringVar(&certFile, "tlscert", "", "TLS cert file path")
	fs.StringVar(&keyFile, "tlskey", "", "TLS key file path")
	fs.BoolVar(&binaryMode, "b", false, "Use binary frames instead of text frames")
	fs.BoolVar(&binaryMode, "binary", false, "Use binary frames instead of text frames")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <tcpTargetAddress>\n", os.Args[0])
		fs.PrintDefaults()
	}

	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "No address specified")
		os.Exit(1)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stderr, "Usage: %s <tcpTargetAddress>\n", os.Args[0])
		fs.PrintDefaults()
		os.Exit(0)
	}
	fs.Parse(os.Args[2:])
}

func main() {
	portString := fmt.Sprintf(":%d", port)

	tcpAddress = os.Args[1]

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

	if binaryMode {
		ws.PayloadType = websocket.BinaryFrame
	}

	doneCh := make(chan bool)

	go copyWorker(conn, ws, doneCh)
	go copyWorker(ws, conn, doneCh)

	<-doneCh
	conn.Close()
	ws.Close()
	<-doneCh
}
