package Assigner

import hw "../Driver/elevio"
// import "fmt"

// Jeg har bare loka litt her for å forstå hva vi trenger, ikke noe kode vi trenger å bruke :) 

type HardwareChannels struct {
	hwButtons chan hw.ButtonEvent
	hwFloor chan int
	hwObstruction chan bool
	hwStop chan bool
}

func assign_order(ch HardwareChannels) {
	// In main 
	/*go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr) // Should these go directly into Order_executer? No need for order_assigner to know about obstruction??
	go elevio.PollStopButton(drv_stop)*/

	for {
		select {
		case buttonPress := <- ch.hwButtons:
			newOrder := Order{Floor: buttonPress.Floor, Finished: false, Confirmed: false, Button: buttonPress.Button}
			// Cost func

		case atFloor := <- ch.hwFloor:

		case obstructionPress := <- ch.hwObstruction:

		case stopPress := <- ch.hwStop:
		}
	}
}

func cost_func(order Order, elevator_info []Elev) {
	//...
}

func update_elevator_info() {

}

func update_order_backup() {

}