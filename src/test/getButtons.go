//get Buttons
package test

import (
	"constants"
	"driver"
	"structs"
	"time"
)

/*Iterate over all floors, all types of buttons (through the button_channel_matrix) and add to channel*/
func Get_Button_Press(c chan structs.Button) {
	var button_pressed structs.Button
	for {
		for floor := 0; floor < constants.N_FLOORS; floor++ {
			for button := 0; button < constants.N_BUTTONS; button++ {
				if driver.Elev_get_button_signal(button, floor) == 1 {
					button_pressed.Floor = floor
					button_pressed.B_type = button
					button_pressed.Value = true
					c <- button_pressed
					time.Sleep(500 * time.Millisecond)
					//for driver.Elev_get_button_signal(button, floor) == 1 {
					//TODO Don't do this
					//}
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
