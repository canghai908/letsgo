/*=============================================================================
#     FileName: chatclient.go
#         Desc: chat client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-13 17:48:50
#      History:
=============================================================================*/
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
    "strconv"
    "flag"

    "./protos"
    lnet "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/helper"
)

var endian = lnet.BigEndian

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *lnet.ClientPool) {
    reader := bufio.NewReader(os.Stdin)
    for {
        if (*cid)==0 {
            fmt.Print("you no connect anyone server,please input conn cmd,\n")
        }
        fmt.Print("you> ")
        input, _ := reader.ReadBytes('\n')
        cmd := string(input[:len(input)-1])

        var text string

        if cmd[0] == '/' {
            cmds := strings.Split(cmd," ")
            switch cmds[0]{
            case "/conn":
                var name,addr string
                if len(cmds)>2 {
                    name = cmds[1]
                    addr = cmds[2]
                }else {
                    name = "c_" + strconv.Itoa(*cid)
                    addr = cmds[1]
                }

                p := client.Clients.GetByName(name)
                if p != nil {
                    fmt.Println(name," is exists !")
                    continue
                }

                go client.Start(name,addr)


                fmt.Print("please input your name:")
                input, _ := reader.ReadBytes('\n')
                input =input[0:len(input)-1]

                for true {
                    b := client.Clients.GetByName(name)
                    if b!=nil{
                        change(cid,client,name)
                        break
                    }
                    time.Sleep(2*1e3)
                }

                text = string(input)

            case "/change":
                name := cmds[1]
                change(cid,client,name)

            case "/quit\n":
                text = "/quit"

            default:
                text = string(input[:len(input)-1])
            }
        } else {
            text = string(input[:len(input)-1])
        }

        msg := lnet.NewMessageWriter(endian)
        msg.SetCode(1011,0)
        msg.WriteString(text,0)

        log.Trace("has %v clients",client.Clients.Len())
        client.Clients.Get(*cid).SendMessage(msg)
    }
}

func change(cid *int,client *lnet.ClientPool,name string,) {
    b:= client.Clients.GetByName(name)
    if b!=nil{
        _cid := b.GetTransport().Cid
        *cid = _cid
        fmt.Println("current connection change:")
    }

    for c,p:=range client.Clients.All() {
        if p.GetName() != name {
            fmt.Println(" ",c,p.GetName())
        } else {
            fmt.Println("*",c,p.GetName())
        }
    }
}

var (
    loglevel = flag.Int("loglevel",0,"log level")
)

func main() {
    flag.Parse()

    log.SetLevel(*loglevel)

    datagram := lnet.NewDatagram(lnet.BigEndian)

    cid := 0
    client := lnet.NewClientPool(protos.MakeClient, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    running :=1
    for running==1 {
        time.Sleep(3*1e3)
    }
}
