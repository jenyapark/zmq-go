package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pebbe/zmq4" // zmq
)

type ServerTask struct {
	num_server int
}

type ServerWoker struct {
	ctx *zmq4.Context
	id  int
}

func (s *ServerTask) Run() {

	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	// ROUTER 소켓 생성
	frontend, err := ctx.NewSocket(zmq4.ROUTER)
	if err != nil {
		panic(err)
	}

	// 서버에 연결
	err = frontend.Bind("tcp://*:5570")
	if err != nil {
		panic(err)
	}

	// DEALER 소켓 생성
	backend, err := ctx.NewSocket(zmq4.DEALER)
	if err != nil {
		panic(err)
	}

	// 워커와 연결될 inproc 엔드포인트
	err = backend.Bind("inproc://backend")
	if err != nil {
		panic(err)
	}

	// 워커 스레드 생성
	var workers []*ServerWoker
	for i := 0; i < s.num_server; i++ {
		worker := &ServerWoker{ctx: ctx, id: i}
		go worker.Run()
		workers = append(workers, worker)
	}

	// router <-> dealer 자동 중계
	err = zmq4.Proxy(frontend, backend, nil)
	if err != nil {
		panic(err)
	}

	frontend.Close()
	backend.Close()
	ctx.Term()
}

// 각 워커 스레드 실행
func (w *ServerWoker) Run() {

	//워커용 DEALER 소켓 생성
	worker, err := w.ctx.NewSocket(zmq4.DEALER)
	if err != nil {
		panic(err)
	}

	// 서버의 backend에 연결
	err = worker.Connect("inproc://backend")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Worker#%d started\n", w.id)

	for {
		parts, err := worker.RecvMessageBytes(0)
		if err != nil {
			panic(err)
		}

		if len(parts) < 2 {
			continue
		}

		ident := parts[0]
		msg := parts[1]

		fmt.Printf("Worker#%d received %s from %s\n", w.id, string(msg), string(ident))

		_, err = worker.SendMessage(ident, msg)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	// 워커 개수
	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	server := &ServerTask{num_server: num}
	go server.Run()
	select {}
}
