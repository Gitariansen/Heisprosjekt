package main

import (
	"fmt"
	//"time"
	"test"
	"driver"
	"order_manager"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	go test.Get_Button_Press()
	driver.Elev_init()
	order_manager.Queue_init()
	//driver.Elev_set_motor_direction(1)
  //driver.Elev_set_button_lamp(1, 1, true)

	/*  for {
	    for floor := 0; floor < driver.N_FLOORS; floor++{
	      for button := 0; button < driver.N_BUTTONS; button ++{
	        fmt.Println("New iteration", floor, button)
	        fmt.Println(driver.Elev_get_button_signal(button, floor))
	            time.Sleep(100*time.Millisecond)
	      }
	    }

	  }*/

		for {
			if driver.Elev_get_stop_signal() == 1 {
				driver.Elev_set_motor_direction(0)
			}
		}
}
