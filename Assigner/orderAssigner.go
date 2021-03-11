package Assigner

import . "../Common"

// import "fmt"

func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			newOrder := Order{Floor: buttonPress.Floor, Finished: false, Button: buttonPress.Button}
			orderChan.NewOrder <- newOrder

			/*case updatedLocalElev := <-exChan.stateUpdate:
			updateElevatorInfo(updatedLocalElev)*/
		}
	}
}

/*
func costFunc(order Order, elevatorInfo []Elevator) {
	//...
}

func getRecommendedExecuter()

func updateElevatorInfo(elev Elevator) {

}

func updateOrderBackup() {
	// Make a map with id and order
}
*/
