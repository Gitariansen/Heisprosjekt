package eventManager

import (
	"driver"
	"time"
)

func GetButtonPress(newButton chan driver.Button) {
	var buttonPressed driver.Button
	for {
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for button := 0; button < driver.N_BUTTONS; button++ {
				if driver.ElevGetButtonSignal(button, floor) == 1 {
					buttonPressed.Floor = floor
					buttonPressed.BtnType = button
					buttonPressed.Value = true
					newButton <- buttonPressed
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func GetNewFloor(newFloor chan int) {
	prevFloor := driver.ElevGetFloorSensorSignal()
	for {
		currFloor := driver.ElevGetFloorSensorSignal()
		if currFloor != -1 && currFloor != prevFloor {
			newFloor <- currFloor
		}
		prevFloor = currFloor
	}
}
