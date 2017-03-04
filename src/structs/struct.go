package structs

type Button struct {
	floor  int
	b_type int
}

type Channels struct {
	newButtonchan chan Button
}

type UDP_message struct{}
type TCP_message struct{}

type Fsm_channels struct { //TODO implement this
	Door_timeout chan bool
	Door_reset   chan bool
	New_order    chan bool
	New_floor    chan int
	Sync_lights  chan bool
}
