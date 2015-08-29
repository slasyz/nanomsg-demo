package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/op/go-nanomsg"
    "log"
    "os"
)

func createLogReqSocket(logAddr string) *nanomsg.Socket {
    logSocket, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.REQ)
    if err != nil {
        log.Fatalln("Error while creating log socket.")
    }

    _, err = logSocket.Connect(logAddr)
    if err != nil {
        log.Fatalln("Error while binding log socket to address.")
    }

    return logSocket
}

func readJSONData(r *bufio.Reader) []byte {
    var inputString []byte

    piece, _ := r.ReadBytes(byte('\n'))
    for len(piece) > 1 {
        inputString = append(inputString, piece...)
        piece, _ = r.ReadBytes(byte('\n'))
    }

    return bytes.Trim(inputString, "\n\t ")
}

func nodeInput(addr, logAddr string) {
    // Read stdin, send it to A, B or C, and then to LOG
    socket, err := nanomsg.NewPubSocket()
    if err != nil {
        log.Fatalln("Error while creating socket.")
    }
    defer socket.Close()

    _, err = socket.Bind(addr)
    if err != nil {
        log.Fatalln("Error while binding socket to address.")
    }

    // Create log socket
    logSocket := createLogReqSocket(logAddr)
    defer logSocket.Close()

    log.Println("Now you can enter here JSON messages (type two new lines in the end of each message). Enter empty message to exit.")
    fmt.Println(split)

    reader := bufio.NewReader(os.Stdin)
    for {
        // Read another message
        data := readJSONData(reader)
        if len(data) == 0 {
            break
        }

        // Convert it to JSON to read "route" value
        v := make(map[string]interface{})
        err := json.Unmarshal(data, &v)
        if err != nil {
            log.Println("Error while parsing JSON.")
            continue
        }

        // Send "A{...}", "B{...}" or "C{...}" message
        if val, ok := v["route"]; ok {
            dest := []byte(val.(string))
            msg := append(dest, data...)
            socket.Send(msg, 0)
            log.Println("Ok, message sent.")
            fmt.Println(split)
        } else {
            log.Println("Error: there must be \"route\" key in JSON array.")
        }

        // send data to LOG
        _, err = logSocket.Send(data, 0)
        if err != nil {
            log.Println("Error while logging this.")
        }
    }
}
