package fsm

import(
  "order_manager"
)


//elevator states, TODO move to config

const (
  IDLE elevator_state = iota
  MOVING
  DOOR_OPEN
)

type Channels struct{

}
/*
const (
	BUTTON_CALL_UP elev_button_type_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)
*/

func arriving_at_floor(){
  switch elevator_state {
  case IDLE:
    //Do nothing
  case MOVING:
    if order_manager.Should_stop(){
      driver.Elev_set_motor_direction(0)
      //DOOR OPEN
      //set stoplights
    }

  case DOOR_OPEN:
    //Do nothing
  default:
}
func door_timeout(){
  switch elevator_state {
  case IDLE:
    //Do nothing
  case MOVING:
    //Do nothing
  case DOOR_OPEN:
    //call a set_direction func. set_direction should put elevator in correct state
  default:
}
func new_order_in_queue(){
  switch elevator_state {
  case IDLE:
    //check if the order is in this elevator's current floor. If it is: call open_door func, switch state to DOOR_OPEN
  case MOVING:
    //do nothing
  case DOOR_OPEN:
    //check if the order is in this elevator's current floor. If it is: call open_door func, switch state to DOOR_OPEN
  default:
  }

}
