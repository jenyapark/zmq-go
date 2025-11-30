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

	publisher, err := ctx.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}

	err = publisher.Bind("tcp://*:5556")
	if err != nil {
		panic(err)
	}

	collector, err := ctx.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}

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
