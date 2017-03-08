package order_manager

import (
	"constants"
	"driver"
	"structs"
)

type Queue struct {
	Queue_matrix [constants.N_FLOORS][constants.N_BUTTONS]bool
}

func Make_empty_queue() Queue {
	var ret Queue
	for f := 0; f < constants.N_FLOORS; f++ {
		for b := 0; b < constants.N_BUTTONS; b++ {
			ret.Queue_matrix[f][b] = false //TODO The init must check the backup, and update correctly
		}
	}
	return ret
}

/*func (q *Queue) Add_order_to_queue(c chan structs.Button, o chan bool) {
	for {
		select {
		case button_pressed := <-c:
			fmt.Println("Adding order to queue")
			q.Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
			q.set_lights()
			o <- true
		}
	}
}*/

func (q *Queue) Add_order_to_queue(button_pressed structs.Button, o chan bool) {
	q.Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
	//q.set_lights()
	q.set_local_lights()
	o <- true
}

func (q *Queue) Clear_lights_at_floor(floor, dir int) (structs.Button, structs.Button) {
	var button_up structs.Button
	button_up.B_type = constants.B_UP
	var button_down structs.Button
	button_down.B_type = constants.B_DOWN
	button_up.Value = true
	button_down.Value = true

	switch dir {
	case constants.UP:
		button_up.Value = false
		button_up.Floor = floor
		if !(q.order_above(floor)) {
			button_down.Value = false
			button_down.Floor = floor
		}
	case constants.DOWN:
		button_down.Value = false
		button_down.Floor = floor
		if !(q.order_below(floor)) {
			button_up.Value = false
			button_up.Floor = floor
		}
	case constants.STOP:
		button_up.Value = false
		button_up.Floor = floor
		button_down.Value = false
		button_down.Floor = floor
	}
	return button_up, button_down
}

func (q *Queue) Clear_orders_at_floor(floor, dir int) {
	q.Queue_matrix[floor][constants.B_CMD] = false
	switch dir {
	case constants.UP:
		q.Queue_matrix[floor][constants.B_UP] = false
		if !(q.order_above(floor)) {
			q.Queue_matrix[floor][constants.B_DOWN] = false
		}
	case constants.DOWN:
		q.Queue_matrix[floor][constants.B_DOWN] = false
		if !(q.order_below(floor)) {
			q.Queue_matrix[floor][constants.B_UP] = false
		}
	case constants.STOP:
		q.Queue_matrix[floor][constants.B_UP] = false
		q.Queue_matrix[floor][constants.B_DOWN] = false
	}
	//q.set_lights()
	q.set_local_lights()
}

func (q *Queue) set_local_lights() {
	for f := 0; f < constants.N_FLOORS; f++ {
		if q.Queue_matrix[f][constants.B_CMD] {
			driver.Elev_set_button_lamp(constants.B_CMD, f, true)
		} else {
			driver.Elev_set_button_lamp(constants.B_CMD, f, false)
		}
	}
}

func (q *Queue) set_lights() {
	for f := 0; f < constants.N_FLOORS; f++ {
		for b := 0; b < constants.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				driver.Elev_set_button_lamp(b, f, true)
			} else {
				driver.Elev_set_button_lamp(b, f, false)
			}
		}
	}
}
func (q *Queue) Clear_order(floor int, button int) {
	q.Queue_matrix[floor][button] = false
}
func (q *Queue) Should_stop(floor int, dir int) bool {
	switch dir {
	case constants.UP:
		return q.Queue_matrix[floor][constants.B_UP] ||
			q.Queue_matrix[floor][constants.B_CMD] ||
			floor == constants.N_FLOORS-1 ||
			!q.order_above(floor)
	case constants.DOWN:
		if q.Queue_matrix[floor][constants.B_DOWN] || q.Queue_matrix[floor][constants.B_CMD] || floor == 0 || !q.order_below(floor) {
			return true
		}
	default:
		return true
	}
	return false
}

func (q *Queue) order_above(floor int) bool {
	for f := floor + 1; f < constants.N_FLOORS; f++ {
		for b := 0; b < constants.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				return true
			}
		}
	}
	return false
}

func (q *Queue) Is_order(f int, b int) bool {
	return q.Queue_matrix[f][b]
}

func (q *Queue) order_below(floor int) bool {
	for f := floor - 1; f >= 0; f-- {
		for b := 0; b < constants.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				return true
			}
		}
	}
	return false
}

func (q *Queue) Is_empty() bool {
	for f := 0; f < constants.N_FLOORS; f++ {
		for b := 0; b < constants.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				return false
			}
		}
	}
	return true
}

func (q *Queue) Choose_dir(floor, dir int) int {
	if q.Is_empty() {
		return constants.STOP
	}
	switch dir {
	case constants.UP:
		if q.order_above(floor) {
			return constants.UP
		} else if q.order_below(floor) {
			return constants.DOWN
		} else {
			return constants.STOP
		}
	case constants.DOWN:
		if q.order_below(floor) {
			return constants.DOWN
		} else if q.order_above(floor) {
			return constants.UP
		} else {
			return constants.STOP
		}

	case constants.STOP:
		if q.order_above(floor) {
			return constants.UP
		} else if q.order_below(floor) {
			return constants.DOWN
		} else {
			return constants.STOP
		}
	default:
		return constants.STOP
	}
}
