package main

import (
	"config"
	"driver"
	"event_manager"
	"fmt"
	"fsm"
	"network"
	"order_manager/order_manager"
	"time"
)

func main() {

	fmt.Println("You are in main")
	newButton := make(chan driver.Button, 10)
	newOrder := make(chan bool, 10)
	newFloor := make(chan int)
	doorTimeout := make(chan bool)
	doorReset := make(chan bool)
	floorReset := make(chan bool)

	transmitQueue := make(chan config.UDP_queue, 10)
	transmitLight := make(chan driver.Button, 10)

	config.Elev_init()
	fsm.Init(newOrder, newFloor, doorTimeout, doorReset, floorReset, transmitLight)
	network.Init(newButton, newOrder, transmitQueue, transmitLight)

	go event_manager.Get_Button_Press(newButton)
	go event_manager.Get_new_floor(newFloor)
	go order_manager.Order_manager(newButton, newOrder, transmitQueue, transmitLight)

	go func() {
		for {
			config.ElevatorMap[config.LocalElev.ID] = config.LocalElev
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for { //TODO remove this completely
		if driver.Elev_get_stop_signal() == 1 {
			//driver.Elev_set_motor_direction(0)
			fmt.Println(config.LocalElev.Queue)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
