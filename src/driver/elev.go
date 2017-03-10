package driver

const (
	DIR_DOWN = -1
	DIR_UP   = 1
	DIR_STOP = 0
)

const (
	BTN_UP   = 0
	BTN_DOWN = 1
	BTN_CMD  = 2
)

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

const (
	MOTOR_SPEED = 2800
)

type buttonType int

const (
	BUTTON_CALL_UP buttonType = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type Button struct {
	Floor   int
	BtnType int
	Value   bool
}

var lampChannelMatrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannelMatrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func ElevSetMotorDirection(direction int) {
	if direction == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if direction > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	} else if direction < 0 {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	}
}

func ElevSetButtonLamp(button int, floor int, value bool) {
	if value {
		IoSetBit(lampChannelMatrix[floor][button])
	} else {
		IoClearBit(lampChannelMatrix[floor][button])
	}
}

func ElevSetFloorIndicator(floor int) {
	//binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		IoSetBit(LIGHT_FLOOR_IND1)
	} else {
		IoClearBit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 > 0 {
		IoSetBit(LIGHT_FLOOR_IND2)
	} else {
		IoClearBit(LIGHT_FLOOR_IND2)
	}
}

func ElevSetDoorOpenLamp(value bool) {
	if value {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevGetButtonSignal(button int, floor int) int {
	return int(IoReadBit(buttonChannelMatrix[floor][button]))
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR_FLOOR1) == 1 {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2) == 1 {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3) == 1 {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4) == 1 {
		return 3
	} else {
		return -1
	}
}

func ElevGetStopSignal() int {
	return IoReadBit(STOP)
}
