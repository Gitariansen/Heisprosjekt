//get Buttons
package test

import (
	"conf"
	"cost"
	"driver"
	"fmt"
	"fsm"
	"structs"
	"time"
)

/*Iterate over all floors, all types of buttons (through the button_channel_matrix) and add to channel*/
func Get_Button_Press(c chan structs.Button) /*<-chan bool*/ {
	//ret := make(chan bool, 5)
	var button_pressed structs.Button
	for {
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for button := 0; button < driver.N_BUTTONS; button++ {
				if driver.Elev_get_button_signal(button, floor) == 1 {
					button_pressed.Floor = floor
					button_pressed.B_type = button
					c <- button_pressed
					/*for driver.Elev_get_button_signal(button, floor) == 1 {
							//TODO Don't do this
					}*/
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func Get_new_floor(ch chan int) {
	prev_floor := driver.Elev_get_floor_sensor_signal()
	for {
		curr_floor := driver.Elev_get_floor_sensor_signal()
		if curr_floor != -1 && curr_floor != prev_floor {
			ch <- curr_floor
		}
		prev_floor = curr_floor
	}
}

func Update_queues(b chan structs.Button, o chan bool, q chan structs.UDP_queue) {
	for {
		select {
		case button_pressed := <-b:
			if button_pressed.B_type == conf.B_CMD {
				fsm.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
			} else {
				fmt.Println("initiating choose elevator")
				fmt.Println("LOCAL: ", cost.TimeToIdle(fsm.LocalElev))
				ret_ID := fsm.LocalElev.ID
				localElevTime := cost.TimeToIdle(fsm.LocalElev)
				shortestTime := localElevTime
				for i := 0; i < len(fsm.Elevator_list); i++ {
					if fsm.Elevator_list[i].Active {
						new_time := cost.TimeToIdle(fsm.Elevator_list[i])
						fmt.Println("New Time ", new_time)
						if new_time < shortestTime {
							shortestTime = new_time
							ret_ID = fsm.Elevator_list[i].ID
						}
					}
				}
				fmt.Println("ret: ", ret_ID)
				fmt.Println(fsm.Elevator_list[0].ID, fsm.LocalElev.ID)
				if ret_ID == fsm.LocalElev.ID {
					fmt.Println("Local elevator was chosen")
					fsm.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
				} else {
					fmt.Println("Elevator with ID: ", ret_ID) //FIX THIS TODO
					var temp structs.UDP_queue
					temp.IP = ret_ID
					temp.Button = button_pressed
					q <- temp
					//send TCP message to IP hadde vært det beste
				}
			}
		}
	}

}
