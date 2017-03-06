package peers

import (
	"fmt"
	"fsm"
	"net"
	"sort"
	"time"

	"network/conn"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

func UpdateOnlineElevators(p PeerUpdate) {
	var elev fsm.Elevator
	elev.ID = p.New

	//ADDS NEW ELEVATORS TO MAP AND ACTIVATES/DEACTIVATES ELEVATORS

	if len(p.New) != 0 { //IF THERE IS NEW ELEVATOR
		fsm.Add_elevator_to_map(elev)
		elev.Active = true
		fsm.Update_elevator_map(elev)

	}
	if len(p.Lost) != 0 { //IF THERE IS AN ELEVATOR LOST
		for i := 0; i < len(p.Lost); i++ {
			elev = fsm.Melevator[p.Lost[i]]
			elev.Active = false
			fsm.Melevator[p.Lost[i]] = elev
		}
	}
	j := 0
	for _, v := range fsm.Melevator {
		j++
		fmt.Println("j: ", j)
		fmt.Println(v.ID)
		if v.Active == true {
			fmt.Println("Active elevators ID's", j, v.ID)
		}

		/*for _, v := range fsm.Melevator {
			fmt.Println("All elevators: ", v.ID)
		}*/
	}

	/*if len(p.New) != 0 { //IF THERE IS NEW ELEVATOR
		inArr, _ := (network.InArray(elev, fsm.Elevator_list))
		fmt.Println("b: ", inArr)

		if !inArr {
			fsm.Elevator_list = append(fsm.Elevator_list, elev)
			fsm.IP_list = append(fsm.IP_list, elev.ID)
			fmt.Println("appendend ")
		}
		for i := 0; i < len(fsm.Elevator_list); i++ {
			if fsm.Elevator_list[i].ID == p.New {
				fsm.Elevator_list[i].Active = true
				fmt.Println("Elevator Active ", fsm.Elevator_list[i].ID, fsm.Elevator_list[i].Active)
			}
		}

		fmt.Println("fsm.Elevator_list is now: ", fsm.IP_list)
	}
	if len(p.Lost) != 0 { //IF THERE IS AN ELEVATOR LOST
		for i := 0; i < len(p.Lost); i++ {
			fmt.Println(i)
			for j := 0; j < len(fsm.Elevator_list); j++ {
				if fsm.Elevator_list[j].ID == p.Lost[i] {
					fsm.Elevator_list[j].Active = false
					fmt.Println("Elevator inactive")
				}
			}
		}
	}*/
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
