package config

import (
	"driver"
	"fmt"
	"network/localip"
	"order_manager/queue_methods"
)

var LocalElev Elevator
var ElevatorMap map[string]Elevator

type Elevator struct {
	Floor  int
	State  elevator_state
	Queue  queue_methods.Queue
	ID     string
	Dir    int
	Tic    int
	Active bool
}

type elevator_state int

const (
	IDLE elevator_state = iota
	MOVING
	DOOR_OPEN
)

func Elev_init() {
	driver.Io_init()
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.Elev_set_button_lamp(b, f, false)
		}
	}
	if driver.Elev_get_floor_sensor_signal() == -1 {
		driver.Elev_set_motor_direction(driver.DIR_DOWN)
	}
	for driver.Elev_get_floor_sensor_signal() == -1 {
		//wait
	}
	driver.Elev_set_motor_direction(driver.DIR_STOP)
	driver.Elev_set_floor_indicator(LocalElev.Floor)
	driver.Elev_set_door_open_lamp(false)

	ElevatorMap = make(map[string]Elevator)
	LocalElev.ID, _ = localip.LocalIP()
	LocalElev.State = IDLE
	LocalElev.Dir = driver.DIR_STOP
	LocalElev.Floor = driver.Elev_get_floor_sensor_signal()
	LocalElev.Active = true
	LocalElev.Queue = queue_methods.Make_empty_queue()
	Add_elevator_to_map(LocalElev)
	fmt.Println("Elevator is initialized in floor: ", LocalElev.Floor+1)

}

func Add_elevator_to_map(e Elevator) {
	if _, ok := ElevatorMap[e.ID]; ok {
		fmt.Println("Elevator already in map")
	} else {
		ElevatorMap[e.ID] = e
		fmt.Println("Elevator added to map")
	}
}

func Update_elevator_map(e Elevator) {
	ElevatorMap[e.ID] = e
	fmt.Println("UPDATED: ")
}

//TODO FIX UNDER HERE

type Light struct {
	Floor  int
	B_type int
	Value  bool
}

type Fsm_channels struct { //TODO implement this
	Door_timeout chan bool
	Door_reset   chan bool
	New_order    chan bool
	New_floor    chan int
	Sync_lights  chan bool
}

type UDP_queue struct {
	IP     string
	Button driver.Button
}
