package network

import (
	"config"
	"driver"
	"flag"
	"fmt"
	"network/bcast"
	"network/localip"
	"network/peers"
	"time"
)

var transmitStatus chan config.Elevator
var receiveStatus chan config.Elevator
var transmitPeer chan bool
var receivePeer chan peers.PeerUpdate

var receiveQueue chan config.UDP_queue

var receiveLight chan driver.Button

func Init(newButton chan driver.Button, newOrder chan bool, transmitQueue chan config.UDP_queue, transmitLight chan driver.Button) {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf(localIP)
		fmt.Println(fmt.Println("ID is: ", id))
	}

	peer_port := 20142
	alive_port := 30142
	queue_port := 20413
	light_port := 29444

	transmitStatus = make(chan config.Elevator, 10)
	receiveStatus = make(chan config.Elevator)
	receivePeer = make(chan peers.PeerUpdate)
	transmitPeer = make(chan bool)
	//transmitQueue = make(chan config.UDP_queue, 10)
	receiveQueue = make(chan config.UDP_queue, 10)
	receiveLight = make(chan driver.Button, 10)
	//transmitLight = make(chan driver.Button, 10)

	go bcast.Transmitter(alive_port, transmitStatus)
	go bcast.Receiver(alive_port, receiveStatus)
	go bcast.Transmitter(queue_port, transmitQueue)
	go bcast.Receiver(queue_port, receiveQueue)
	go bcast.Transmitter(light_port, transmitLight)
	go bcast.Receiver(light_port, receiveLight)
	go peers.Transmitter(peer_port, id, transmitPeer)
	go peers.Receiver(peer_port, receivePeer)

	go handleIncomingMessages(newButton, newOrder, transmitQueue, transmitLight)
	go periodicStatusUpdate()

}

func handleIncomingMessages(newButton chan driver.Button, newOrder chan bool, transmitQueue chan config.UDP_queue, transmitLight chan driver.Button) {
	for {
		select {
		case p := <-receivePeer:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			peers.UpdatePeers(p, newButton, transmitQueue, transmitLight)

		case l := <-receiveLight:
			if l.B_type != driver.B_CMD {
				driver.Elev_set_button_lamp(l.B_type, l.Floor, l.Value)
			}

		case a := <-receiveStatus:
			if a.ID != config.LocalElev.ID {
				config.Update_elevator_map(a)
			}
		case q := <-receiveQueue:
			if q.IP == config.LocalElev.ID {
				fmt.Println("RECIEVED A QUEUE", q)
				config.LocalElev.Queue.Add_order_to_queue(q.Button, newOrder)
				transmitLight <- q.Button
			}
		}
	}
}

func periodicStatusUpdate() {
	time.Sleep(1 * time.Second) //wait for other incoming messages
	fmt.Println("Started Alive-spam")
	for {
		transmitStatus <- config.LocalElev
		time.Sleep(1 * time.Second)
	}
}
