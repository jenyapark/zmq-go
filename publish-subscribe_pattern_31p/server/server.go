package main

import (
	"fmt" // print()

	"github.com/pebbe/zmq4" // zmq

	"math/rand"
)

func randrange(a, b int) int {
	return a + rand.Intn(b-a)
}

func main() {

	fmt.Printf("Publishing updates at weather server...")

	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	// PUB 소켓 생성
	socket, err := ctx.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}

	// 바인드
	err = socket.Bind("tcp://*:5556")
	if err != nil {
		panic(err)
	}

	for {
		zipcode := randrange(1, 100000)
		temperature := randrange(-80, 135)
		relhumidity := randrange(10, 60)
		msg := fmt.Sprintf("%d %d %d", zipcode, temperature, relhumidity)
		socket.Send(msg, 0)
	}
}
