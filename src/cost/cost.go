package cost

import (
	"constants"
	"fsm"
	"time"
)

var TravelTime = (2500 * time.Millisecond)
var DoorOpenTime = (3000 * time.Millisecond)

func TimeToIdle(elev fsm.Elevator) time.Duration {
	dur := 0 * time.Millisecond
	e := elev
	switch e.State {
	case fsm.IDLE:
		e.Dir = e.Queue.Choose_dir(e.Floor, e.Dir)
		if e.Dir == constants.STOP {
			return dur
		}
	case fsm.MOVING:
		e.Floor = e.Floor + int(e.Dir)
		dur += TravelTime / 2
	case fsm.DOOR_OPEN:
		dur += DoorOpenTime / 2
	}
	for {
		if e.Queue.Should_stop(e.Floor, e.Dir) {
			e.Queue.Clear_orders_at_floor(e.Floor, e.Dir)
			dur += DoorOpenTime
			e.Dir = e.Queue.Choose_dir(e.Floor, e.Dir)
			if e.Dir == constants.STOP {
				return dur
			}
		}
		e.Floor = e.Floor + int(e.Dir)
		dur += TravelTime
		if e.Floor > 4 || e.Floor < 1 {
			return dur
		}
	}
}
