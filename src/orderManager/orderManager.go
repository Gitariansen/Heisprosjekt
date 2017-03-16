package orderManager

import (
	"config"
	"driver"
	"fmt"
	"time"
)

func OrderManager(newButton chan driver.Button, newOrder chan bool, transmitQueue chan config.QueueMessage, transmitLight chan driver.Button) {
	for {
	loop:
		select {
		case buttonPressed := <-newButton:
			if buttonPressed.BtnType == driver.BTN_CMD {
				config.LocalElev.Queue.AddOrderToQueue(buttonPressed)
				newOrder <- true
			} else {
				if config.LocalElev.Active {
					for _, elev := range config.ElevatorMap {
						if elev.Active && elev.Queue.IsOrder(buttonPressed.Floor, buttonPressed.BtnType) {
							break loop
						}
					}
					fmt.Println("----------------------------------------")
					returnID := calculateFastestElevator(buttonPressed)
					if returnID == config.LocalElev.ID {
						fmt.Println("Local elevator was chosen")
						config.LocalElev.Queue.AddOrderToQueue(buttonPressed)
						newOrder <- true
						transmitLight <- buttonPressed
					} else {
						fmt.Println("Elevator with ID: ", returnID, "chosen")
						var tempOrder config.QueueMessage
						tempOrder.IP = returnID
						tempOrder.Button = buttonPressed
						//HERE
						for i = 0; i < 3; i++ {
							transmitQueue <- tempOrder
							time.Sleep(10 * time.Millisecond)
						}

					}
					fmt.Println("----------------------------------------")
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TransmitLightSignal(transmitLight chan driver.Button) {
	var btnUp, btnDown driver.Button
	btnUp, btnDown = config.LocalElev.Queue.ClearLightsAtFloor(config.LocalElev.Floor, config.LocalElev.Dir)
	//HERE
	for i = 0; i < 3; i++ {
		transmitLight <- btnUp
		transmitLight <- btnDown
		time.Sleep(10 * time.Millisecond)
	}
}

func SyncHallLights(transmitLight chan driver.Button) {
	var ret driver.Button
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS-1; b++ {
			ret.Value = config.LocalElev.Queue.IsOrder(f, b)
			ret.Floor = f
			ret.BtnType = b
			transmitLight <- ret
		}
	}
}

func calculateFastestElevator(buttonPressed driver.Button) string {
	returnID := config.LocalElev.ID
	shortestTime := timeToIdle(config.LocalElev, buttonPressed)
	for IP, elev := range config.ElevatorMap {
		if elev.Active {
			elevCopy := elev
			elevCopy.Queue.QueueMatrix[buttonPressed.Floor][buttonPressed.BtnType] = true
			newTime := timeToIdle(elevCopy, buttonPressed)
			if newTime < shortestTime {
				shortestTime = newTime
				returnID = IP
			}
		}
	}
	return returnID
}

/*
Modified code provided by Anders Rønning Petersen
https://piazza.com/class/iyj1u7jrrfe2mo?cid=32
*/

func timeToIdle(elev config.Elevator, buttonPressed driver.Button) time.Duration {
	const travelTime = (2500 * time.Millisecond)
	const doorOpenTime = (3000 * time.Millisecond)
	dur := 0 * time.Millisecond
	elevCopy := elev
	elevCopy.Queue.AddOrderToQueue(buttonPressed)
	switch elevCopy.State {
	case config.IDLE:
		elevCopy.Dir = elevCopy.Queue.ChooseDir(elevCopy.Floor, elevCopy.Dir)
		if elevCopy.Dir == driver.DIR_STOP {
			return dur
		}
	case config.MOVING:
		elevCopy.Floor = elevCopy.Floor + int(elevCopy.Dir)
		dur += travelTime / 2
	case config.DOOR_OPEN:
		dur += doorOpenTime / 2
	}
	for {
		if elevCopy.Queue.ShouldStop(elevCopy.Floor, elevCopy.Dir) {
			elevCopy.Queue.ClearOrdersAtFloor(elevCopy.Floor, elevCopy.Dir)
			dur += doorOpenTime
			elevCopy.Dir = elevCopy.Queue.ChooseDir(elevCopy.Floor, elevCopy.Dir)
			if elevCopy.Dir == driver.DIR_STOP {
				return dur
			}
		}
		elevCopy.Floor = elevCopy.Floor + int(elevCopy.Dir)
		dur += travelTime
		if elevCopy.Floor > 4 || elevCopy.Floor < 1 {
			return dur
		}
	}
}
