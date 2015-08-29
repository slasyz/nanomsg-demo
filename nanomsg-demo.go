package main

import (
    "fmt"
    "os"
)

const (
    pubSubAddress = "ipc:///tmp/nanomsg-demo-ps.ipc"
    logAddress    = "ipc:///tmp/nanomsg-demo-log.ipc"
    split         = "=================================================="
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println(`Usage: ./nanomsg-demo.go TYPE

Where TYPE is one of these values: IN, A, B, C, LOG`)
        os.Exit(1)
    }

    switch os.Args[1] {
    case "IN":
        nodeInput(pubSubAddress, logAddress)
    case "A", "B", "C":
        nodeOutput(os.Args[1], pubSubAddress)
    case "LOG":
        nodeLog(logAddress)
    default:
        fmt.Println("Error. First argument should be one of these: IN, A, B, C, LOG")
        os.Exit(1)
    }
}
