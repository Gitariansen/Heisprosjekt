package queue

import (
	"driver"
)

type Queue struct {
	QueueMatrix [driver.N_FLOORS][driver.N_BUTTONS]bool
}

func MakeEmptyQueue() Queue {
	var ret Queue
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			ret.QueueMatrix[f][b] = false
		}
	}
	return ret
}

func (q *Queue) AddOrderToQueue(buttonPressed driver.Button) {
	q.QueueMatrix[buttonPressed.Floor][buttonPressed.BtnType] = true
	q.setCabLights()
	//o <- true
}

func (q *Queue) ClearLightsAtFloor(floor, dir int) (driver.Button, driver.Button) {
	var buttonUP driver.Button
	buttonUP.BtnType = driver.BTN_UP
	var buttonDown driver.Button
	buttonDown.BtnType = driver.BTN_DOWN
	buttonUP.Value = true
	buttonDown.Value = true
	switch dir {
	case driver.DIR_UP:
		buttonUP.Value = false
		buttonUP.Floor = floor
		if !(q.orderAbove(floor)) {
			buttonDown.Value = false
			buttonDown.Floor = floor
		}
	case driver.DIR_DOWN:
		buttonDown.Value = false
		buttonDown.Floor = floor
		if !(q.orderBelow(floor)) {
			buttonUP.Value = false
			buttonUP.Floor = floor
		}
	case driver.DIR_STOP:
		buttonUP.Value = false
		buttonUP.Floor = floor
		buttonDown.Value = false
		buttonDown.Floor = floor
	}
	return buttonUP, buttonDown
}

func (q *Queue) ClearOrdersAtFloor(floor, dir int) {
	q.QueueMatrix[floor][driver.BTN_CMD] = false
	switch dir {
	case driver.DIR_UP:
		q.QueueMatrix[floor][driver.BTN_UP] = false
		if !(q.orderAbove(floor)) {
			q.QueueMatrix[floor][driver.BTN_DOWN] = false
		}
	case driver.DIR_DOWN:
		q.QueueMatrix[floor][driver.BTN_DOWN] = false
		if !(q.orderBelow(floor)) {
			q.QueueMatrix[floor][driver.BTN_UP] = false
		}
	case driver.DIR_STOP:
		q.QueueMatrix[floor][driver.BTN_UP] = false
		q.QueueMatrix[floor][driver.BTN_DOWN] = false
	}
	q.setCabLights()
}

func (q *Queue) ClearOrder(floor int, button int) {
	q.QueueMatrix[floor][button] = false
}

func (q *Queue) ShouldStop(floor int, dir int) bool {
	switch dir {
	case driver.DIR_UP:
		return q.QueueMatrix[floor][driver.BTN_UP] ||
			q.QueueMatrix[floor][driver.BTN_CMD] ||
			floor == driver.N_FLOORS-1 ||
			!q.orderAbove(floor)
	case driver.DIR_DOWN:
		return q.QueueMatrix[floor][driver.BTN_DOWN] ||
			q.QueueMatrix[floor][driver.BTN_CMD] ||
			floor == 0 ||
			!q.orderBelow(floor)
	default:
		return true
	}
}

func (q *Queue) IsOrder(f int, b int) bool {
	return q.QueueMatrix[f][b]
}

func (q *Queue) IsEmpty() bool {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.QueueMatrix[f][b] {
				return false
			}
		}
	}
	return true
}

func (q *Queue) ChooseDir(floor, dir int) int {
	if q.IsEmpty() {
		return driver.DIR_STOP
	}
	switch dir {
	case driver.DIR_UP:
		if q.orderAbove(floor) {
			return driver.DIR_UP
		} else if q.orderBelow(floor) {
			return driver.DIR_DOWN
		} else {
			return driver.DIR_STOP
		}
	case driver.DIR_DOWN:
		if q.orderBelow(floor) {
			return driver.DIR_DOWN
		} else if q.orderAbove(floor) {
			return driver.DIR_UP
		} else {
			return driver.DIR_STOP
		}
	case driver.DIR_STOP:
		if q.orderAbove(floor) {
			return driver.DIR_UP
		} else if q.orderBelow(floor) {
			return driver.DIR_DOWN
		} else {
			return driver.DIR_STOP
		}
	default:
		return driver.DIR_STOP
	}
}

func (q *Queue) setCabLights() {
	for f := 0; f < driver.N_FLOORS; f++ {
		if q.QueueMatrix[f][driver.BTN_CMD] {
			driver.ElevSetButtonLamp(driver.BTN_CMD, f, true)
		} else {
			driver.ElevSetButtonLamp(driver.BTN_CMD, f, false)
		}
	}
}

func (q *Queue) orderBelow(floor int) bool {
	for f := floor - 1; f >= 0; f-- {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.QueueMatrix[f][b] {
				return true
			}
		}
	}
	return false
}

func (q *Queue) orderAbove(floor int) bool {
	for f := floor + 1; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if q.QueueMatrix[f][b] {
				return true
			}
		}
	}
	return false
}
