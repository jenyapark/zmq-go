package main

import (
	"fmt"  // print()
	"time" // 시간 관련 함수(sleep())

	"github.com/pebbe/zmq4" // zmq
)

func main() {

	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil { // err == nil -> 에러 없음
		panic(err)
	}

	// REP 소켓 생성
	socket, err := ctx.NewSocket(zmq4.REP)
	if err != nil {
		panic(err)
	}

	// 바인드
	err = socket.Bind("tcp://*:5555")
	if err != nil {
		panic(err)
	}

	for {
		// 메시지 수신
		msg, err := socket.Recv(0) // 인자 0은 기본 blocking: 메시지가 올 때까지 멈춤
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received request: %s\n", msg)

		// 1초 대기
		time.Sleep(1 * time.Second)

		// 응답 전송
		_, err = socket.Send("World", 0)
		if err != nil {
			panic(err)
		}
	}
}
