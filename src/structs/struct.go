package structs

import (
	"order_manager"
	"driver"
	"conf"
	"fmt"
)

var TheElev Elevator

func Elev_init_own() {
	driver.Io_init() //CHECK THIS
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			driver.Elev_set_button_lamp(b, f, false)
		}
	}
	if driver.Elev_get_floor_sensor_signal() == -1{
		driver.Elev_set_motor_direction(conf.DOWN)
	}
	for driver.Elev_get_floor_sensor_signal() == -1{

	}
	driver.Elev_set_motor_direction(conf.STOP)
	TheElev.Floor = driver.Elev_get_floor_sensor_signal()
	fmt.Println("Elevators floor is ", TheElev.Floor)
	driver.Elev_set_floor_indicator(TheElev.Floor)
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
  ID    int //?
	Dir		int
}

type UDP_message struct {}
type TCP_message struct {}
