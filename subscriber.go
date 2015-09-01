package main

import (
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
)

func nodeSubscriber() {
    // Create SUB socket
    socket, err := nanomsg.NewSubSocket()
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    // Subscribe to ""
    err = socket.Subscribe("")
    if err != nil {
        log.Fatalln("Error while subscribing to \"\" topic.")
    }

    // Connect to socket.
    _, err = socket.Connect(addressPS)
    if err != nil {
        log.Fatalln("Error while connecting to address.")
    }

    log.Println("Ok, waiting for something...")
    fmt.Println(delim)

    // Read PUB from server
    for {
        msg, err := socket.Recv(0)
        if err != nil {
            log.Println("Error while receiving message.")
            continue
        }

        // Output it to stdout
        log.Println("Got a message:")
        fmt.Println(string(msg))
        fmt.Println(delim)
    }

}
