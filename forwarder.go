package main

import (
    "encoding/json"
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
)

func bindClientSocket() *nanomsg.RepSocket {
    // Create a REP socket
    //socket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.REP)
    socket, err := nanomsg.NewRepSocket()
    if err != nil {
        log.Fatalln("Error while creating client socket.")
    }

    // Bind the socket
    _, err = socket.Bind(addressForwarder)
    if err != nil {
        log.Fatalln("Error while binding client socket to address.")
    }

    return socket
}

func connectServerSocket(num, address string) *nanomsg.ReqSocket {
    // Create REQ socket (to A, B, C)
    socket, err := nanomsg.NewReqSocket()
    if err != nil {
        log.Fatalf("Error while creating server %s socket.\n", num)
    }

    // Connect REQ socket
    _, err = socket.Connect(address)
    if err != nil {
        log.Fatalln("Error while connecting to address.")
    }

    return socket
}

func connectPSSocket() *nanomsg.PubSocket {
    // Create PUB socket
    socket, err := nanomsg.NewPubSocket()
    if err != nil {
        log.Fatalln("Error while creating PUB socket.")
    }

    // Bind PUB socket
    _, err = socket.Bind(addressPS)
    if err != nil {
        log.Fatalln("Error while binding socket to address.")
    }

    return socket
}

func nodeForwarder() {
    // Create REP socket (for client)
    socketClient := bindClientSocket()
    defer socketClient.Close()

    // Create sockets and connect to servers A, B and C
    socketServer := make(map[string]*nanomsg.ReqSocket)
    socketServer["A"] = connectServerSocket("A", addressServerA)
    socketServer["B"] = connectServerSocket("B", addressServerB)
    socketServer["C"] = connectServerSocket("C", addressServerC)
    defer socketServer["A"].Close()
    defer socketServer["B"].Close()
    defer socketServer["C"].Close()

    // Create PUB socket
    socketPS := connectPSSocket()
    defer socketPS.Close()

    log.Println("Ok, waiting for something...")
    fmt.Println(delim)

    for {
        // Read REQ from client in form of:
        // {"route":"A"/"B"/"C", "data": "Ping A/B/C"}
        msg, err := socketClient.Recv(0)
        if err != nil {
            log.Println("Error while receiving message from client.")
            log.Println(err)
            continue
        }
        var msgJSON map[string]string
        err = json.Unmarshal(msg, &msgJSON)
        dest := msgJSON["route"]

        // Send REQ to server in form of:
        // {"route":"A"/"B"/"C", "data": "Ping A/B/C"}
        _, err = socketServer[dest].Send(msg, 0)
        if err != nil {
            log.Printf("Error while sending message to server %s.\n", dest)
            continue
        }

        // Send PUB to subscribers in form of:
        // client -> server A/B/C: {...}
        msg = []byte(fmt.Sprintf("client -> server %s:\n%s", dest, msg))
        _, err = socketPS.Send(msg, 0)
        if err != nil {
            log.Println("Error while sending PUB message.")
            continue
        }

        // Read REP from server in form of:
        // {"data": "Pong A/B/C"}
        msg, err = socketServer[dest].Recv(0)
        if err != nil {
            log.Println("Error while receiving message from server.")
            continue
        }

        // Send REP to client in form of:
        // {"data": "Pong A/B/C"}
        _, err = socketClient.Send(msg, 0)
        if err != nil {
            log.Println("Error while sending message to client.")
            continue
        }

        // Send PUB to subscribers in form of:
        // server A/B/C -> client: {...}
        msg = []byte(fmt.Sprintf("server %s -> client:\n%s", dest, msg))
        _, err = socketPS.Send(msg, 0)
        if err != nil {
            log.Println("Error while sending PUB message.")
            continue
        }
    }
}
