package main

import (
	"constants"
	"driver"
	"event_manager"
	"flag"
	"fmt"
	"fsm"
	"network/bcast"
	"network/localip"
	"network/peers"
	"os"
	"os/signal"
	"structs"
	"test"
	"time"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan structs.Button, 10)
	newOrderchan := make(chan bool, 10)
	newFloorchan := make(chan int)
	doorTimeoutchan := make(chan bool)
	doorResetchan := make(chan bool)
	floorResetchan := make(chan bool)
	// We make a channel for receiving updates on the id's of the peers that are
	// Alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	transmit_QUEUE := make(chan structs.UDP_queue, 10)
	recieve_QUEUE := make(chan structs.UDP_queue, 10)

	light_transmit := make(chan structs.Button, 10)
	//lightClear := make(chan structs.Button, 10)
	fmt.Println("Initialising FSM")

	fsm.Elev_init_own()
	//fmt.Println("Initialising Network")
	//network.Init(conn, msg_chan)

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
	peer_port := 20142
	go peers.Transmitter(peer_port, id, peerTxEnable)
	go peers.Receiver(peer_port, peerUpdateCh)

	transmit_status := make(chan fsm.Elevator, 10)
	received_status := make(chan fsm.Elevator)

	light_recieve := make(chan structs.Button, 10)
	//light_transmit := make(chan structs.Button)
	alive_port := 30142
	queue_port := 20413
	light_port := 29444
	// TODO Setup network init
	go bcast.Transmitter(alive_port, transmit_status)
	go bcast.Receiver(alive_port, received_status)

	go bcast.Receiver(queue_port, recieve_QUEUE)
	go bcast.Transmitter(queue_port, transmit_QUEUE)

	go bcast.Receiver(light_port, light_recieve)
	go bcast.Transmitter(light_port, light_transmit)

	go test.Get_Button_Press(newButtonchan)
	go test.Get_new_floor(newFloorchan)
	go event_manager.Order_manager(newButtonchan, newOrderchan, transmit_QUEUE, light_transmit)
	go fsm.Run(newOrderchan, newFloorchan, doorTimeoutchan, doorResetchan, light_transmit, floorResetchan)
	//go safeKill()
	go fsm.CheckMotorResponse(floorResetchan)

	//ALIVE BROADCAST
	go func() {
		time.Sleep(5 * time.Second) //wait for other incoming messages
		fmt.Println("Started Alive-spam")
		for {
			transmit_status <- fsm.LocalElev
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			fsm.Melevator[fsm.LocalElev.ID] = fsm.LocalElev
			time.Sleep(100 * time.Millisecond)
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

				peers.UpdateOnlineElevators(p, newButtonchan, transmit_QUEUE)

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
			case q := <-recieve_QUEUE:
				if q.IP == fsm.LocalElev.ID {
					fmt.Println("RECIEVED A QUEUE", q)
					//newButtonchan <- q.Button
					fsm.LocalElev.Queue.Add_order_to_queue(q.Button, newOrderchan)
					fmt.Println("sent light 2")
					light_transmit <- q.Button
				}
			}
		}
	}()

	for {
		if driver.Elev_get_stop_signal() == 1 {
			//driver.Elev_set_motor_direction(0)
			fmt.Println(fsm.LocalElev.Queue)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func safeKill() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(constants.STOP)
	fsm.LocalElev.Active = false
	os.Exit(0)
}
