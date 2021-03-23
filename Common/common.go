package Common

import hw "../Driver/elevio"

// import "fmt"

const (
	NumFloors    = 4
	NumButtons   = 3
	NumElevators = 1
)

type Order struct {
	Floor  int
	Button hw.ButtonType
	Id     int
}

type ElevState int

const (
	IDLE ElevState = iota
	MOVING
	DOOROPEN
)

type Elevator struct {
	Id         int
	Floor      int
	Dir        hw.MotorDirection //Both direction and elevator behaviour in this variable?
	State      ElevState
	Online     bool
	OrderQueue [NumFloors][NumButtons]bool // Order_queue?
	Obstructed bool
}

type LocalElevChannels struct {
	LocalOrder      chan Order
	LocalElevUpdate chan Elevator
}

// Have not changed these in localOrderExe

type HardwareChannels struct {
	HwButtons     chan hw.ButtonEvent
	HwFloor       chan int
	HwObstruction chan bool
}

type Acknowledge int

const (
	NotAck = iota - 1
	Ack
)

type MessageType struct {
	Order
	Elevator
	Acknowledge
}

type Message struct {
	Msg        MessageType
	MessageId  int
	ElevatorId int
}

type NetworkChannels struct {
	//PeerUpdateCh chan peers.PeerUpdate
	PeerTxEnable   chan bool
	BcastMessage   chan Message
	RecieveMessage chan Message
}

/*
type AssignerChannels struct {
	RecieveElevUpdate chan Elevator
	SendElevUpdate    chan Elevator
	OrderBackupUpdate chan Order
	SendOrder         chan Order
}
*/
type OrderChannels struct {
	//From assigner to distributer
	SendOrder chan Order
	//From distributer to assigner
	OrderBackupUpdate chan Order
	RecieveElevUpdate chan Elevator
	//From distributor to executer
	LocalOrder chan Order
	//From executer to distributor
	LocalElevUpdate chan Elevator
}
