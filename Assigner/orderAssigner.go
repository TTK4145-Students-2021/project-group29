package Assigner

import co "../Common"

// import "fmt"

func assignOrder(hwChan co.HardwareChannels, exChan co.ExecuterChannels) {
	// In main
	/*go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr) // Should these go directly into Order_executer? No need for order_assigner to know about obstruction??
	go elevio.PollStopButton(drv_stop)*/

	for {
		select {
		case buttonPress := <-hwChan.hwButtons:
			newOrder := co.Order{Floor: buttonPress.Floor, Finished: false, Button: buttonPress.Button}
			exChan.newOrder <- newOrder
			// Cost func

			// case atFloor := <-hwChan.hwFloor:

			//case obstructionPress := <-hwChan.hwObstruction:

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
