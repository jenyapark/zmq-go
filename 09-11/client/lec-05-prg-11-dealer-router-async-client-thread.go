package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pebbe/zmq4"
)

type ClientTask struct {
	id       string
	identity string
	socket   *zmq4.Socket
	poller   *zmq4.Poller
}

func (c *ClientTask) recvHandler() {
	for {
		sockets, _ := c.poller.Poll(1000 * time.Millisecond)
		for _, s := range sockets {
			if s.Socket == c.socket {
				msg, _ := c.socket.Recv(0)
				fmt.Printf("%s received: %s\n", c.identity, msg)
			}
		}
	}
}

func (c *ClientTask) run() {
	//DEALER 소켓 생성 + ID 설정 + 서버에 연결
	context, _ := zmq4.NewContext()
	socket, _ := context.NewSocket(zmq4.DEALER)
	c.socket = socket
	c.identity = c.id
	c.socket.SetIdentity(c.identity)
	c.socket.Connect("tcp://localhost:5570")

	fmt.Printf("Client %s started\n", c.identity)

	// 폴러
	c.poller = zmq4.NewPoller()
	c.poller.Add(c.socket, zmq4.POLLIN)

	// 수신 전용 고루틴 시작
	go c.recvHandler()

	reqs := 0
	for {
		reqs++
		fmt.Printf("Req #%d sent..\n", reqs)
		c.socket.Send(fmt.Sprintf("request #%d", reqs), 0)
		time.Sleep(1 * time.Second)
	}

	// useless
	c.socket.Close()
	context.Term()
}

func NewClientTask(id string) *ClientTask {
	return &ClientTask{id: id}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dealer_router_async_client.go <client_id>")
		return
	}

	client := NewClientTask(os.Args[1])
	client.run()
}
