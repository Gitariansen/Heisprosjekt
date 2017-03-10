package queue_methods

import (
	"driver"
)

type Queue struct {
	Queue_matrix [driver.N_FLOORS][driver.N_BUTTONS]bool
}

func Make_empty_queue() Queue {
	var ret Queue
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			ret.Queue_matrix[f][b] = false //TODO The init must check the backup, and update correctly
		}
	}
	return ret
}

func (q *Queue) Add_order_to_queue(button_pressed driver.Button, o chan bool) {
	q.Queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
	q.set_local_lights()
	o <- true
}

func (q *Queue) Clear_lights_at_floor(floor, dir int) (driver.Button, driver.Button) {
	var button_up driver.Button
	button_up.B_type = driver.B_UP
	var button_down driver.Button
	button_down.B_type = driver.B_DOWN
	button_up.Value = true
	button_down.Value = true

	switch dir {
	case driver.DIR_UP:
		button_up.Value = false
		button_up.Floor = floor
		if !(q.order_above(floor)) {
			button_down.Value = false
			button_down.Floor = floor
		}
	case driver.DIR_DOWN:
		button_down.Value = false
		button_down.Floor = floor
		if !(q.order_below(floor)) {
			button_up.Value = false
			button_up.Floor = floor
		}
	case driver.DIR_STOP:
		button_up.Value = false
		button_up.Floor = floor
		button_down.Value = false
		button_down.Floor = floor
	}
	return button_up, button_down
}

func (q *Queue) Clear_orders_at_floor(floor, dir int) {
	q.Queue_matrix[floor][driver.B_CMD] = false
	switch dir {
	case driver.DIR_UP:
		q.Queue_matrix[floor][driver.B_UP] = false
		if !(q.order_above(floor)) {
			q.Queue_matrix[floor][driver.B_DOWN] = false
		}
	case driver.DIR_DOWN:
		q.Queue_matrix[floor][driver.B_DOWN] = false
		if !(q.order_below(floor)) {
			q.Queue_matrix[floor][driver.B_UP] = false
		}
	case driver.DIR_STOP:
		q.Queue_matrix[floor][driver.B_UP] = false
		q.Queue_matrix[floor][driver.B_DOWN] = false
	}
	q.set_local_lights()
}

func (q *Queue) set_local_lights() {
	for f := 0; f < driver.N_FLOORS; f++ {
		if q.Queue_matrix[f][driver.B_CMD] {
			driver.Elev_set_button_lamp(driver.B_CMD, f, true)
		} else {
			driver.Elev_set_button_lamp(driver.B_CMD, f, false)
		}
	}
}

func (q *Queue) set_lights() {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
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
	case driver.DIR_UP:
		return q.Queue_matrix[floor][driver.B_UP] ||
			q.Queue_matrix[floor][driver.B_CMD] ||
			floor == driver.N_FLOORS-1 ||
			!q.order_above(floor)
	case driver.DIR_DOWN:
		if q.Queue_matrix[floor][driver.B_DOWN] || q.Queue_matrix[floor][driver.B_CMD] || floor == 0 || !q.order_below(floor) {
			return true
		}
	default:
		return true
	}
	return false
}

func (q *Queue) order_above(floor int) bool {
	for f := floor + 1; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
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
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				return true
			}
		}
	}
	return false
}

func (q *Queue) Is_empty() bool {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.Queue_matrix[f][b] {
				return false
			}
		}
	}
	return true
}

func (q *Queue) Choose_dir(floor, dir int) int {
	if q.Is_empty() {
		return driver.DIR_STOP
	}
	switch dir {
	case driver.DIR_UP:
		if q.order_above(floor) {
			return driver.DIR_UP
		} else if q.order_below(floor) {
			return driver.DIR_DOWN
		} else {
			return driver.DIR_STOP
		}
	case driver.DIR_DOWN:
		if q.order_below(floor) {
			return driver.DIR_DOWN
		} else if q.order_above(floor) {
			return driver.DIR_UP
		} else {
			return driver.DIR_STOP
		}

	case driver.DIR_STOP:
		if q.order_above(floor) {
			return driver.DIR_UP
		} else if q.order_below(floor) {
			return driver.DIR_DOWN
		} else {
			return driver.DIR_STOP
		}
	default:
		return driver.DIR_STOP
	}
}
