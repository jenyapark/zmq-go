package main

import (
	"fmt" // print()

	"github.com/pebbe/zmq4" // zmq
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

	// 바인드
	err = publisher.Bind("tcp://*:5557")
	if err != nil {
		panic(err)
	}

	// PULL 소켓 생성
	collector, err := ctx.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}

	// 바인드
	err = collector.Bind("tcp://*:5558")
	if err != nil {
		panic(err)
	}

	for {

		message, err := collector.RecvBytes(0)
		if err != nil {
			panic(err)
		}

		fmt.Println("I: publishing update", message)

		publisher.SendBytes(message, 0)

	}

}
