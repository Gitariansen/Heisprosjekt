package peers

import (
	"config"
	"driver"
	"fmt"
	"net"
	"network/conn"
	"order_manager/order_manager"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

func UpdatePeers(p PeerUpdate, newButton chan driver.Button, transmit_queue chan config.UDP_queue, transmitLight chan driver.Button) {
	newID := p.New
	lostIDs := p.Lost
	if len(p.New) != 0 { //IF THERE IS NEW ELEVATOR
		if newID == config.LocalElev.ID {
			fmt.Println("I am the new elevator")
			config.LocalElev.Active = true
		} else if _, elevatorInMap := config.ElevatorMap[newID]; elevatorInMap {
			newElevCopy := config.ElevatorMap[newID]
			for f := 0; f < driver.N_FLOORS; f++ {
				if newElevCopy.Queue.Is_order(f, driver.B_CMD) {
					for i := 0; i < 5; i++ {
						orderButton := driver.Button{Floor: f, B_type: driver.B_CMD, Value: true}
						order := config.UDP_queue{IP: newElevCopy.ID, Button: orderButton}
						transmit_queue <- order
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
			order_manager.SyncHallLights(transmitLight)
		} else {
			var elev config.Elevator
			elev.ID = newID
			config.Add_elevator_to_map(elev)
		}
	}
	if len(p.Lost) != 0 { //IF THERE IS AN ELEVATOR LOST
		for i := 0; i < len(lostIDs); i++ {
			elev := config.ElevatorMap[lostIDs[i]]
			elev.Active = false
			config.ElevatorMap[p.Lost[i]] = elev
		}
		if len(p.Peers) == 0 {
			fmt.Println("I am alone on the network")
			config.LocalElev.Active = false
		} else {
			for i := 0; i < len(p.Lost); i++ {
				redistributeQueue := config.ElevatorMap[p.Lost[i]].Queue
				for b := 0; b < driver.N_BUTTONS-1; b++ {
					for f := 0; f < driver.N_FLOORS; f++ {
						if redistributeQueue.Is_order(f, b) {
							newOrder := driver.Button{Floor: f, B_type: b, Value: true}
							newButton <- newOrder
							redistributeQueue.Clear_order(f, b)
						}
					}
				}
				elev := config.ElevatorMap[p.Lost[i]]
				elev.Queue = redistributeQueue
				config.ElevatorMap[p.Lost[i]] = elev
			}
		}

	}
}

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
