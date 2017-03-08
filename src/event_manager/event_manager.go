package event_manager

import (
	"constants"
	"cost"
	"fmt"
	"fsm"
	"structs"
	"time"
)

func Order_manager(b chan structs.Button, o chan bool, q chan structs.UDP_queue, l chan structs.Button) {

	for {
	loop:
		select {
		case button_pressed := <-b:
			if button_pressed.B_type == constants.B_CMD {
				fsm.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
			} else if fsm.LocalElev.Active {
				for _, v := range fsm.Melevator {
					if v.Active && v.Queue.Is_order(button_pressed.Floor, button_pressed.B_type) {
						fmt.Println("Order is already in someone's queue")
						break loop
					}
				}
				fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
				//fmt.Println("a")
				return_ID := fsm.LocalElev.ID
				fmt.Println("a")
				shortestTime := cost.TimeToIdle(fsm.LocalElev)
				for IP, v := range fsm.Melevator {
					//fmt.Println("Evaluating elevator: ", IP)
					if v.Active {
						v2 := v
						v2.Queue.Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
						new_time := cost.TimeToIdle(v2)
						if new_time < shortestTime {
							shortestTime = new_time
							return_ID = IP
						}
					}
					fmt.Println("return ID: ", return_ID)
				}
				if return_ID == fsm.LocalElev.ID {
					fmt.Println("Local elevator was chosen")
					fsm.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
					fmt.Println("Light sent")
					l <- button_pressed
					time.Sleep(10 * time.Millisecond)
				} else {
					fmt.Println("Elevator with ID chosen: ", return_ID)
					var temp structs.UDP_queue
					temp.IP = return_ID
					temp.Button = button_pressed
					q <- temp
				}

			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
