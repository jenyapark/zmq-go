package main

import (
	"fmt" // print()

	"github.com/pebbe/zmq4" // zmq
)

func main() {

	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil { // err == nil -> 에러 없음
		panic(err)
	}

	fmt.Println("Connecting to hello world server...")

	// REP 소켓 생성
	socket, err := ctx.NewSocket(zmq4.REQ)
	if err != nil {
		panic(err)
	}

	// 서버에 연결
	err = socket.Connect("tcp://localhost:5555")
	if err != nil {
		panic(err)
	}

	// Do 10 requests, waiting each time for a response
	for request := 0; request < 10; request++ {
		fmt.Printf("Sending request %d ...\n", request)

		_, err = socket.Send("Hello", 0) // go는 원래 바이트 시퀀스
		if err != nil {
			panic(err)
		}

		message, err := socket.Recv(0)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received reply %d [ %s ]\n", request, message)

	}

}
