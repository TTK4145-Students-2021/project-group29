package Assigner

import (
	"fmt"

	. "../Common"
)

// halla
// import "fmt"
// Handles all states
var elevatorInfo Elevator
var allElevators [NumElevators]Elevator

func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			// Lots of cost functions
			// Send newOrder to Distribution
			// Id = 1 // Here we find id to the one taking the order
			newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: 123}
			fmt.Printf("%+v\n", newOrder)
			orderChan.SendOrder <- newOrder
			fmt.Println("Sending order from assigner to distributer")
			/* Implement again when more elevators
			case peerUpdate := PeerHandler:
				// Reassign all orders here
				AssignerChannels.SendOrder <- newOrder
			*/
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

func updateElevatorInfo(elev Elevator) {
	elevatorInfo.Floor = elev.Floor
	elevatorInfo.Dir = elev.Dir
	elevatorInfo.State = elev.State
	elevatorInfo.Online = elev.Online
	elevatorInfo.OrderQueue = elev.OrderQueue

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
