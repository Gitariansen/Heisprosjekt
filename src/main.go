package main

import (
	"config"
	"driver"
	"eventManager"
	"fmt"
	"fsm"
	"network"
	"orderManager"
	"os"
)

func main() {

	newButton := make(chan driver.Button, 10)
	newOrder := make(chan bool, 10)
	newFloor := make(chan int)
	doorTimeout := make(chan bool)
	doorReset := make(chan bool)
	responseTimerReset := make(chan bool)
	transmitQueue := make(chan config.QueueMessage, 10)
	transmitLight := make(chan driver.Button, 10)

	config.ElevInit()
	fsm.Init(newOrder, newFloor, doorTimeout, doorReset, responseTimerReset, transmitLight)
	network.Init(transmitQueue, transmitLight)

	go eventManager.GetButtonPress(newButton)
	go eventManager.GetNewFloor(newFloor)
	go network.HandleIncomingMessages(newButton, newOrder, transmitQueue, transmitLight)
	go orderManager.OrderManager(newButton, newOrder, transmitQueue, transmitLight)

	for {
		if driver.ElevGetStopSignal() == 1 {
			fmt.Println("Program terminated")
			os.Exit(0)
		}
	}
}
