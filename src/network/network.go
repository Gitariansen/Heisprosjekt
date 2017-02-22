package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

/*

ANDERS SE HER!!!
https://github.com/danielbmx/heisprosjekt/tree/master/src/networkmodule

*/

type UDPMessage struct {
	msg   string
	queue []bool
}

//var serverAddr *net.UDPAddr
//var localIP string
var port = ":30042"
var msg = "Hello from the other side"
var conn *net.UDPConn

func check_error(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

// All credz Anders.
// Brukes til å sammenligne om det er en melingbuffer[0:n],&save fra deg selv.
// I såfall ignorer denne.
func get_local_IP() string {
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
	check_error(err)
	defer conn.Close()

	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP
}

func Init() (*net.UDPAddr, string) {
	localIP := get_local_IP()

	// setting up UDP server for broadcasting
	serverAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255"+port)
	check_error(err)
	fmt.Println("Server adress: ", serverAddr)
	fmt.Println("Local adress: ", localIP)

	return serverAddr, localIP
}

func Recive_msg_UDP(msg_chan chan UDPMessage) {
	var msg UDPMessage

	serverAddr, err := net.ResolveUDPAddr("udp", port)
	check_error(err)

	conn, err := net.ListenUDP("udp", serverAddr)
	check_error(err)
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, address, err := conn.ReadFromUDP(buffer)
		check_error(err)
		fmt.Println("Got message from ", address, " with n = ", n)
		json.Unmarshal(buffer[0:n], &msg)
		check_error(err)
		msg_chan <- msg
	}
}

func Broadcast_UDP(serverAddr *net.UDPAddr) {
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	check_error(err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	check_error(err)
	defer conn.Close()

	fmt.Println("Sending message...")
	_, err = conn.Write([]byte(msg))
	check_error(err)
	msg = msg + "e"
	time.Sleep(1000 * time.Millisecond)
}
