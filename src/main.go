package main

import (
	"fmt"
	//"time"
	//../network"
	//"../test"
	"driver"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	//go test.Get_Button_Press()
	driver.Elev_init()
	driver.Elev_set_motor_direction(-1)

	for {
		if driver.Elev_get_stop_signal() == 1 {
			driver.Elev_set_motor_direction(0)
		}
	}
	/*  for {
	    for floor := 0; floor < driver.N_FLOORS; floor++{
	      for button := 0; button < driver.N_BUTTONS; button ++{
	        fmt.Println("New iteration", floor, button)
	        fmt.Println(driver.Elev_get_button_signal(1, 1))
	        driver.Elev_set_button_lamp(1, 1, false)
	            time.Sleep(100*time.Millisecond)
	      }
	    }

	  }*/
}
