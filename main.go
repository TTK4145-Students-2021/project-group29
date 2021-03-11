package main

import (
	assigner "./Assigner"
	. "./Common"
	hw "./Driver/elevio"
	executer "./Executer"
)

func main() {

	// Making all channels (evt. make a function "InitializeChannels")
	orderChan := OrderChannels{
		NewOrder:    make(chan Order),
		StateUpdate: make(chan Elevator),
	}

	hwChan := HardwareChannels{
		HwButtons:     make(chan hw.ButtonEvent),
		HwFloor:       make(chan int),
		HwObstruction: make(chan bool),
		HwStop:        make(chan bool),
	}

	// Init hardware??
	hw.Init("localhost:15657", NumFloors)

	// Goroutine of runElevator
	go executer.RunElevator(hwChan, orderChan)

	// Goroutine of Assigner
	go assigner.AssignOrder(hwChan, orderChan)

}
