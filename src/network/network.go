package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"fsm"
)

/*

ANDERS SE HER!!!
https://github.com/danielbmx/heisprosjekt/tree/master/src/networkmodule

*/


type UDPMessage struct {
	Msg   string
	Queue []bool
}

type TCPMessage struct {
	Msg  string
	Flag bool
}

type Connection struct {
	ConnUDP *net.UDPConn
	ConnTCP *net.TCPConn
}

type Peer struct {
	Address  string
	LastSeen time.Time
	Active   bool
}

//var serverAddr *net.UDPAddr
//var localIP string
var localIP = "129.241.187.255" //TODO Endre dette navnet til no som har med broadcast å gjøre
var port = ":30042"
var peerChan = make(chan string)
var peers []Peer

func check_error(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

// All credz Anders.
// Brukes til å sammenligne om det er en melingbuffer[0:n],&save fra deg selv.
// I såfall ignorer denne. 129.241.187.141
func Get_local_IP() string {
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
	check_error(err)
	defer conn.Close()

	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP
}

func Peers(peerChan chan string) {
	var peer Peer
	address := <-peerChan
	for i := 0; i < len(peers); i++ {
		if address == peers[i].Address {
			peers[i].Active = true
			peers[i].LastSeen = time.Now()
		} else {
			peer.Active = true
			peer.LastSeen = time.Now()
			peer.Address = address
			peers = append(peers, peer)
		}
	}
}

func Check_if_connected() {

}

func Init(c chan Connection, msg_chan chan fsm.Elevator){
	var store_conn Connection


	fmt.Println("Your local IP: ", Get_local_IP())

	// setting up UDP server for broadcasting
	serverAddr_UDP, err := net.ResolveUDPAddr("udp", localIP+port)
	check_error(err)
	ConnUDP, err := net.DialUDP("udp", nil, serverAddr_UDP)
	check_error(err)
	store_conn.ConnUDP = ConnUDP
	store_conn.ConnTCP = nil

	var test_msg UDPMessage
	test_msg.Msg = "I'm aliiiiiiveeee"

  go Recive_msg_UDP(msg_chan)
	go Broadcast_UDP(c, msg_chan)

	c <- store_conn
}

func Recive_msg_UDP(msg_chan chan fsm.Elevator) {
	var message fsm.Elevator

	serverAddr, err := net.ResolveUDPAddr("udp", localIP+port)
	check_error(err)
	fmt.Println(serverAddr)

	conn, err := net.ListenUDP("udp", serverAddr)
	check_error(err)
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, address, err := conn.ReadFromUDP(buffer)
		check_error(err)
		fmt.Println("Got message from ", address)
		err = json.Unmarshal(buffer[0:n], &message)
		check_error(err)
		fmt.Println(message.Queue)
		check_error(err)
		msg_chan <- message
	}
}

func Broadcast_UDP(c chan Connection, msg_chan chan fsm.Elevator) {
	conn_store := <-c
  var msg fsm.Elevator

	for {
		msg = fsm.TheElev
		json_msg, err := json.Marshal(msg)
		check_error(err)
		fmt.Println("Sending message...")
		fmt.Println(msg.Queue)
		_, err = conn_store.ConnUDP.Write([]byte(json_msg))
		check_error(err)
		time.Sleep(1000 * time.Millisecond)
		<-msg_chan
	}
}
