package main

import (
	"fmt"
	"time"
	"driver"
	"test"
	"structs"
	"fsm"
)

func main() {
	//network.init()
	fmt.Println("You are in main")
	newButtonchan := make(chan driver.Button)
	newOrderchan	:= make(chan bool)
	newFloorchan  := make(chan int)
	doorTimeoutchan := make(chan bool)
	doorResetchan := make(chan bool)



	structs.Elev_init_own()

	//structs.TheElev.Queue = order_manager.Make_empty_queue() //TODO change make-empty-queue to update-queue

	go test.Get_Button_Press(newButtonchan)
	go test.Get_new_floor(newFloorchan)
	go structs.TheElev.Queue.Add_order_to_queue(newButtonchan, newOrderchan)
	go fsm.Run(newOrderchan, newFloorchan, doorTimeoutchan, doorResetchan)
	//go structs.TheElev.Queue.Set_lights()

	for {
		if driver.Elev_get_stop_signal() == 1 {
			driver.Elev_set_motor_direction(0)
			fmt.Println(structs.TheElev.Queue)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
