package main

import (
	"driver"
	"fmt"
	"fsm"
	"network"
	"test"
	"time"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan driver.Button)
	newOrderchan := make(chan bool)
	newFloorchan := make(chan int)
	doorTimeoutchan := make(chan bool)
	doorResetchan := make(chan bool)

	conn := make(chan network.Connection)
	msg_chan := make(chan fsm.Elevator)
	//var msg_recieved = make(chan bool)

	fmt.Println("Initialising FSM")
	fsm.Elev_init_own()
	fmt.Println("Initialising Network")
	network.Init(conn, msg_chan)

	go test.Get_Button_Press(newButtonchan)
	go fsm.TheElev.Queue.Add_order_to_queue(newButtonchan, newOrderchan)
	go test.Get_new_floor(newFloorchan)

	go fsm.Run(newOrderchan, newFloorchan, doorTimeoutchan, doorResetchan)

	//go fsm.TheElev.Queue.Set_lights()

	// Test connection functions

	/*time.Sleep(100 * time.Millisecond)
	go network.Recive_msg_UDP(msg_chan)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("stuck test 1")

	go network.Broadcast_UDP(conn, test_msg, msg_chan)*/

	/*
		// Test Peers
		var testPeers []Peer
		peer1 := Test_struct{Address: "test1"}
		peer2 := Test_struct{Address: "test2"}
		peer3 := Test_struct{Address: "test3"}

		testPeers := append(testPeers, peer1)
		testPeers := append(testPeers, peer2)
		testPeers := append(testPeers, peer3)
	*/

	for {
		if driver.Elev_get_stop_signal() == 1 {
			driver.Elev_set_motor_direction(0)
			fmt.Println(fsm.TheElev.Queue)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
