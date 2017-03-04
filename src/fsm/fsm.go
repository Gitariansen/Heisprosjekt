package fsm

import (
	"conf"
	"driver"
	"fmt"
	"order_manager"
	"time"
)

/**/
//elevator states, TODO move to config
type elevator_state int

const (
	IDLE elevator_state = iota
	MOVING
	DOOR_OPEN
)

var TheElev Elevator

type Elevator struct {
	Floor  int
	State  int
	Queue  order_manager.Queue
	ID     string //get_local_IP
	Dir    int
	Tic    int
	Active bool
}

func Elev_init_own() {
	driver.Io_init()
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.Elev_set_button_lamp(b, f, false)
		}
	}
	if driver.Elev_get_floor_sensor_signal() == -1 {
		driver.Elev_set_motor_direction(conf.DOWN)
	}
	for driver.Elev_get_floor_sensor_signal() == -1 {
		//wait
	}
	driver.Elev_set_motor_direction(conf.STOP)
	TheElev.Floor = driver.Elev_get_floor_sensor_signal()
	driver.Elev_set_floor_indicator(TheElev.Floor)
	TheElev.Queue = order_manager.Make_empty_queue() //TODO change make-empty-queue to update-queue
	fmt.Println("Elevator is initialized in floor: ", TheElev.Floor+1)
}

var states_test elevator_state

//hardkokt
func open_door(Door_timeout, Door_reset chan bool) {
	const length = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-Door_reset:
			timer.Reset(length)
		case <-timer.C:
			timer.Stop()
			Door_timeout <- true
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func Run(newOrder chan bool, newFloor chan int, Door_timeout chan bool, Door_reset chan bool) {
	go open_door(Door_timeout, Door_reset)
	for {
		select {
		case <-newOrder:
			new_order_in_queue(Door_reset)
		case floor := <-newFloor:
			arriving_at_floor(floor, Door_reset)
		case <-Door_timeout:
			door_timeout()
		}
	}
}

func new_order_in_queue(Door_reset chan bool) {
	switch states_test {
	case IDLE:
		TheElev.Dir = TheElev.Queue.Choose_dir(TheElev.Floor, TheElev.Dir)
		if TheElev.Dir == conf.STOP {
			states_test = DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			TheElev.Queue.Clear_orders_at_floor(TheElev.Floor, TheElev.Dir)
		} else {
			driver.Elev_set_motor_direction(TheElev.Dir)
			states_test = MOVING
		}
	case MOVING:
		//do nothing
	case DOOR_OPEN:
		TheElev.Dir = TheElev.Queue.Choose_dir(TheElev.Floor, TheElev.Dir)
		if TheElev.Dir == conf.STOP {
			Door_reset <- true
			TheElev.Queue.Clear_orders_at_floor(TheElev.Floor, TheElev.Dir)
		}
	}
}

func arriving_at_floor(f int, Door_reset chan bool) {
	TheElev.Floor = f
	driver.Elev_set_floor_indicator(f)
	switch states_test {
	case IDLE:
		//Do nothing
	case MOVING:
		if TheElev.Queue.Should_stop(f, TheElev.Dir) {
			driver.Elev_set_motor_direction(conf.STOP)
			TheElev.Queue.Clear_orders_at_floor(f, TheElev.Dir)
			states_test = DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
		}
	case DOOR_OPEN:
		//Do nothing
	default:
	}
}
func door_timeout() {
	switch states_test {
	case IDLE:
		//Do nothing
	case MOVING:
		//Do nothing
	case DOOR_OPEN:
		driver.Elev_set_door_open_lamp(false)
		TheElev.Dir = TheElev.Queue.Choose_dir(TheElev.Floor, TheElev.Dir)
		if TheElev.Dir == conf.STOP {
			states_test = IDLE
		} else {
			driver.Elev_set_motor_direction(TheElev.Dir)
			states_test = MOVING
		}
	default:
	}
}
