package peers

import (
	"config"
	"driver"
	"fmt"
	"net"
	"network/conn"
	"orderManager"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

func UpdatePeers(receivedPeer PeerUpdate, newButton chan driver.Button, transmitQueue chan config.QueueMessage, transmitLight chan driver.Button) {
	newID := receivedPeer.New
	lostIDs := receivedPeer.Lost
	if len(receivedPeer.New) != 0 {
		if newID == config.LocalElev.ID {
			config.LocalElev.Active = true
		} else if _, elevatorInMap := config.ElevatorMap[newID]; elevatorInMap {
			newElevCopy := config.ElevatorMap[newID]
			for f := 0; f < driver.N_FLOORS; f++ {
				if newElevCopy.Queue.IsOrder(f, driver.BTN_CMD) {
					for i := 0; i < 5; i++ {
						orderButton := driver.Button{Floor: f, BtnType: driver.BTN_CMD, Value: true}
						order := config.QueueMessage{IP: newElevCopy.ID, Button: orderButton}
						transmitQueue <- order
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
			orderManager.SyncHallLights(transmitLight)
		} else {
			var elev config.Elevator
			elev.ID = newID
			config.AddElevatorToMap(elev)
		}
	}
	if len(receivedPeer.Lost) != 0 {
		for i := 0; i < len(lostIDs); i++ {
			lostElev := config.ElevatorMap[lostIDs[i]]
			lostElev.Active = false
			config.ElevatorMap[receivedPeer.Lost[i]] = lostElev
		}
		if len(receivedPeer.Peers) == 0 {
			config.LocalElev.Active = false
		} else {
			for i := 0; i < len(receivedPeer.Lost); i++ {
				redistributeQueue := config.ElevatorMap[receivedPeer.Lost[i]].Queue
				for b := 0; b < driver.N_BUTTONS-1; b++ {
					for f := 0; f < driver.N_FLOORS; f++ {
						if redistributeQueue.IsOrder(f, b) {
							newOrder := driver.Button{Floor: f, BtnType: b, Value: true}
							newButton <- newOrder
							redistributeQueue.ClearOrder(f, b)
						}
					}
				}
				lostElev := config.ElevatorMap[receivedPeer.Lost[i]]
				lostElev.Queue = redistributeQueue
				config.ElevatorMap[receivedPeer.Lost[i]] = lostElev
			}
		}
	}
}

/*
Code belowd provided by klasbo
https://github.com/TTK4145/Network-go/tree/master/network
*/
const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}
