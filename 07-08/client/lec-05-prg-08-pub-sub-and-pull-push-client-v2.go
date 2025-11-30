package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/pebbe/zmq4"
)

func main() {
	clientID := os.Args[1]

	// Context 생성
	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	// SUB 소켓 생성
	subscriber, err := ctx.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}

	// 모든 메시지 구독
	err = subscriber.SetSubscribe("")
	if err != nil {
		panic(err)
	}

	// SUB 연결
	err = subscriber.Connect("tcp://localhost:5557")
	if err != nil {
		panic(err)
	}

	// PUSH 소켓 생성
	publisher, err := ctx.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}

	// PUSH 연결
	err = publisher.Connect("tcp://localhost:5558")
	if err != nil {
		panic(err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	poller := zmq4.NewPoller()
	poller.Add(subscriber, zmq4.POLLIN)

	for {
		a, err := poller.Poll(100 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		// 메시지 수신
		if len(a) > 0 {
			msgBytes, err := subscriber.RecvBytes(0)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s: receive status => %s\n", clientID, string(msgBytes))
		} else { // 메시지 없을 때 랜덤 작업 수행함
			val := r.Intn(100) + 1

			if val < 10 {
				time.Sleep(1 * time.Second)
				msg := "(" + clientID + ":ON)"
				_, err := publisher.Send(msg, 0)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s: send status - activated\n", clientID)

			} else if val > 90 {
				time.Sleep(1 * time.Second)
				msg := "(" + clientID + ":OFF)"
				_, err := publisher.Send(msg, 0)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s: send status - deactivated\n", clientID)
			}
		}
	}

}
