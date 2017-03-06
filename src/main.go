package main

import (
	"constants"
	"driver"
	"flag"
	"fmt"
	"fsm"
	"network/bcast"
	"network/localip"
	"network/peers"
	"structs"
	"test"
	"time"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan structs.Button)
	newOrderchan := make(chan bool, 10)
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
	recieve_QUEUE := make(chan structs.UDP_queue)

	lightSync := make(chan structs.Button, 10)
	//lightClear := make(chan structs.Button, 10)
	fmt.Println("Initialising FSM")

	fsm.Elev_init_own()
	//fmt.Println("Initialising Network")
	//network.Init(conn, msg_chan)

	go test.Get_Button_Press(newButtonchan, lightSync)
	go test.Get_new_floor(newFloorchan)
	go test.Update_queues(newButtonchan, newOrderchan, transmit_QUEUE)
	go fsm.Run(newOrderchan, newFloorchan, doorTimeoutchan, doorResetchan, lightSync)
	go test.Update_lights(lightSync)

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	//TODO Move to netowkr init
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

	go peers.Transmitter(20142, id, peerTxEnable)
	go peers.Receiver(20142, peerUpdateCh)

	transmit_status := make(chan fsm.Elevator)
	received_status := make(chan fsm.Elevator)

	light_recieve := make(chan structs.Button)
	//light_transmit := make(chan structs.Button)

	go bcast.Transmitter(30142, transmit_status)
	go bcast.Receiver(30142, received_status)
	go bcast.Receiver(20143, recieve_QUEUE)
	go bcast.Receiver(29444, light_recieve)
	//ALIVE BROADCAST
	go func() {
		for {
			transmit_status <- fsm.LocalElev
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			a := transmit_QUEUE
			fmt.Println("SENDING A QUEUE, ", a)
			bcast.Transmitter(20143, a)
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
			case l := <-light_recieve:
				fmt.Println("Recieved a light")
				if l.B_type != constants.B_CMD {
					driver.Elev_set_button_lamp(l.B_type, l.Floor, l.Value)
				}
			case a := <-received_status:
				//fmt.Println(a.ID)
				//UPDATES ELEVATOR MAP
				//TODO Finn ut hvorfor dette ikke funker:
				if a.ID != fsm.LocalElev.ID {
					fsm.Update_elevator_map(a)
				}

				/*if a.ID != fsm.LocalElev.ID {
					fmt.Println("INFO FROM ", a.ID)
					for i := 0; i < len(fsm.Elevator_list); i++ {
						if a.ID == fsm.Elevator_list[i].ID {
							fsm.Elevator_list[i] = a
							fmt.Println("Updated: ", fsm.Elevator_list[i].ID)
						}
					}
				}

				/*case v := <-transmit_QUEUE:
				fmt.Println("SENDING A QUEUE, ", transmit_QUEUE)
				bcast.Transmitter(20143, v)*/
			case q := <-recieve_QUEUE:
				if q.IP == fsm.LocalElev.ID {
					fmt.Println("RECIEVED A QUEUE", q)
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
