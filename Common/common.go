package Common

import (
	"os"

	hw "../Driver/elevio"
	p "../Network/network/peers"

	"fmt"

	localip "../Network/network/localip"
)

// import "fmt"

func GetElevIP() string {
	// Adds elevator-ID (localIP + process ID)
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%s-%d", localIP, os.Getpid())
	return id
}

const (
	NumFloors    = 4
	NumButtons   = 3
	TravelTime   = 2500
	DoorOpenTime = 3000
)

var NumElevators = 0

type ClearOrdersParams struct {
	Elev    Elevator
	IfEqual func(hw.ButtonType, int)
}

type Order struct {
	Floor  int
	Button hw.ButtonType
	Id     string
}

type ElevState int

const (
	IDLE ElevState = iota
	MOVING
	DOOROPEN
)

type Elevator struct {
	Id         string
	Floor      int
	Dir        hw.MotorDirection //Both direction and elevator behaviour in this variable?
	State      ElevState
	Online     bool
	OrderQueue [NumFloors][NumButtons]bool // Order_queue?
	Mobile     bool
}

type HardwareChannels struct {
	HwButtons     chan hw.ButtonEvent
	HwFloor       chan int
	HwObstruction chan bool
}

type MessageType int

const (
	ORDER MessageType = iota
	ELEVSTATUS
	CONFIRMATION
)

type Message struct {
	OrderMsg    Order
	ElevatorMsg Elevator
	MsgType     MessageType
	MessageId   int
	ElevatorId  string
}

type NetworkChannels struct {
	PeerUpdateCh  	chan p.PeerUpdate
	PeerTxEnable   	chan bool
	BcastMessage   	chan Message
	RecieveMessage 	chan Message
	IsOnline 	   	chan bool
	InMobileElev		chan Elevator
}

type OrderChannels struct {
	//From assigner to distributer
	SendOrder 			chan Order
	//From distributer to assigner
	OrderBackupUpdate  	chan Order
	RecieveElevUpdate 	chan Elevator
	//From distributor to executer
	LocalOrder 			chan Order
	//From executer to distributor
	LocalElevUpdate 	chan Elevator
	//ReassignOrders  chan string
}
