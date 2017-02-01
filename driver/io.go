package driver

/* need io.go, io.c, io,h, channels.go, channels.h and driver.go*/
import "C"

func io_init(){
  C.io_init()
}

func io_set_bit(channel int){
  C.io_set_bit(C.int(channel))
}

func io_clear_bit(channel int){
  C.io_clear_bit(C.int(channel))
}

func io_read_bit(channel int) int {
  return int(C.io_read_bit(C.int(channel)))
}

func io_read_analog(channel int) int {
  return int(C.oi_read_analog(C.int(channel)))
}

func io_write_analog(channel int, value int){
  C.io_write_analog(C.int(channel),C.int(value))
}
