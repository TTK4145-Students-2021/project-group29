package main

import (
	hw "./Driver/elevio"
	as "./Assigner"
	ex "./Executer"
	co "./Common"
)

func main() {

	// Making all channels (evt. make a fcuntion "InitializeChannels")
	exChan := co.ExecuterChannels{
		newOrder:       make(chan Order),
		arrivedAtFloor: make(chan int),
		stateUpdate:    make(chan Elevator),
	}

	hwChan := co.HardWareChannels{
		hwButtons:     make(chan hw.buttonEvent),
		hwFloor:       make(chan int),
		hwObstruction: make(chan bool),
		hwStop:        make(chan bool),
	}

	// Init hardware??
	hw.Init("localhost:15657", co.NumFloors)

	// Goroutine of hardware
	go hw.PollButtons(hwChan.hwButtons)
	go hw.PollFloorSensor(hwChan.hwFloor)
	go hw.PollObstructionSwitch(hwChan.hwObstruction)
	go hw.PollStopButton(hwChan.hwStop)

	// Goroutine of runElevator
	go ex.RunElevator(exChan)

	// Goroutine of Assigner
	go as.AssignOrder(hwChan, exChan)

}
