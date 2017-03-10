package config

import (
	"driver"
	"network/localip"
	"orderManager/queue"
)

var LocalElev Elevator
var ElevatorMap = make(map[string]Elevator)

type Elevator struct {
	Floor  int
	State  elevatorState
	Queue  queue.Queue
	ID     string
	Dir    int
	Tic    int
	Active bool
}

type QueueMessage struct {
	IP     string
	Button driver.Button
}

type elevatorState int

const (
	IDLE elevatorState = iota
	MOVING
	DOOR_OPEN
)

func ElevInit() {
	driver.IoInit()
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.ElevSetButtonLamp(b, f, false)
		}
	}
	if driver.ElevGetFloorSensorSignal() == -1 {
		driver.ElevSetMotorDirection(driver.DIR_DOWN)
	}
	for driver.ElevGetFloorSensorSignal() == -1 {
		//wait
	}
	driver.ElevSetMotorDirection(driver.DIR_STOP)
	driver.ElevSetFloorIndicator(LocalElev.Floor)
	driver.ElevSetDoorOpenLamp(false)
	LocalElev.ID, _ = localip.LocalIP()
	LocalElev.State = IDLE
	LocalElev.Dir = driver.DIR_STOP
	LocalElev.Floor = driver.ElevGetFloorSensorSignal()
	LocalElev.Active = true
	LocalElev.Queue = queue.MakeEmptyQueue()
	AddElevatorToMap(LocalElev)
}

func AddElevatorToMap(e Elevator) {
	if _, elevatorInMap := ElevatorMap[e.ID]; elevatorInMap {
		//Do nothing
	} else {
		ElevatorMap[e.ID] = e
	}
}

func UpdateElevatorMap(e Elevator) {
	ElevatorMap[e.ID] = e
}
