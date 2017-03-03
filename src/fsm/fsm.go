package fsm

import(
  "structs"
  "driver"
  "time"
  "conf"
  "fmt"
)
/**/
//elevator states, TODO move to config
type elevator_state int

const (
  IDLE elevator_state = iota
  MOVING
  DOOR_OPEN
)

var states_test elevator_state

type Fsm_channels struct{ //TODO implement this
    Door_timeout  chan bool
    Door_reset    chan bool
    New_order     chan bool
    New_floor     chan int

}

func open_door(Door_timeout, Door_reset chan bool){
  const length = 3 * time.Second
  timer := time.NewTimer(0)
  timer.Stop()
  for{
    select{
    case <- Door_reset:
      timer.Reset(length)
    case <- timer.C:
      timer.Stop()
      Door_timeout <- true
    }
    time.Sleep(10*time.Millisecond)
  }
}



func Run(newOrder chan bool, newFloor chan int, Door_timeout chan bool, Door_reset chan bool, sync_lights chan bool){
  go open_door(Door_timeout, Door_reset)
  for {
    select{
    case <-newOrder:
      new_order_in_queue(Door_reset, sync_lights)
    case floor := <-newFloor:
      arriving_at_floor(floor, Door_reset, sync_lights)
    case <- Door_timeout:
      door_timeout()
    }
  }
}

func new_order_in_queue(Door_reset chan bool, sync_lights chan bool){
  sync_lights <- true
  fmt.Println("New order in queue")
  switch states_test {
  case IDLE:
    structs.TheElev.Dir = structs.TheElev.Queue.Choose_dir(structs.TheElev.Floor, structs.TheElev.Dir)
    if(structs.TheElev.Dir == conf.STOP){
      states_test = DOOR_OPEN
      Door_reset <- true
      driver.Elev_set_door_open_lamp(true)
    } else{
      driver.Elev_set_motor_direction(structs.TheElev.Dir)
      states_test = MOVING
    }
  case MOVING:
    fmt.Println("Was MOVING - did nothing")
    //do nothing
  case DOOR_OPEN:
    if(structs.TheElev.Dir == conf.STOP){
      Door_reset <- true
      driver.Elev_set_door_open_lamp(true) //don't need this
    }
  }
}

func arriving_at_floor(f int, Door_reset chan bool, sync_lights chan bool){
  structs.TheElev.Floor = f
  driver.Elev_set_floor_indicator(f)
  fmt.Println("Arriving at floor: ", f)
  switch states_test {
  case IDLE:
    //Do nothing
  case MOVING:
    if(structs.TheElev.Queue.Should_stop(f, structs.TheElev.Dir)){
      driver.Elev_set_motor_direction(conf.STOP)
      fmt.Println("Clearing orders for direction: ", structs.TheElev.Dir)
      structs.TheElev.Queue.Clear_orders_at_floor(f, structs.TheElev.Dir)
      sync_lights <- true
      states_test = DOOR_OPEN
      Door_reset<- true
      driver.Elev_set_door_open_lamp(true)
    }
  case DOOR_OPEN:
    //Do nothing
  default:
  }
}
func door_timeout(){
  switch states_test {
  case IDLE:
    //Do nothing
  case MOVING:
    //Do nothing
  case DOOR_OPEN:
    driver.Elev_set_door_open_lamp(false)
    fmt.Println("THIS")
    structs.TheElev.Dir = structs.TheElev.Queue.Choose_dir(structs.TheElev.Floor, structs.TheElev.Dir)

    if(structs.TheElev.Dir == conf.STOP){
      states_test = IDLE
    } else {
      driver.Elev_set_motor_direction(structs.TheElev.Dir)
      states_test = MOVING
    }
  default:
  }
}
