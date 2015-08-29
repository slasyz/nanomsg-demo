package main

import (
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
)

func nodeOutput(num, addr string) {
    socket, err := nanomsg.NewSubSocket()
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    err = socket.Subscribe(num)
    if err != nil {
        log.Fatalf("Error while subscribing to %s topic.\n", num)
    }

    _, err = socket.Connect(addr)
    if err != nil {
        log.Fatalln("Error while connecting to address.")
    }

    log.Println("Ok, waiting for something...")
    fmt.Println(split)

    for {
        msg, err := socket.Recv(0)
        if err != nil {
            log.Println("Error while receiving message.")
            continue
        }

        log.Println("Got a letter for you!")
        fmt.Println(string(msg[1:]))
        fmt.Println(split)
    }
}
