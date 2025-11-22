package main

import (
	"fmt" // print()
	"strings"

	"github.com/pebbe/zmq4" // zmq

	"os"
	"strconv"
)

func main() {

	// 컨텍스트 생성
	ctx, err := zmq4.NewContext()
	if err != nil {
		panic(err)
	}

	//SUB 소켓 생성
	socket, err := ctx.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Collecting updates from weather server...")

	// 연결
	err = socket.Connect("tcp://localhost:5556")
	if err != nil {
		panic(err)
	}

	zip_filter := "10001"
	if len(os.Args) > 1 {
		zip_filter = os.Args[1]
	}

	socket.SetSubscribe(zip_filter)

	total_temp := 0

	var update_nbr int

	for update_nbr = 0; update_nbr < 20; update_nbr++ {
		msg, err := socket.Recv(0)
		if err != nil {
			panic(err)
		}
		parts := strings.Split(msg, " ")
		zipcode := parts[0]
		temperature := parts[1]
		relhumidity := parts[2]

		_ = zipcode // zipcode와 relhumidity는 선언되었지만 사용되지 않아서 _로 버림
		_ = relhumidity

		temp, err := strconv.Atoi(temperature)
		if err != nil {
			panic(err)
		}
		total_temp += temp

		fmt.Printf("Received temperature for zipcode '%s' was %s F\n", zip_filter, temperature)
	}

	avg := float64(total_temp) / float64(update_nbr)
	fmt.Printf("Average temperature for zipcode '%s' was %.2f F\n", zip_filter, avg)

}
