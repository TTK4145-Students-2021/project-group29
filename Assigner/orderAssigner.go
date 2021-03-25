package Assigner

import (
	"fmt"

	. "../Common"

	hw "../Driver/elevio"
)

// import "fmt"
// Handles all states
var elevatorInfo Elevator
var allElevators [NumElevators]Elevator

func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			// Cost function returning ID of elevator taking the order

			newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: 123}
			//fmt.Printf("%+v\n", newOrder)
			orderChan.SendOrder <- newOrder
			/* Implement again when more elevators
			case peerUpdate := PeerHandler:
				// Reassign all orders here
				AssignerChannels.SendOrder <- newOrder
			*/
		}
	}
}

func UpdateAssigner(orderChan OrderChannels) {
	for {
		select {
		case updateLocalElev := <-orderChan.RecieveElevUpdate:
			allElevators[0] = updateLocalElev
			fmt.Printf("%v", allElevators)
			setAllLights(allElevators[0])

		}
	}
}

func setAllLights(elev Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] == true {
				hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
			} else {
				hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
			}
		}
	}
}

/*
func UpdateAssigner() {
	for {
		select {
		case updateLocalElev := <-LocalElevChannels.LocalElevUpdate:
			updateLocalElevator(updateLocalElev)
			AssignerChannels.SendElevUpdate <- updateLocalElev

		case updateExternalElev := <-AssignerChannels.RecieveElevUpdate:
			updateElevators(updateExternalElev)

		case updateOrderList := <-AssignerChannels.OrderBackupUpdate:
			updateOrderBackup(updateOrderList)

		}
	}
}
*/
/*
func costFunc(order Order, elevatorInfo []Elevator) {
	//...
}

*/

/*
func getRecommendedExecuter() {

}



func updateAllElevatorInfo(msg net.Message) {

}

func updateOrderBackup() {
	// Make a map with id and order
}


func RemoveElevFromNetwork() {
	// If PeersUpdate (p.Lost)
	// Remove that elevator from network
	// Only needs ID to the elevator that is lost
	elev = ElevList[ID]
	PeerHandler <- elev
	// Sends all orders to AssignCer through a channel
}
func AddElevToNetwork() {
	// If PeersUpdate (p.New)
	// Add Elevator to network
}
*/
