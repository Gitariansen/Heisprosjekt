package structs

type Button struct {
	Floor  int
	B_type int
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

type UDP_queue struct {
	IP     string
	Button Button
}
