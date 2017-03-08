package peers

import (
	"constants"
	"fmt"
	"fsm"
	"net"
	"network/conn"
	"sort"
	"structs"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

func UpdateOnlineElevators(p PeerUpdate, button_chan chan structs.Button, transmit_queue chan structs.UDP_queue) {
	var elev fsm.Elevator
	elev.ID = p.New

	//ADDS NEW ELEVATORS TO MAP AND ACTIVATES/DEACTIVATES ELEVATORS

	if len(p.New) != 0 { //IF THERE IS NEW ELEVATOR

		if p.New == fsm.LocalElev.ID {
			fmt.Println("I am the new elevator")
		} else if _, ok := fsm.Melevator[p.New]; ok {
			temp := fsm.Melevator[p.New]
			fmt.Println(temp.Queue)
			for f := 0; f < constants.N_FLOORS; f++ {
				if temp.Queue.Is_order(f, constants.B_CMD) {
					return_btn := structs.Button{f, constants.B_CMD, true}
					order := structs.UDP_queue{temp.ID, return_btn}
					fmt.Println("Sending queue, ", order)
					transmit_queue <- order //TODO failsafe this
				}
			}
		}
		elev.Active = true
		fsm.Update_elevator_map(elev)
	}
	if len(p.Lost) != 0 { //IF THERE IS AN ELEVATOR LOST
		for i := 0; i < len(p.Lost); i++ {
			elev := fsm.Melevator[p.Lost[i]]
			elev.Active = false
			fsm.Melevator[p.Lost[i]] = elev
		}
		if len(p.Peers) == 4 {
			fsm.LocalElev.Active = false
		} else {
			for i := 0; i < len(p.Lost); i++ {
				temp_queue := fsm.Melevator[p.Lost[i]].Queue
				for b := 0; b < constants.N_BUTTONS-1; b++ {
					for f := 0; f < constants.N_FLOORS; f++ {
						if temp_queue.Is_order(f, b) {
							fmt.Println("Adding order to queue, floor: ", f, " button: ", b)
							return_btn := structs.Button{f, b, true}
							button_chan <- return_btn
							temp_queue.Clear_order(f, b)
						}
					}
				}
				v2 := fsm.Melevator[p.Lost[i]]
				v2.Queue = temp_queue
				fsm.Melevator[p.Lost[i]] = v2
				fmt.Println(v2.Queue)
			}
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
