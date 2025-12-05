package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pebbe/zmq4"
)

// local IP 조회
func get_local_ip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	host, err := os.Hostname()
	if err == nil {
		ip, _ := net.LookupIP(host)
		if len(ip) > 0 {
			return ip[0].String()
		}
	}
	return "127.0.0.1"
}

// SUB 소켓으로 1~254 대역 스캔
func search_nameserver(ipMask string, localIP string, port int) string {
	context, _ := zmq4.NewContext()
	defer context.Term()

	sub, _ := context.NewSocket(zmq4.SUB)
	defer sub.Close()

	sub.SetSubscribe("NAMESERVER")
	sub.SetRcvtimeo(2 * time.Second)

	for last := 1; last < 255; last++ {
		target := fmt.Sprintf("tcp://%s.%d:%d", ipMask, last, port)
		sub.Connect(target)
	}

	msg, err := sub.Recv(0)
	if err != nil {
		return ""
	}

	parts := strings.Split(msg, ":")
	if len(parts) == 2 && parts[0] == "NAMESERVER" {
		return parts[1]
	}
	return ""
}

// 1초마다 자신을 네임서버로 알림
func beacon_nameserver(localIP string, port int) {
	context, _ := zmq4.NewContext()
	pub, _ := context.NewSocket(zmq4.PUB)
	bindAddr := fmt.Sprintf("tcp://%s:%d", localIP, port)
	pub.Bind(bindAddr)
	fmt.Println("local p2p name server bind to", bindAddr)

	for {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("NAMESERVER:%s", localIP)
		pub.Send(msg, 0)
	}
}

// REQ/REP로 사용자 등록 처리
func user_manager_nameserver(localIP string, port int) {
	db := [][]string{}
	context, _ := zmq4.NewContext()
	rep, _ := context.NewSocket(zmq4.REP)
	bindAddr := fmt.Sprintf("tcp://%s:%d", localIP, port)
	rep.Bind(bindAddr)
	fmt.Println("local p2p db server activated at", bindAddr)

	for {
		req, err := rep.Recv(0)
		if err != nil {
			continue
		}
		parts := strings.Split(req, ":")
		db = append(db, parts)
		fmt.Printf("user registration '%s' from '%s'.\n", parts[1], parts[0])
		rep.Send("ok", 0)
	}
}

// PULL -> PUB 중계
func relay_server_nameserver(localIP string, pubPort int, pullPort int) {
	context, _ := zmq4.NewContext()

	pub, _ := context.NewSocket(zmq4.PUB)
	pub.Bind(fmt.Sprintf("tcp://%s:%d", localIP, pubPort))

	pull, _ := context.NewSocket(zmq4.PULL)
	pull.Bind(fmt.Sprintf("tcp://%s:%d", localIP, pullPort))

	fmt.Printf("local p2p relay server activated at tcp://%s:%d & %d.\n",
		localIP, pubPort, pullPort)

	for {
		msg, err := pull.Recv(0)
		if err != nil {
			continue
		}
		fmt.Println("p2p-relay:<==>", msg)
		pub.Send("RELAY:"+msg, 0)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run dechat.go _user-name_")
		return
	}

	userName := os.Args[1]

	portNameServer := 9001
	portChatPublisher := 9002
	portChatCollector := 9003
	portSubscribe := 9004

	localIP := get_local_ip()
	ipMask := localIP[:strings.LastIndex(localIP, ".")]

	fmt.Println("starting p2p chatting program.")
	fmt.Println("searching for p2p server.")

	// 네임서버 탐색
	foundServer := search_nameserver(ipMask, localIP, portNameServer)
	var serverIP string

	if foundServer == "" {
		// 스스로 서버 역할
		serverIP = localIP
		fmt.Println("p2p server is not found, activating server mode.")

		go beacon_nameserver(localIP, portNameServer)
		fmt.Println("p2p beacon server activated.")

		go user_manager_nameserver(localIP, portSubscribe)
		fmt.Println("p2p subscriber database server activated.")

		go relay_server_nameserver(localIP, portChatPublisher, portChatCollector)
		fmt.Println("p2p message relay server activated.")

	} else {
		serverIP = foundServer
		fmt.Printf("p2p server found at %s, client mode activated.\n", serverIP)
	}

	// 사용자 등록
	fmt.Println("starting user registration procedure.")
	ctxDB, _ := zmq4.NewContext()
	dbReq, _ := ctxDB.NewSocket(zmq4.REQ)
	dbReq.Connect(fmt.Sprintf("tcp://%s:%d", serverIP, portSubscribe))

	dbReq.Send(fmt.Sprintf("%s:%s", localIP, userName), 0)

	ack, _ := dbReq.Recv(0)
	if ack == "ok" {
		fmt.Println("user registration to p2p server completed.")
	} else {
		fmt.Println("user registration failed.")
	}

	fmt.Println("starting message transfer procedure.")

	// 메시지 send/recv 소켓
	ctxRelay, _ := zmq4.NewContext()

	sub, _ := ctxRelay.NewSocket(zmq4.SUB)
	sub.SetSubscribe("RELAY")
	sub.Connect(fmt.Sprintf("tcp://%s:%d", serverIP, portChatPublisher))

	push, _ := ctxRelay.NewSocket(zmq4.PUSH)
	push.Connect(fmt.Sprintf("tcp://%s:%d", serverIP, portChatCollector))

	fmt.Println("starting autonomous message transmit/receive scenario.")

	for {
		sub.SetRcvtimeo(100 * time.Millisecond)
		msg, err := sub.Recv(0)
		if err == nil {
			parts := strings.Split(msg, ":")
			fmt.Printf("p2p-recv::<<== %s:%s\n", parts[1], parts[2])
		} else {
			r := rand.Intn(100)
			if r < 10 {
				time.Sleep(3 * time.Second)
				out := fmt.Sprintf("(%s,%s:ON)", userName, localIP)
				push.Send(out, 0)
				fmt.Println("p2p-send::==>>", out)

			} else if r > 90 {
				time.Sleep(3 * time.Second)
				out := fmt.Sprintf("(%s,%s:OFF)", userName, localIP)
				push.Send(out, 0)
				fmt.Println("p2p-send::==>>", out)
			}
		}
	}
}
