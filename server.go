/* Entry point for echo server.
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
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// File handle
var pFile *os.File

// Command-line flags
var pPort = flag.String ("p", "", "the port number to listen on.");
var pFileName = flag.String ("f", "", "the file name to write the receive requests to.");
var Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s: run the echo server.  Usage:\n", os.Args[0])
        flag.PrintDefaults()
    }

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Handle a connection
func clientConnection(listener net.Listener) chan net.Conn {
    channel := make(chan net.Conn)
    i := 0
    
    go func() {
        for {
            client, err := listener.Accept()
            if (client != nil) {
                i++
                fmt.Printf("%d: %v <-> %v.\n", i, client.LocalAddr(), client.RemoteAddr())
                channel <- client
            } else {
                fmt.Printf("Couldn't accept connection (%s).\n", err.Error())
            }
        }
    }()
    
    return channel
}

func handleConnection(client net.Conn) {
    var err error 
    var line []byte
    
    b := bufio.NewReader(client)
    for err == nil { // err will be set on EOF when the client closes
        line, err = b.ReadBytes('\n')
        if err == nil {
            num, err1:= client.Write(line)
            if err1 != nil {
                fmt.Printf("Error, only %d of %d byte(s) could be echoed (%s).\n", num, len(line), err1.Error())    
            }
            if pFile != nil {
                pFile.Write(line)
            }
        }
    }
    fmt.Printf("Connection dropped (%s).\n", err.Error())    
}

// Entry point
func main() {
    var err error

    // Deal with the command-line parameters
    flag.Parse()
    
    if *pPort != "" {
        // Open the output file for append
        if *pFileName != "" {
            pFile, err = os.OpenFile(*pFileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
        }    
        
        if err == nil {        
            // Set up the echo server
            fmt.Printf("Echoing packets received on port %s", *pPort)
            if (pFile != nil) {
                fmt.Printf(" and writing received packets to \"%s\"", pFile.Name())
            }
            fmt.Printf(".\n")
            
            pServer, err := net.Listen("tcp", ":" + *pPort)
            if pServer != nil {
                connection := clientConnection(pServer)
                for {
                    go handleConnection(<-connection)
                }
            } else {
                fmt.Printf("Couldn't start listening server on port %s (%s).\n", *pPort, err.Error())
            }
            
        } else {
            fmt.Printf("Couldn't open file %s (%s).\n", *pFileName, err.Error())
            os.Exit(-1);
        }
    } else {
        fmt.Printf("Must specify a port number.\n")
        flag.PrintDefaults()
        os.Exit(-1);
    }
}

// End Of File
