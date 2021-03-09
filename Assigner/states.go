package Assigner

import "../Driver/elevio"
// import "fmt"

const (
	NumFloors  = 4;
	NumButtons = 10;
)

type Order struct {
	Floor int
	Finished bool
	Confirm bool
	Status int 
	ButtonType int
	Id int //??? to identify which elevator has taken the order
	// Examples, do we need this? 
}


type ElevState int
const (
    Idle ElevState = iota
    Moving
    DoorOpen
)

type Elev struct {
	Floor int
	Dir elevio.MotorDirection //Both direction and elevator behaviour in this variable?
	State ElevState
	Online bool
	OrderQueue []bool
	// Order_queue?
}