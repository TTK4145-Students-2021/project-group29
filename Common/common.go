package Common

import hw "../Driver/elevio"

// import "fmt"

const (
	NumFloors  = 4
	NumButtons = 10
)

type Order struct {
	Floor     int
	Finished  bool
	// Confirmed bool
	Button    hw.ButtonType
	// Id int 
}

type ElevState int

const (
	Idle ElevState = iota
	Moving
	DoorOpen
)

type Elevator struct {
	Floor      int
	Dir        hw.MotorDirection //Both direction and elevator behaviour in this variable?
	State      ElevState
	Online     bool
	OrderQueue [NumFloors][NumButtons]bool // Order_queue?
}

type ExecuterChannels struct {
	newOrder       chan Order
	arrivedAtFloor chan int
	stateUpdate    chan Elevator
}

type HardwareChannels struct {
	hwButtons     chan hw.ButtonEvent
	hwFloor       chan int
	hwObstruction chan bool
	hwStop        chan bool
}
