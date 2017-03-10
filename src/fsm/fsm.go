package fsm

import (
	"config"
	"driver"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func Init(newOrder chan bool, newFloor chan int, doorTimeout chan bool, doorReset chan bool, floorReset chan bool, transmitLight chan driver.Button) {
	go runFSM(newOrder, newFloor, doorTimeout, doorReset, transmitLight, floorReset)
	go openDoor(doorTimeout, doorReset)
	go checkMotorResponse(floorReset)
	go safeKill()
}

func runFSM(newOrder chan bool, newFloor chan int, Door_timeout chan bool, Door_reset chan bool, light_clear chan driver.Button, floor_reset chan bool) {
	for {
		select {
		case <-newOrder:
			newOrderInQueue(Door_reset, light_clear)
			floor_reset <- true //TODO RENAME THIS sHIT
		case floor := <-newFloor:
			arrivingAtFloor(floor, Door_reset, light_clear)
			floor_reset <- true
		case <-Door_timeout:
			doorTimeout()
			floor_reset <- true
		}
	}
}

func newOrderInQueue(Door_reset chan bool, light_clear chan driver.Button) {
	switch config.LocalElev.State {
	case config.IDLE:
		config.LocalElev.Dir = config.LocalElev.Queue.Choose_dir(config.LocalElev.Floor, config.LocalElev.Dir)
		if config.LocalElev.Dir == driver.DIR_STOP {
			config.LocalElev.State = config.DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			config.LocalElev.Queue.Clear_orders_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			//TODO make subfunction
			var bup, bdwn driver.Button
			bup, bdwn = config.LocalElev.Queue.Clear_lights_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			light_clear <- bup
			light_clear <- bdwn
		} else {
			driver.Elev_set_motor_direction(config.LocalElev.Dir)
			config.LocalElev.State = config.MOVING
		}
	case config.MOVING:
		//do nothing
	case config.DOOR_OPEN:
		if config.LocalElev.Queue.Should_stop(config.LocalElev.Floor, config.LocalElev.Dir) {
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			config.LocalElev.Queue.Clear_orders_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			var bup, bdwn driver.Button
			bup, bdwn = config.LocalElev.Queue.Clear_lights_at_floor(config.LocalElev.Floor, config.LocalElev.Dir)
			light_clear <- bup
			light_clear <- bdwn
		}
	}
}

func arrivingAtFloor(f int, Door_reset chan bool, light_clear chan driver.Button) {
	config.LocalElev.Floor = f
	var bup, bdwn driver.Button
	driver.Elev_set_floor_indicator(f)
	switch config.LocalElev.State {
	case config.IDLE:
		//Do nothing
	case config.MOVING:
		if config.LocalElev.Queue.Should_stop(f, config.LocalElev.Dir) {
			driver.Elev_set_motor_direction(driver.DIR_STOP)
			config.LocalElev.Queue.Clear_orders_at_floor(f, config.LocalElev.Dir)

			config.LocalElev.State = config.DOOR_OPEN
			Door_reset <- true
			driver.Elev_set_door_open_lamp(true)
			bup, bdwn = config.LocalElev.Queue.Clear_lights_at_floor(f, config.LocalElev.Dir)

			light_clear <- bup

			light_clear <- bdwn

		}
	case config.DOOR_OPEN:
		//Do nothing
	default:
	}
}
func doorTimeout() {
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
func openDoor(Door_timeout, Door_reset chan bool) {
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

func checkMotorResponse(floor_reset chan bool) {
	const length = 5 * time.Second
	timer2 := time.NewTimer(0)
	timer2.Stop()
	for {

		fmt.Println("1")
		select {
		case <-floor_reset:
			fmt.Println("Resetting timer")
			timer2.Reset(length)
		case <-timer2.C:
			fmt.Println("Gother")
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

func safeKill() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(driver.DIR_STOP)
	config.LocalElev.Active = false
	os.Exit(0)
}
