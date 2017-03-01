//get Buttons
package test

import (
	"driver"
	"fmt"
	"time"
)

/*Iterate over all floors, all types of buttons (through the button_channel_matrix) and add to channel*/
func Get_Button_Press(c chan driver.Button) /*<-chan bool*/ {
	//ret := make(chan bool, 5)
	var button_pressed driver.Button
	for {
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for button := 0; button < driver.N_BUTTONS; button++ {
				if driver.Elev_get_button_signal(button, floor) == 1 {
					fmt.Println("Button Pressed in floor: ", floor+1)
					button_pressed.Floor = floor
					button_pressed.B_type = button
					c <- button_pressed
					for driver.Elev_get_button_signal(button, floor) == 1 {

					}
					//order_manager.Queue_matrix[floor][button] = true
					//

				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	//return ret
}
