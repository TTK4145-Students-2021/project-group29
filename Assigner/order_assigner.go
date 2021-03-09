package Assigner

import "./elevio"
import "fmt"

// Jeg har bare loka litt her for å forstå hva vi trenger, ikke noe kode vi trenger å bruke :) 


type Order struct {
	floor int
	finished bool
	confirm bool
	status int 
	button_type int
	id int //??? to identify which elevator has taken the order
	// Examples
}


type ElevState int
const (
    Idle ElevState = iota
    Moving
    DoorOpen
) // This is an enum in Go


type Elev struct {
	Floor int
	Dir MotorDirection //Both direction and elevator behaviour in this variable?
	State ElevState
	Online bool
	Order_queue []Order
}

drv_buttons := make(chan elevio.ButtonEvent)
drv_floors  := make(chan int)
drv_obstr   := make(chan bool)
drv_stop    := make(chan bool)    

go elevio.PollButtons(drv_buttons)
go elevio.PollFloorSensor(drv_floors)
go elevio.PollObstructionSwitch(drv_obstr) // Should these go directly into Order_executer? No need for order_assigner to know about obstruction??
go elevio.PollStopButton(drv_stop)

func assign_order() {
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