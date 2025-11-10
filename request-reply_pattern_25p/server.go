package main

import (
	"fmt" // print()
	"time" // 시간 관련 함수(sleep()) 
	"github.com/pebbe/zmq4" // zmq
)

func main() {

	// 소켓 생성
	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil { // err == nil -> 에러 없음
		panic(err)
	}

	// 바인드
	err = socket.Bind("tcp://*:5555")
	if err != nil {
		panic(err)
	}

	for {
		// wait for next request from client
		msg, err := socket.Recv(0) // 인자 0은 기본 blocking: 메시지가 올 때까지 멈춤
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received request: %s\n", msg)
	}
}