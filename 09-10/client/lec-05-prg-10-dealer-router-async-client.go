package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pebbe/zmq4"
)

type ClientTask struct {
	id string
}

func (c *ClientTask) Run() {

	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	// DEALER 소켓 생성
	socket, err := ctx.NewSocket(zmq4.DEALER)
	if err != nil {
		panic(err)
	}

	// 클라이언트 식별자
	identity := c.id
	socket.SetIdentity(identity)

	// 연결
	err = socket.Connect("tcp://localhost:5570")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client %s started\n", identity)

	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)

	reqs := 0

	for {

		reqs++
		fmt.Printf("Req #%d sent..\n", reqs)

		// 서버로 메시지 전송
		_, err = socket.Send(fmt.Sprintf("request #%d", reqs), 0)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)

		sockets, err := poller.Poll(1000 * time.Millisecond)
		if err != nil {
			panic(err)
		}

		if len(sockets) > 0 {
			msg, err := socket.Recv(0)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s received: %s\n", identity, msg)
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		panic("usage: go run main.go <client_id>")
	}

	client := ClientTask{id: os.Args[1]}

	go client.Run()

	select {}
}
