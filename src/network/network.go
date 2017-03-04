package network

import (
	//"encoding/binary"
	"encoding/json"
	"fmt"
	"fsm"
	"net"
	"os"
	"strings"
	"time"
)

/*

ANDERS SE HER!!!
https://github.com/danielbmx/heisprosjekt/tree/master/src/networkmodule

*/

type Connection struct {
	ConnUDP *net.UDPConn
	ConnTCP *net.TCPConn
}

//var serverAddr *net.UDPAddr
var localIP string
var broadcastIP = "129.241.187.255"
var port = ":30083"
var peerChan = make(chan string)
var peers []fsm.Elevator

func check_error(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func inArray(a string, array []string) (bool, int) {
	for i := 0; i < len(array); i++ {
		if a == array[i] {
			return true, i
		}
	}
	return false, 0
}

func Init(c chan Connection, msg_chan chan fsm.Elevator) {
	var store_conn Connection
	localIP = Get_local_IP()

	fmt.Println("Your local IP: ", localIP)
	fsm.TheElev.ID = localIP

	find_other_elevators(peers, c, msg_chan)

	Peers(fsm.TheElev, peers)

	// setting up UDP server for broadcasting
	serverAddr_UDP, err := net.ResolveUDPAddr("udp", broadcastIP+port)
	check_error(err)

	ConnUDP, err := net.DialUDP("udp", nil, serverAddr_UDP)
	check_error(err)

	store_conn.ConnUDP = ConnUDP
	store_conn.ConnTCP = nil

	go Recive_msg_UDP(msg_chan)
	go Broadcast_UDP(c, msg_chan)

	c <- store_conn
}

func find_other_elevators(peers []fsm.Elevator, c chan Connection, msg_chan chan fsm.Elevator) {
	var message fsm.Elevator
	i := false
	//making a temporary connection
	serverAddr, err := net.ResolveUDPAddr("udp", broadcastIP+port)
	check_error(err)
	conn, err := net.ListenUDP("udp", serverAddr)
	check_error(err)

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	for !(i) {
		fmt.Println("iter")
		n, _, err := conn.ReadFromUDP(buffer)
		err = json.Unmarshal(buffer[0:n], &message)
		if err == nil {
			fmt.Println("err = nil, I recieved a message")
			Peers(message, peers)
		} else {
			i = true
			fmt.Println("There was an error, I am the only elevator")
		}
	}
	conn.Close()

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

func Peers(elev fsm.Elevator, peers []fsm.Elevator) {
	in_peers_flag := false
	if len(peers) == 0 {
		fmt.Println("Peers was empty, adding this elevator as peer with address: ", elev.ID)
		peers = append(peers, fsm.TheElev)
		fmt.Println("length of peers is now ", len(peers))
	} else {
		for i := 0; i < len(peers); i++ {
			if elev == peers[i] {
				in_peers_flag = true
				fmt.Println("Peer already exists")
				elev.Active = true
				elev.Tic = 0
			}
			if !in_peers_flag {
				fmt.Println("Adding new peer with address:", elev.ID)
				elev.Active = true
				elev.Tic = 0
				peers = append(peers, elev)
			}
		}
	}
	fmt.Println("Currently active peers: ")
	for i := 0; i < len(peers); i++ {
		fmt.Println("	", peers[i].ID)

	}
}

func Check_if_connected() {
	// Remember mutex
	for {
		fsm.TheElev.Tic++
		if fsm.TheElev.Tic >= 30 {
			fsm.TheElev.Active = false
		}
	}
}

func Recive_msg_UDP(msg_chan chan fsm.Elevator) {
	var message fsm.Elevator

	serverAddr, err := net.ResolveUDPAddr("udp", broadcastIP+port)
	check_error(err)

	conn, err := net.ListenUDP("udp", serverAddr)
	check_error(err)
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, address, err := conn.ReadFromUDP(buffer)
		check_error(err)
		//Peers(address.IP.String())
		fmt.Println("Got message from ", address)
		err = json.Unmarshal(buffer[0:n], &message)
		check_error(err)
		//fmt.Println(message.Queue)
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
		//fmt.Println(msg.Queue)
		_, err = conn_store.ConnUDP.Write([]byte(json_msg))
		check_error(err)
		time.Sleep(1000 * time.Millisecond)
		<-msg_chan
	}
}
