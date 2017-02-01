package network


import (
	"fmt"
	"net"
  "os"
	"time"
	"strings"
)

var serverAddr *net.UDPAddr
var localIP string
var port = "30042"



func check_error(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(0)
    }
}


// All credz Anders
func get_local_IP() (string) {
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
	go check_error(err)
	defer conn.Close()

	localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP
}



func init() {
    // setting up UDP server for broadcasting
    serverAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort("255.255.255.255", port))
		localIP = get_local_IP()
    check_error(err)
		fmt.Println("Server adress: ", serverAddr)
		fmt.Println("Local adress: ", localIP)
}


func recive_msg_UDP(port string) string {
	conn, err := net.ListenUDP("udp", serverAddr)
	check_error(err)
	defer conn.Close()

	var buffer []byte
	for {
		time.Sleep(100*time.Millisecond)
		n, address, err := conn.ReadFromUDP(buffer)
		check_error(err)
		fmt.Println("Got message from ", address, " with n = ", address, n)
		if n > 0 {
			fmt.Println("From address: ", address, " got message: ", string(buffer[0:n]), n)
		}
	}
}


func broadcast_UDP(port, msg string, length int) {
	IP, err := net.ResolveUDPAddr("udp", localIP)
	conn, err := net.DialUDP("udp", IP, serverAddr)
	check_error(err)
	defer conn.Close()

	for {
		time.Sleep(1000*time.Microsecond)
		_, err := conn.WriteToUDP([]byte("Hello from the other siiiiiideeee \n"), serverAddr)
		check_error(err)
	}
}
