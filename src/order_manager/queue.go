package order_manager

import (
	"conf"
	"driver"
	"structs"
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
	q.set_lights()
	o <- true
}

func (q *Queue) Clear_orders_at_floor(floor, dir int) {
	q.Queue_matrix[floor][conf.B_CMD] = false
	switch dir {
	case conf.UP:
		q.Queue_matrix[floor][conf.B_UP] = false
		if !(q.order_above(floor)) {
			q.Queue_matrix[floor][conf.B_DOWN] = false
		}
	case conf.DOWN:
		q.Queue_matrix[floor][conf.B_DOWN] = false
		if !(q.order_below(floor)) {
			q.Queue_matrix[floor][conf.B_UP] = false
		}
	case conf.STOP:
		q.Queue_matrix[floor][conf.B_UP] = false
		q.Queue_matrix[floor][conf.B_DOWN] = false
	}
	q.set_lights()
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

func (q *Queue) Should_stop(floor int, dir int) bool {
	switch dir {
	case conf.UP:
		if q.Queue_matrix[floor][conf.B_UP] || q.Queue_matrix[floor][conf.B_CMD] || floor == driver.N_FLOORS-1 || !q.order_above(floor) {
			return true
		}
	case conf.DOWN:
		if q.Queue_matrix[floor][conf.B_DOWN] || q.Queue_matrix[floor][conf.B_CMD] || floor == 0 || !q.order_below(floor) {
			return true
		}
	case conf.STOP:
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.Queue_matrix[floor][b] {
				return true
			}
		}
	default:
		return false
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

func (q *Queue) is_empty() bool {
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
	if q.is_empty() {
		return conf.STOP
	}
	switch dir {
	case conf.UP:
		if q.order_above(floor) {
			return conf.UP
		} else if q.order_below(floor) {
			return conf.DOWN
		} else {
			return conf.STOP
		}
	case conf.DOWN:
		if q.order_below(floor) {
			return conf.DOWN
		} else if q.order_above(floor) {
			return conf.UP
		} else {
			return conf.STOP
		}

	case conf.STOP:
		if q.order_above(floor) {
			return conf.UP
		} else if q.order_below(floor) {
			return conf.DOWN
		} else {
			return conf.STOP
		}
	default:
		return conf.STOP
	}
}
