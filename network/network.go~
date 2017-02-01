package network


import (
	"fmt"
	"net"
    "os"
)


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
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port:53})
	go check_error(err)
	defer conn.Close()

	IP = strings.Split(conn.LocalAddr().String(), ":")[0]
	fmt.Println("Your local IP address is: ", IP)
	return IP
}



func init() {
    // setting up UDP server
    localIP, err := net.ResolveUDPAddr("udp", net.JoinHostPort(get_local_IP(), port))
    check_error(err)
}
