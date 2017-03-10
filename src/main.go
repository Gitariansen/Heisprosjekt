package main

import (
	"config"
	"driver"
	"event_manager"
	"fsm"
	"network"
	"order_manager/order_manager"
	"os"
)

func main() {

	newButton := make(chan driver.Button, 10)
	newOrder := make(chan bool, 10)
	newFloor := make(chan int)
	doorTimeout := make(chan bool)
	doorReset := make(chan bool)
	responseTimerReset := make(chan bool)

	transmitQueue := make(chan config.UDP_queue, 10)
	transmitLight := make(chan driver.Button, 10)

	config.Elev_init()
	fsm.Init(newOrder, newFloor, doorTimeout, doorReset, responseTimerReset, transmitLight)
	network.Init(transmitQueue, transmitLight)

	go event_manager.GetButtonPress(newButton)
	go event_manager.GetNewFloor(newFloor)
	go network.HandleIncomingMessages(newButton, newOrder, transmitQueue, transmitLight)
	go order_manager.OrderManager(newButton, newOrder, transmitQueue, transmitLight)

	for {
		if driver.Elev_get_stop_signal() == 1 {
			os.Exit(0)
		}
	}
}
