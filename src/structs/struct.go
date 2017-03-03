package structs

import (
	"order_manager"
	"driver"
	"conf"
	"fmt"
)

var TheElev Elevator

func Elev_init_own() {
	driver.Io_init()
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.Elev_set_button_lamp(b, f, false)
		}
	}
	if driver.Elev_get_floor_sensor_signal() == -1{
		driver.Elev_set_motor_direction(conf.DOWN)
	}
	for driver.Elev_get_floor_sensor_signal() == -1{
		//wait
	}
	driver.Elev_set_motor_direction(conf.STOP)
	TheElev.Floor = driver.Elev_get_floor_sensor_signal()
	driver.Elev_set_floor_indicator(TheElev.Floor)
	TheElev.Queue = order_manager.Make_empty_queue() //TODO change make-empty-queue to update-queue
	fmt.Println("Elevator is initialized in floor: ", TheElev.Floor+1)
}

type Button struct {
  floor int
  b_type int
}

type Channels struct {
  newButtonchan chan Button
}

type Elevator struct {
  Floor int
  State int
  Queue order_manager.Queue
  ID    int //int?
	Dir		int
}

type UDP_message struct {}
type TCP_message struct {}
