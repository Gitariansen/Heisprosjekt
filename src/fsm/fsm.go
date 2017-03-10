package fsm

import (
	"config"
	"driver"
	"fmt"
	"order_manager/order_manager"
	"os"
	"os/signal"
	"time"
)

func Init(newOrder chan bool, newFloor chan int, doorTimeout chan bool, doorReset chan bool, floorReset chan bool, transmitLight chan driver.Button) {
	go runFSM(newOrder, newFloor, doorTimeout, doorReset, transmitLight, floorReset)
	go openDoor(doorTimeout, doorReset)
	go checkMotorResponse(floorReset)
	go safeKill()
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
		config.LocalElev.Dir = config.LocalElev.Queue.Choose_dir(config.LocalElev.Floor, config.LocalElev.Dir)
		if config.LocalElev.Dir == driver.DIR_STOP {
			config.LocalElev.State = config.DOOR_OPEN
			doorReset <- true
			driver.Elev_set_door_open_lamp(true)
			config.LocalElev.Queue.Clear_orders_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			order_manager.TransmitLightSignal(transmitLight)
		} else {
			driver.Elev_set_motor_direction(config.LocalElev.Dir)
			config.LocalElev.State = config.MOVING
		}
	case config.MOVING:
		//do nothing
	case config.DOOR_OPEN:
		if config.LocalElev.Queue.Should_stop(config.LocalElev.Floor, config.LocalElev.Dir) {
			doorReset <- true
			driver.Elev_set_door_open_lamp(true)
			config.LocalElev.Queue.Clear_orders_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			order_manager.TransmitLightSignal(transmitLight)
		}
	}
}

func arrivingAtFloor(floor int, doorReset chan bool, transmitLight chan driver.Button) {
	config.LocalElev.Floor = floor
	driver.Elev_set_floor_indicator(floor)
	switch config.LocalElev.State {
	case config.IDLE:
		//Do nothing
	case config.MOVING:
		if config.LocalElev.Queue.Should_stop(floor, config.LocalElev.Dir) {
			driver.Elev_set_motor_direction(driver.DIR_STOP)
			config.LocalElev.Queue.Clear_orders_at_floor(floor, config.LocalElev.Dir)
			config.LocalElev.State = config.DOOR_OPEN
			doorReset <- true
			order_manager.TransmitLightSignal(transmitLight)
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
		driver.Elev_set_door_open_lamp(false)
		config.LocalElev.Dir = config.LocalElev.Queue.Choose_dir(config.LocalElev.Floor, config.LocalElev.Dir)
		if config.LocalElev.Dir == driver.DIR_STOP {
			config.LocalElev.State = config.IDLE
		} else {
			driver.Elev_set_motor_direction(config.LocalElev.Dir)
			config.LocalElev.State = config.MOVING
		}
	default:
	}
}

//hardkokt
func openDoor(Door_timeout, doorReset chan bool) {
	const length = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-doorReset:
			timer.Reset(length)
		case <-timer.C:
			timer.Stop()
			Door_timeout <- true
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func checkMotorResponse(floor_reset chan bool) {
	const length = 5 * time.Second
	timer2 := time.NewTimer(0)
	timer2.Stop()
	for {
		select {
		case <-floor_reset:
			timer2.Reset(length)
		case <-timer2.C:
			timer2.Stop()
			if !config.LocalElev.Queue.Is_empty() && config.LocalElev.State != config.DOOR_OPEN {
				fmt.Println("Motor has stoped")
				driver.Elev_set_motor_direction(driver.DIR_STOP)
				config.LocalElev.Active = false
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
func safeKill() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(driver.DIR_STOP)
	config.LocalElev.Active = false
	os.Exit(0)
}
