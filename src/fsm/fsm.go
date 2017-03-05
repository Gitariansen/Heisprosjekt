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

var LocalElev Elevator
var Elevator_list []Elevator
var IP_list []string

type Elevator struct {
	Floor  int
	State  elevator_state
	Queue  order_manager.Queue
	ID     string //get_local_IP
	Dir    int
	Tic    int
	Active bool
	Dur    int
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
	LocalElev.State = IDLE
	LocalElev.Dir = conf.STOP
	LocalElev.Floor = driver.Elev_get_floor_sensor_signal()
	driver.Elev_set_floor_indicator(LocalElev.Floor)
	LocalElev.Queue = order_manager.Make_empty_queue() //TODO change make-empty-queue to update-queue
	fmt.Println("Elevator is initialized in floor: ", LocalElev.Floor+1)
}

func update_list() {

}

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
	switch LocalElev.State {
	case IDLE:
		LocalElev.Dir = LocalElev.Queue.Choose_dir(LocalElev.Floor, LocalElev.Dir)
		if LocalElev.Dir == conf.STOP {
			LocalElev.State = DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			LocalElev.Queue.Clear_orders_at_floor(LocalElev.Floor, LocalElev.Dir)
		} else {
			driver.Elev_set_motor_direction(LocalElev.Dir)
			LocalElev.State = MOVING
		}
	case MOVING:
		//do nothing
	case DOOR_OPEN:
		if LocalElev.Queue.Should_stop(LocalElev.Floor, LocalElev.Dir) {
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			LocalElev.Queue.Clear_orders_at_floor(LocalElev.Floor, LocalElev.Dir)
		}
	}
}

func arriving_at_floor(f int, Door_reset chan bool) {
	LocalElev.Floor = f
	driver.Elev_set_floor_indicator(f)
	switch LocalElev.State {
	case IDLE:
		//Do nothing
	case MOVING:
		if LocalElev.Queue.Should_stop(f, LocalElev.Dir) {
			driver.Elev_set_motor_direction(conf.STOP)
			LocalElev.Queue.Clear_orders_at_floor(f, LocalElev.Dir)
			LocalElev.State = DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
		}
	case DOOR_OPEN:
		//Do nothing
	default:
	}
}
func door_timeout() {
	switch LocalElev.State {
	case IDLE:
		//Do nothing
	case MOVING:
		//Do nothing
	case DOOR_OPEN:
		driver.Elev_set_door_open_lamp(false)
		LocalElev.Dir = LocalElev.Queue.Choose_dir(LocalElev.Floor, LocalElev.Dir)
		if LocalElev.Dir == conf.STOP {
			LocalElev.State = IDLE
		} else {
			driver.Elev_set_motor_direction(LocalElev.Dir)
			LocalElev.State = MOVING
		}
	default:
	}
}
