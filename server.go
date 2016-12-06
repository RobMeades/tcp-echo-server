/* Entry point for TCP echo server.
 *
 * This code based:
 * https://gist.github.com/paulsmith/775764
 */
 
package main

import (
    "net"
    "bufio"
    "fmt"
    "os"
    "flag"
    "log"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// File handle
var pWriter *bufio.Writer
var numConnections int
var numPackets int

// Command-line flags
var pPort = flag.String ("p", "", "the port number to listen on.");
var Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s: run the TCP echo server.  Usage:\n", os.Args[0])
        flag.PrintDefaults()
    }

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Handle a connection
func clientConnection(listener net.Listener) chan net.Conn {
    channel := make(chan net.Conn)
    
    go func() {
        for {
            client, err := listener.Accept()
            if (client != nil) {
                numPackets = 0;
                numConnections++
                log.Printf("Connection %d: %v <-> %v.\n", numConnections, client.LocalAddr(), client.RemoteAddr())
                channel <- client
            } else {
                log.Printf("Couldn't accept connection (%s).\n", err.Error())
            }
        }
    }()
    
    return channel
}

func handleConnection(client net.Conn) {
    var err error
    line := make([]byte, 1024)

    b := bufio.NewReader(client)
    for err == nil { // err will be set on EOF when the client closes
        line, err = b.ReadBytes('\n')
        if err == nil {
            numPackets++;
            num, err1:= client.Write(line)
            log.Printf("%d.%d: %v <-> %v: %q", numConnections, numPackets, client.LocalAddr(), client.RemoteAddr(), line)
            if err1 != nil {
                log.Printf("Error, only %d of %d byte(s) could be echoed (%s).\n", num, len(line), err1.Error())    
            }
        }
    }
    log.Printf("Connection dropped (%s).\n", err.Error())    
}

// Entry point
func main() {
    // Deal with the command-line parameters
    flag.Parse()
    
    // Set up logging
    log.SetFlags(log.LstdFlags)
    
    if *pPort != "" {
        // Say what we're doing
        fmt.Printf("Echoing TCP packets received on port %s.\n", *pPort)
        
        // Set up the server
        pServer, err := net.Listen("tcp", ":" + *pPort)
        if pServer != nil {
            connection := clientConnection(pServer)
            for {
                go handleConnection(<-connection)
            }
        } else {
            fmt.Printf("Couldn't start TCP listening server on port %s (%s).\n", *pPort, err.Error())
        }            
    } else {
        fmt.Printf("Must specify a port number.\n")
        flag.PrintDefaults()
        os.Exit(-1);
    }
}

// End Of File
