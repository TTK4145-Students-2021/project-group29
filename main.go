package main

import (
	assigner "./Assigner"
	. "./Common"
	hw "./Driver/elevio"
	executer "./Executer"
	"./Network/network/peers"
	"./Network/network/bcast"
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

	netChan := NetworkChannels {
		PeerUpdateCh: make(chan peers.PeerUpdate),
		PeerTxEnable: make(chan bool),
		BcastMessage: make(chan Message),
		RecieveMessage: make(chan Message)

	} 

	// Init hardware??
	hw.Init("localhost:15657", NumFloors)

	// Goroutine of runElevator
	go executer.RunElevator(hwChan, orderChan)

	// Goroutine of Assigner
	go assigner.AssignOrder(hwChan, orderChan)

	// Goroutine from Network module
	go peers.Reciever(42035, netChan.PeerUpdateCh) 
	go peers.Transmitter(42035), netChan.PeerTxEnable)

	go bcast.Reciever(42034, netChan.RecieveMessage)
	go bcast.Transmitter(42034, netChan.BcastMessage)

	select{}
}
