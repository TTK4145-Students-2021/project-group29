package Common

import (
	"os"
	"fmt"
	hw "../Driver/elevio"
	p "../Network/network/peers"
	localip "../Network/network/localip"
)

var NumElevs = 0

const (
	NUMFLOORS    		= 4
	NUMBUTTONS   		= 3
	TRAVELTIME   		= 2500
	DOOROPENTIME 		= 3000
	LOSTPACKAGECOUNTER	= 50
)

type ClearOrdersParams struct {
	Elev    	Elevator
	IfEqual 	func(hw.ButtonType, int)
}

type Order struct {
	Floor  		int
	Button 		hw.ButtonType
	Id     		string
}

type ElevatorState int

const (
	IDLE ElevatorState = iota
	MOVING
	DOOROPEN
)

type Elevator struct {
	Id         			string
	Floor      			int
	Dir        			hw.MotorDirection 
	State  			    ElevatorState
	Online     			bool
	OrderQueue 			[NUMFLOORS][NUMBUTTONS]bool 
	Mobile     			bool
}

type HardwareChannels struct {
	HwButtons     		chan hw.ButtonEvent
	HwFloor       		chan int
	HwObstruction 		chan bool
}

type MessageType int

const (
	ORDER MessageType = iota
	ELEVSTATUS
	CONFIRMATION
)

type Message struct {
	OrderMsg    		Order
	ElevMsg 			Elevator
	MsgType     		MessageType
	MsgId   			int
	ElevId  			string
}

type NetworkChannels struct {
	PeerUpdateCh  		chan p.PeerUpdate
	PeerTxEnable   		chan bool
	BcastMsg   			chan Message
	RecieveMsg 			chan Message
	IsOnline 	   		chan bool
	InmobileElev		chan Elevator
}

type OrderChannels struct {
	SendOrder 			chan Order
	RecieveElevUpdate 	chan Elevator
	LocalOrder 			chan Order
	LocalElevUpdate 	chan Elevator
}

func GetElevIP() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%s-%d", localIP, os.Getpid())
	return id
}
