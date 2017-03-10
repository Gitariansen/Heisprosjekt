package order_manager

import (
	"config"
	"driver"
	"fmt"
	"order_manager/cost"
	"time"
)

func OrderManager(b chan driver.Button, o chan bool, q chan config.UDP_queue, light_transmit chan driver.Button) {
	for {
	loop:
		select {
		case button_pressed := <-b:
			if button_pressed.B_type == driver.B_CMD {
				config.LocalElev.Queue.Add_order_to_queue(button_pressed)
				o <- true
			} else {
				if config.LocalElev.Active {
					for _, v := range config.ElevatorMap {
						if v.Active && v.Queue.Is_order(button_pressed.Floor, button_pressed.B_type) {
							break loop
						}
					}
					fmt.Println("----------------------------------------")
					returnID := calculateFastestElevator(button_pressed)
					if returnID == config.LocalElev.ID {
						fmt.Println("Local elevator was chosen")
						config.LocalElev.Queue.Add_order_to_queue(button_pressed)
						o <- true
						light_transmit <- button_pressed
					} else {
						fmt.Println("Elevator with ID: ", returnID, "chosen")
						var temp config.UDP_queue
						temp.IP = returnID
						temp.Button = button_pressed
						q <- temp
					}
					fmt.Println("----------------------------------------")
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func calculateFastestElevator(button_pressed driver.Button) string {
	returnID := config.LocalElev.ID
	shortestTime := cost.TimeToIdle(config.LocalElev, button_pressed)
	for IP, elev := range config.ElevatorMap {
		if elev.Active {
			elevCopy := elev
			elevCopy.Queue.Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
			newTime := cost.TimeToIdle(elevCopy, button_pressed)
			if newTime < shortestTime {
				shortestTime = newTime
				returnID = IP
			}
		}
	}
	return returnID
}

func TransmitLightSignal(transmitLight chan driver.Button) {
	var bup, bdwn driver.Button
	bup, bdwn = config.LocalElev.Queue.Clear_lights_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
	transmitLight <- bup
	transmitLight <- bdwn
}

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
