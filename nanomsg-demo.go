package main

import (
    "fmt"
    "os"
)

const (
    delim = "=================================================="

    addressForwarder = "ipc:///tmp/nanomsg-demo-forwarder.ipc"
    addressServerA   = "ipc:///tmp/nanomsg-demo-server-a.ipc"
    addressServerB   = "ipc:///tmp/nanomsg-demo-server-b.ipc"
    addressServerC   = "ipc:///tmp/nanomsg-demo-server-c.ipc"
    addressPS        = "ipc:///tmp/nanomsg-demo-ps.ipc"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println(`Usage: ./nanomsg-demo.go TYPE

Where TYPE is one of these values:
    FRW - to start forwarder
    A, B or C - to start server A, B or C
    CLIENT - to start a client
    SUB - to subscribe for pub-sub messages`)
        os.Exit(1)
    }

    switch os.Args[1] {
    case "FRW":
        nodeForwarder()
    case "A", "B", "C":
        nodeServer(os.Args[1])
    case "CLIENT":
        nodeClient()
    case "SUB":
        nodeSubscriber()
    default:
        fmt.Println("Error. First argument should be one of these: IN, A, B, C, LOG")
        os.Exit(1)
    }
}
