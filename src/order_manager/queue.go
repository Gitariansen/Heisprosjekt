package order_manager

import (
	"driver"
	"fmt"
	"conf"
)

type Queue struct {
	queue_matrix [driver.N_FLOORS][driver.N_BUTTONS]bool
}

func Make_empty_queue() Queue{
	var ret Queue
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			ret.queue_matrix[f][b] = false //TODO The init must check the backup, and update correctly
		}
	}
	return ret
}

func (q *Queue)Add_order_to_queue(c chan driver.Button, o chan bool) {

	for {
		select {
		case button_pressed := <-c:
			q.queue_matrix[button_pressed.Floor][button_pressed.B_type] = true
			q.set_lights()
			o <- true
		}
	}
}

func (q *Queue)Clear_orders_at_floor(floor, dir int) {
	q.queue_matrix[floor][conf.B_CMD] = false
	switch dir{
	case conf.UP:
		q.queue_matrix[floor][conf.B_UP] = false
		if!(q.order_above(floor)){
			q.queue_matrix[floor][conf.B_DOWN] = false
		}
	case conf.DOWN:
	  q.queue_matrix[floor][conf.B_DOWN] = false
		if!(q.order_below(floor)){
			q.queue_matrix[floor][conf.B_UP] = false
		}
	case conf.STOP:
		q.queue_matrix[floor][conf.B_UP] = false
		q.queue_matrix[floor][conf.B_DOWN] = false
	}
	q.set_lights()
}

func (q *Queue)set_lights(){
	  for f:=0;f<driver.N_FLOORS;f++{
	    for b:= 0; b<driver.N_BUTTONS;b++{
	      if q.queue_matrix[f][b]{
	        driver.Elev_set_button_lamp(b,f,true)
	      } else {
	        driver.Elev_set_button_lamp(b,f,false)
	      }
	    }
	  }
}

func (q *Queue)Should_stop(floor int, dir int)bool{
	switch dir{
	case conf.UP:
		if(q.queue_matrix[floor][conf.B_UP] || q.queue_matrix[floor][conf.B_CMD] || floor == driver.N_FLOORS-1 || !q.order_above(floor)){
			return true
		}
	case conf.DOWN:
		if(q.queue_matrix[floor][conf.B_DOWN] || q.queue_matrix[floor][conf.B_CMD] || floor == 0 || !q.order_below(floor)){
			return true
		}
	case conf.STOP:
		for b := 0; b < driver.N_BUTTONS; b++{
			if(q.queue_matrix[floor][b]){
				return true
			}
		}
	default:
		return false
	}
	return false
}

func (q *Queue)order_above(floor int)bool{
	for f := floor+1; f < driver.N_FLOORS; f++{
 		for b := 0; b < driver.N_BUTTONS; b++{
 			if(q.queue_matrix[f][b]){
 				return true
 			}
 		}
	}
	return false
}

func (q *Queue)order_below(floor int)bool{
	for f := floor - 1; f >= 0; f--{
 		for b := 0; b < driver.N_BUTTONS; b++{
 			if(q.queue_matrix[f][b]){
 				return true
 			}
 		}
	}
	return false
}

 func (q *Queue)is_empty()bool{
	 for f:= 0; f<driver.N_FLOORS; f++{
		 for b:=0;b<driver.N_BUTTONS; b++{
			 if(q.queue_matrix[f][b]){
				 return false
			 }
		 }
	 }
	 return true
 }

func (q *Queue)Choose_dir(floor, dir int) int{ //THIS IS NOT COMMPLETE TODO
	if q.is_empty(){
		return conf.STOP
	}
	switch dir{
	case conf.UP:
		if q.order_above(floor){
			return conf.UP
		}else if q.order_below(floor){
			return conf.DOWN
		}else{
			return conf.STOP
		}
	case conf.DOWN:
		if q.order_below(floor){
			return conf.DOWN
		}else if q.order_above(floor){
			return conf.UP
		} else {
			return conf.STOP
		}

	case conf.STOP:
		if q.order_above(floor){
			return conf.UP
		}else if q.order_below(floor){
			return conf.DOWN
		}else{
			return conf.STOP
		}
	default:
		return conf.STOP
	}
}
