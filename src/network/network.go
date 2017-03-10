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

type QueueMessage struct {
	IP     string
	Button driver.Button
}

var transmitStatus = make(chan config.Elevator, 10)
var receiveStatus = make(chan config.Elevator)
var receivePeer = make(chan peers.PeerUpdate)
var transmitPeer = make(chan bool)
var receiveQueue = make(chan config.QueueMessage, 12)
var receiveLight = make(chan driver.Button, 10)

func Init(transmitQueue chan config.QueueMessage, transmitLight chan driver.Button) {
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
	}
	initiateCommunicationGoRoutines(id, transmitQueue, transmitLight)
	go periodicStatusUpdate()

}
func initiateCommunicationGoRoutines(id string, transmitQueue chan config.QueueMessage, transmitLight chan driver.Button) {
	peer_port := 20142
	alive_port := 30142
	queue_port := 20413
	light_port := 29444
	go bcast.Transmitter(alive_port, transmitStatus)
	go bcast.Receiver(alive_port, receiveStatus)
	go bcast.Transmitter(queue_port, transmitQueue)
	go bcast.Receiver(queue_port, receiveQueue)
	go bcast.Transmitter(light_port, transmitLight)
	go bcast.Receiver(light_port, receiveLight)
	go peers.Transmitter(peer_port, id, transmitPeer)
	go peers.Receiver(peer_port, receivePeer)
}

func HandleIncomingMessages(newButton chan driver.Button, newOrder chan bool, transmitQueue chan config.QueueMessage, transmitLight chan driver.Button) {
	for {
		select {
		case receivedPeer := <-receivePeer:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", receivedPeer.Peers)
			fmt.Printf("  New:      %q\n", receivedPeer.New)
			fmt.Printf("  Lost:     %q\n", receivedPeer.Lost)
			peers.UpdatePeers(receivedPeer, newButton, transmitQueue, transmitLight)
		case receivedLight := <-receiveLight:
			if receivedLight.BtnType != driver.BTN_CMD {
				driver.ElevSetButtonLamp(receivedLight.BtnType, receivedLight.Floor, receivedLight.Value)
			}
		case receivedStatus := <-receiveStatus:
			if receivedStatus.ID != config.LocalElev.ID {
				config.UpdateElevatorMap(receivedStatus)
			}
		case receivedQueue := <-receiveQueue:
			if receivedQueue.IP == config.LocalElev.ID {
				config.LocalElev.Queue.AddOrderToQueue(receivedQueue.Button)
				newOrder <- true
				transmitLight <- receivedQueue.Button
			}
		}
	}
}

func periodicStatusUpdate() {
	time.Sleep(1 * time.Second)
	for {
		transmitStatus <- config.LocalElev
		time.Sleep(1 * time.Second)
	}
}
