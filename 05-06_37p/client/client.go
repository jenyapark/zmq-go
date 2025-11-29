package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pebbe/zmq4"
)

func main() {
	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	subscriber, err := ctx.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}

	err = subscriber.SetSubscribe("")
	if err != nil {
		panic(err)
	}

	err = subscriber.Connect("tcp://localhost:5556")
	if err != nil {
		panic(err)
	}

	publisher, err := ctx.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}

	err = publisher.Connect("tcp://localhost:5558")
	if err != nil {
		panic(err)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	poller := zmq4.NewPoller()
	poller.Add(subscriber, zmq4.POLLIN)

	for {
		// a = subscriber.poll(100)
		a, err := poller.Poll(100 * time.Millisecond)
		if err != nil {
			panic(err)
		}

		if len(a) > 0 {
			msg, _ := subscriber.RecvBytes(0)
			fmt.Println("I: received message", msg)
		} else {
			r := rnd.Intn(100) + 1
			if r < 10 {
				data := []byte(fmt.Sprintf("%d", r))
				publisher.SendBytes(data, 0)
				fmt.Println("I: sending message", r)
			}
		}
	}
}
