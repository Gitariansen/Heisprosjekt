package cost

import (
	"config"
	"driver"
	"time"
)

//TODO move to order_manager
var TravelTime = (2500 * time.Millisecond)
var DoorOpenTime = (3000 * time.Millisecond)

func TimeToIdle(elev config.Elevator) time.Duration {
	dur := 0 * time.Millisecond
	e := elev
	switch e.State {
	case config.IDLE:
		e.Dir = e.Queue.Choose_dir(e.Floor, e.Dir)
		if e.Dir == driver.DIR_STOP {
			return dur
		}
	case config.MOVING:
		e.Floor = e.Floor + int(e.Dir)
		dur += TravelTime / 2
	case config.DOOR_OPEN:
		dur += DoorOpenTime / 2
	}
	for {
		if e.Queue.Should_stop(e.Floor, e.Dir) {
			e.Queue.Clear_orders_at_floor(e.Floor, e.Dir)
			dur += DoorOpenTime
			e.Dir = e.Queue.Choose_dir(e.Floor, e.Dir)
			if e.Dir == driver.DIR_STOP {
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
