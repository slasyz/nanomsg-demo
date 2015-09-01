package main

import (
    "bufio"
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
    "os"
)

func nodeClient() {
    // Create REQ socket
    //socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.REQ)
    socket, err := nanomsg.NewReqSocket()
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    // Connect to socket
    _, err = socket.Connect(addressForwarder)
    if err != nil {
        log.Fatalln("Error while connecing to address.")
    }

    reader := bufio.NewReader(os.Stdin)
    log.Println("Now you can enter \"A\", \"B\" or \"C\" to ping corresponding server.")
    fmt.Println(delim)

    for {
        // Wait for input from stdin
        num, _ := reader.ReadString(byte('\n'))
        num = num[:1]

        // Send REQ to forwarder in form of:
        // {"route":"A"/"B"/"C", "data": "Ping A/B/C"}
        msg := []byte(fmt.Sprintf(`{"route": "%s", "data": "Ping %s"}`, num, num))
        _, err = socket.Send(msg, 0)
        if err != nil {
            log.Println("Oops, message hasn't been sent.")
            continue
        }
        log.Println("Ok, message has been sent.")

        // Wait REP from forwarder in form of:
        // {"data": "Ping A/B/C"}
        msg, err := socket.Recv(0)
        if err != nil {
            log.Println("Error while receiving message.")
            continue
        }

        // Output REP to stdout
        log.Println("Got a reply:")
        fmt.Println(string(msg))
        fmt.Println(delim)
    }
}
