package Assigner

import "../Driver/elevio"
// import "fmt"

// Jeg har bare loka litt her for å forstå hva vi trenger, ikke noe kode vi trenger å bruke :) 


func assign_order() {
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	drv_obstr   := make(chan bool)
	drv_stop    := make(chan bool)    

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr) // Should these go directly into Order_executer? No need for order_assigner to know about obstruction??
	go elevio.PollStopButton(drv_stop)
	for {
		select {
		case a := <- drv_buttons:
		
		case a := <- drv_floors:

		case a := <- drv_obstr:

		case a:= <- drv_stop:
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