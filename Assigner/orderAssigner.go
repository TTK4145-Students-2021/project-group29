package Assigner

import (
	. "../Common"
	net "../Distribution"
)

// import "fmt"
// Handles all states
var elevatorInfo Elevator
var allElevatorInfo [NumElevators]Elevator

func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			newOrder := Order{Floor: buttonPress.Floor, Finished: false, Button: buttonPress.Button}
			orderChan.NewOrder <- newOrder

		case updatedLocalElev := <-orderChan.StateUpdate:
			updateElevatorInfo(updatedLocalElev)

		}
	}
}

/*
func costFunc(order Order, elevatorInfo []Elevator) {
	//...
}

*/

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
