package main

import (
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
)

func nodeLog(addrLog string) {
    socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.REP)
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    _, err = socket.Bind(addrLog)
    if err != nil {
        log.Fatalln("Error while binding socket to address.")
    }

    log.Println("Ok, waiting for something...")
    fmt.Println(split)

    for {
        msg, err := socket.Recv(0)
        if err != nil {
            log.Println("Error while receiving message.")
        }

        fmt.Println(string(msg))
        fmt.Println(split)
    }
}
