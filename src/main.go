package main

import (
	"fmt"
	"time"
	//"time"
	"driver"
	"order_manager"
	"test"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan driver.Button)

	driver.Elev_init()
	order_manager.Queue_init()
	driver.Elev_set_motor_direction(0)

	go test.Get_Button_Press(newButtonchan)
	go order_manager.Add_order_to_queue(newButtonchan)

	for {
		if driver.Elev_get_stop_signal() == 1 {
			driver.Elev_set_motor_direction(0)
			fmt.Println(order_manager.Queue_matrix)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
