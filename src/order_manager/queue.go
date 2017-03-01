package order_manager

import (
	"driver"
)

var Queue_matrix = [driver.N_FLOORS][driver.N_BUTTONS]bool{}

func Get_queue() {}

func Add_order_to_queue(c chan driver.Button) {
	var button_pressed driver.Button
	for {
		select {
		case button_pressed = <-c:
			Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
		}
	}
}


func Should_stop(floor int, button int)bool{
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if(Queue_matrix[floor][button] == driver.Elev_get_floor_signal()){
				return true
			}
		}
		return false
	}

	/*func Choose_dir(floor, dir int) int{
		 if queuematrix is empty
		 		return 0
			if order above -> return up

			if order below -> return down

	}*/

	//func order_above(){}
	//func order_below(){}


//func Clear_orders_at_floor() {}

func Queue_init() {

	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			Queue_matrix[f][b] = false //TODO The init must check the backup, and update correctly
		}
	}
}
