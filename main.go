package main

import (
	as "./Assigner"
	. "./Common"
	hw "./Driver/elevio"
	ex "./Executer"
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

	// Goroutine of hardware
	go hw.PollButtons(hwChan.HwButtons)
	go hw.PollFloorSensor(hwChan.HwFloor)
	go hw.PollObstructionSwitch(hwChan.HwObstruction)
	go hw.PollStopButton(hwChan.HwStop)

	// Goroutine of runElevator
	go ex.RunElevator(hwChan, orderChan)

	// Goroutine of Assigner
	go as.AssignOrder(hwChan, orderChan)

}
