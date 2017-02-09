package order_manager

import (
  "driver"
)


type queue_matrix struct{
  Matrix [driver.N_FLOORS][driver.N_BUTTONS]bool
}



func Get_queue(){}

func Add_order_to_queue(){}

func Clear_orders_at_floor(){}

func Queue_init(){
  for f := 0; f< driver.N_FLOORS; f++{
    for b := 0; b< driver.N_BUTTONS; b++{
      //queue_matrix[f][b] = false //TODO this must be changed for fault tolerance
    }
  }
}
