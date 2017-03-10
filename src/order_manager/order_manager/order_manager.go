package order_manager

import (
	"config"
	"driver"
	"fmt"
	"order_manager/cost"
	"time"
)

func SyncHallLights(transmitLight chan driver.Button) {
	var ret driver.Button
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS-1; b++ {
			ret.Value = config.LocalElev.Queue.Is_order(f, b)
			ret.Floor = f
			ret.B_type = b
			transmitLight <- ret
		}
	}
}
func Order_manager(b chan driver.Button, o chan bool, q chan config.UDP_queue, l chan driver.Button) {
	for {
	loop:
		select {
		case button_pressed := <-b:
			if button_pressed.B_type == driver.B_CMD {
				config.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
			} else if config.LocalElev.Active {
				for _, v := range config.ElevatorMap {
					if v.Active && v.Queue.Is_order(button_pressed.Floor, button_pressed.B_type) {
						fmt.Println("Order is already in someone's queue")
						break loop
					}
				}
				fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
				return_ID := config.LocalElev.ID
				fmt.Println("a")
				shortestTime := cost.TimeToIdle(config.LocalElev)
				for IP, v := range config.ElevatorMap {
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
				if return_ID == config.LocalElev.ID {
					fmt.Println("Local elevator was chosen")
					config.LocalElev.Queue.Add_order_to_queue(button_pressed, o)
					fmt.Println("Light sent")
					l <- button_pressed
					time.Sleep(10 * time.Millisecond)
				} else {
					fmt.Println("Elevator with ID chosen: ", return_ID)
					var temp config.UDP_queue
					temp.IP = return_ID
					temp.Button = button_pressed
					q <- temp
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
