package main

import (
	"bcast"
	"driver"
	"flag"
	"fmt"
	"fsm"
	"localip"
	"peers"
	"structs"
	"test"
	"time"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan structs.Button)
	newOrderchan := make(chan bool)
	newFloorchan := make(chan int)
	doorTimeoutchan := make(chan bool)
	doorResetchan := make(chan bool)
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	transmit_QUEUE := make(chan structs.UDP_queue)
	fmt.Println("Initialising FSM")
	fsm.Elev_init_own()

	//fmt.Println("Initialising Network")
	//network.Init(conn, msg_chan)

	go test.Get_Button_Press(newButtonchan)
	go test.Get_new_floor(newFloorchan)
	go test.Update_queues(newButtonchan, newOrderchan, transmit_QUEUE)
	go fsm.Run(newOrderchan, newFloorchan, doorTimeoutchan, doorResetchan)
	fmt.Println("test")

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		fsm.LocalElev.ID = localIP
		id = fmt.Sprintf(localIP)
		fmt.Println(fmt.Println("ID is: ", id))
	}

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	transmit_status := make(chan fsm.Elevator)
	received_status := make(chan fsm.Elevator)

	recieve_QUEUE := make(chan structs.UDP_queue)

	go bcast.Transmitter(16569, transmit_status)
	go bcast.Receiver(16569, received_status)

	go bcast.Receiver(17021, recieve_QUEUE)

	//ALIVE BROADCAST
	go func() {
		for {
			transmit_status <- fsm.LocalElev
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	go func() {
		for {
			select {
			case p := <-peerUpdateCh:
				fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", p.Peers)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)

				peers.UpdateOnlineElevators(p)

			case /*a :=*/ <-received_status:
				//UPDATES ELEVATOR MAP
				//TODO Finn ut hvorfor dette ikke funker:
				/*if a.ID != fsm.LocalElev.ID {
					fmt.Println("INFO FROM ", a.ID)
					for i := 0; i < len(fsm.Elevator_list); i++ {
						if a.ID == fsm.Elevator_list[i].ID {
							fsm.Elevator_list[i] = a
							fmt.Println("Updated: ", fsm.Elevator_list[i].ID)
						}
					}
				}*/

			case v := <-transmit_QUEUE:
				fmt.Println("SENDING A QUEUE, ", v)
				bcast.Transmitter(17021, transmit_QUEUE)
			case q := <-recieve_QUEUE:
				fmt.Println("RECIEVED A QUEUE", q)
				if q.IP == fsm.LocalElev.ID {
					newButtonchan <- q.Button
				}

			}
		}
	}()

	for {
		if driver.Elev_get_stop_signal() == 1 {
			driver.Elev_set_motor_direction(0)
			fmt.Println(fsm.LocalElev.Queue)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
