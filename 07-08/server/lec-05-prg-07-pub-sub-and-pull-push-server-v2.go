package main

import (
	"fmt"

	"github.com/pebbe/zmq4"
)

func main() {
	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	// PUB 소켓 생성
	publisher, err := ctx.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}

	err = publisher.Bind("tcp://*:5557")
	if err != nil {
		panic(err)
	}

	// PULL 소켓 생성
	collector, err := ctx.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}

	err = collector.Bind("tcp://*:5558")
	if err != nil {
		panic(err)
	}

	// 메시지 계속 수신 및 재전송
	for {
		// 메시지 수신
		message, err := collector.RecvBytes(0)
		if err != nil {
			panic(err)
		}

		// 출력
		fmt.Println("server: publishing update =>", string(message))

		// PUB 송신
		_, err = publisher.SendBytes(message, 0) // 전송한 바이트 수는 필요 없으므로 버림
		if err != nil {
			panic(err)
		}
	}

}
