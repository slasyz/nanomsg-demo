package main

import (
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
)

func nodeServer(num string) {
    // Create REP socket
    //socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.REP)
    socket, err := nanomsg.NewRepSocket()
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    // Choose address
    var address string
    switch num {
    case "A":
        address = addressServerA
    case "B":
        address = addressServerB
    case "C":
        address = addressServerC
    }

    // Bind socket to address.
    _, err = socket.Bind(address)
    if err != nil {
        log.Fatalln("Error while binding socket to address.")
    }

    log.Println("Ok, waiting for something...")
    fmt.Println(delim)

    for {
        // Read REQ from forwarder in form of:
        // {"data": "Ping A/B/C"}
        msg, err := socket.Recv(0)
        if err != nil {
            log.Println("Error while receiving message.")
            continue
        }

        // Output it to stdout
        log.Println("Got a message:")
        fmt.Println(string(msg))

        // Send REP to forwarder in form of:
        // {"data": "Pong A/B/C"}
        msg = []byte(fmt.Sprintf("{\"data\": \"Pong %s\"}", num))
        _, err = socket.Send(msg, 0)
        if err != nil {
            log.Println("Oops, reply hasn't been sent.")
            continue
        }
        log.Println("Ok, reply has been sent.")
        fmt.Println(delim)
    }

}
