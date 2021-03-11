package Common

import hw "../Driver/elevio"

// import "fmt"

const (
	NumFloors    = 4
	NumButtons   = 3
	NumElevators = 3
)

type Order struct {
	Floor    int
	Finished bool
	Button   hw.ButtonType
}

type ElevState int

const (
	IDLE ElevState = iota
	MOVING
	DOOROPEN
)

type Elevator struct {
	Floor      int
	Dir        hw.MotorDirection //Both direction and elevator behaviour in this variable?
	State      ElevState
	Online     bool
	OrderQueue [NumFloors][NumButtons]bool // Order_queue?
	Obstructed bool
}

type OrderChannels struct {
	NewOrder    chan Order
	StateUpdate chan Elevator
}

type HardwareChannels struct {
	HwButtons     chan hw.ButtonEvent
	HwFloor       chan int
	HwObstruction chan bool
	HwStop        chan bool
}
