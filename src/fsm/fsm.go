package fsm

import (
	"config"
	"driver"
	"fmt"
	"orderManager"
	"os"
	"os/signal"
	"time"
)

func Init(newOrder chan bool, newFloor chan int, doorTimeout chan bool, doorReset chan bool, responseTimerReset chan bool, transmitLight chan driver.Button) {
	go runFSM(newOrder, newFloor, doorTimeout, doorReset, transmitLight, responseTimerReset)
	go openDoor(doorTimeout, doorReset)
	go checkMotorResponse(responseTimerReset)
	go safeStop()
	go updateElevatorMap()
}

func runFSM(newOrder chan bool, newFloor chan int, doorTimeout chan bool, doorReset chan bool, transmitLight chan driver.Button, responseTimerReset chan bool) {
	for {
		select {
		case <-newOrder:
			newOrderInQueue(doorReset, transmitLight)
			responseTimerReset <- true
		case floor := <-newFloor:
			arrivingAtFloor(floor, doorReset, transmitLight)
			responseTimerReset <- true
		case <-doorTimeout:
			doorTimedOut()
			responseTimerReset <- true
		}
	}
}

func newOrderInQueue(doorReset chan bool, transmitLight chan driver.Button) {
	switch config.LocalElev.State {
	case config.IDLE:
		config.LocalElev.Dir = config.LocalElev.Queue.ChooseDir(config.LocalElev.Floor, config.LocalElev.Dir)
		if config.LocalElev.Dir == driver.DIR_STOP {
			config.LocalElev.State = config.DOOR_OPEN
			doorReset <- true
			driver.ElevSetDoorOpenLamp(true)
			config.LocalElev.Queue.ClearOrdersAtFloor(config.LocalElev.Floor, config.LocalElev.Dir)
			orderManager.TransmitLightSignal(transmitLight)
		} else {
			driver.ElevSetMotorDirection(config.LocalElev.Dir)
			config.LocalElev.State = config.MOVING
		}
	case config.MOVING:
		//Do nothing
	case config.DOOR_OPEN:
		if config.LocalElev.Queue.ShouldStop(config.LocalElev.Floor, config.LocalElev.Dir) {
			doorReset <- true
			driver.ElevSetDoorOpenLamp(true)
			config.LocalElev.Queue.ClearOrdersAtFloor(config.LocalElev.Floor, config.LocalElev.Dir)
			orderManager.TransmitLightSignal(transmitLight)
		}
	}
}

func arrivingAtFloor(floor int, doorReset chan bool, transmitLight chan driver.Button) {
	config.LocalElev.Floor = floor
	driver.ElevSetFloorIndicator(floor)
	switch config.LocalElev.State {
	case config.IDLE:
		//Do nothing
	case config.MOVING:
		if config.LocalElev.Queue.ShouldStop(floor, config.LocalElev.Dir) {
			driver.ElevSetMotorDirection(driver.DIR_STOP)
			config.LocalElev.Queue.ClearOrdersAtFloor(floor, config.LocalElev.Dir)
			config.LocalElev.State = config.DOOR_OPEN
			doorReset <- true
			orderManager.TransmitLightSignal(transmitLight)
		}
	case config.DOOR_OPEN:
		//Do nothing
	default:
	}
}
func doorTimedOut() {
	switch config.LocalElev.State {
	case config.IDLE:
		//Do nothing
	case config.MOVING:
		//Do nothing
	case config.DOOR_OPEN:
		driver.ElevSetDoorOpenLamp(false)
		config.LocalElev.Dir = config.LocalElev.Queue.ChooseDir(config.LocalElev.Floor, config.LocalElev.Dir)
		if config.LocalElev.Dir == driver.DIR_STOP {
			config.LocalElev.State = config.IDLE
		} else {
			driver.ElevSetMotorDirection(config.LocalElev.Dir)
			config.LocalElev.State = config.MOVING
		}
	default:
	}
}

func openDoor(doorTimeout, doorReset chan bool) {
	const length = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-doorReset:
			driver.ElevSetDoorOpenLamp(true)
			timer.Reset(length)
		case <-timer.C:
			timer.Stop()
			doorTimeout <- true
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func checkMotorResponse(responseTimerReset chan bool) {
	const length = 5 * time.Second
	timer2 := time.NewTimer(0)
	timer2.Stop()
	for {
		select {
		case <-responseTimerReset:
			timer2.Reset(length)
		case <-timer2.C:
			timer2.Stop()
			if !config.LocalElev.Queue.IsEmpty() && config.LocalElev.State != config.DOOR_OPEN {
				fmt.Println("Motor has stoped")
				driver.ElevSetMotorDirection(driver.DIR_STOP)
				config.LocalElev.Active = false
				for f := 0; f < driver.N_FLOORS; f++ {
					for b := 0; b < driver.N_BUTTONS; b++ {
						driver.ElevSetButtonLamp(b, f, false)
					}
				}
				os.Exit(0)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func updateElevatorMap() {
	for {
		config.ElevatorMap[config.LocalElev.ID] = config.LocalElev
		time.Sleep(100 * time.Millisecond)
	}
}

func safeStop() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.ElevSetMotorDirection(driver.DIR_STOP)
	config.LocalElev.Active = false
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.ElevSetButtonLamp(b, f, false)
		}
	}
	os.Exit(0)
}
