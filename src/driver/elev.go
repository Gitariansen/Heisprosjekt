package driver

import(
	"conf"
)


type Button struct {
	Floor  int
	B_type int
}

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

const (
	MOTOR_SPEED = 2800
)

type button_type int

const (
	BUTTON_CALL_UP button_type = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}



func Elev_init() {
	Io_init() //CHECK THIS
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			Elev_set_button_lamp(b, f, false)
		}
	}
	if Elev_get_floor_sensor_signal() == -1{
		Elev_set_motor_direction(conf.DOWN)
	}
	for Elev_get_floor_sensor_signal() == -1{

	}
	Elev_set_motor_direction(conf.STOP)
}

func Elev_set_motor_direction(direction int) {
	if direction == 0 {
		Io_write_analog(MOTOR, 0)
	} else if direction > 0 {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	} else if direction < 0 {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_set_button_lamp(button int, floor int, value bool) {
	if value {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}

func Elev_set_floor_indicator(floor int) {
	//binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func Elev_set_door_open_lamp(value bool) {
	if value {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_get_button_signal(button int, floor int) int {
	return int(Io_read_bit(button_channel_matrix[floor][button]))
}

func Elev_get_floor_sensor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) == 1 {
		return 0
	} else if Io_read_bit(SENSOR_FLOOR2) == 1 {
		return 1
	} else if Io_read_bit(SENSOR_FLOOR3) == 1 {
		return 2
	} else if Io_read_bit(SENSOR_FLOOR4) == 1 {
		return 3
	} else {
		return -1
	}
}

func Elev_get_stop_signal() int {
	return Io_read_bit(STOP)
}

func Elev_get_obstruction_signal() int {
	return Io_read_bit(OBSTRUCTION)
}
