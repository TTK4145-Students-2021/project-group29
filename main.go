package main

import (
	. "./Common"
	assigner "./Assigner"
	bcast "./Network/network/bcast"
	distributer "./Distributer"
	executer "./Executer"
	hw "./Driver/elevio"
	peers "./Network/network/peers"
)

func main() {

	executer.InitElev()

	assigner.AllElevs = make(map[string]Elevator)
	assigner.SetLights = make(map[string]bool)
	distributer.PrevRxMsgIds = make(map[string]int)

	orderChan := OrderChannels{
		SendOrder: make(chan Order),
		RecieveElevUpdate: make(chan Elevator),
		LocalOrder: make(chan Order),
		LocalElevUpdate: make(chan Elevator),
	}

	hwChan := HardwareChannels{
		HwButtons: make(chan hw.ButtonEvent),
		HwFloor:       make(chan int),
		HwObstruction: make(chan bool),
	}
	
	netChan := NetworkChannels{
		PeerUpdateCh: make(chan peers.PeerUpdate),
		PeerTxEnable: make(chan bool),		
		BcastMsg:   make(chan Message),
		RecieveMsg: make(chan Message),
		IsOnline: 	  	make(chan bool),
		InmobileElev:	make(chan Elevator),
	}
	
	// Goroutines of Hardware
	go hw.PollButtons(hwChan.HwButtons)
	go hw.PollFloorSensor(hwChan.HwFloor)
	go hw.PollObstructionSwitch(hwChan.HwObstruction)

	// Goroutines of Network
	go bcast.Receiver(42034, netChan.RecieveMsg)
	go bcast.Transmitter(42034, netChan.BcastMsg)
	go peers.Receiver(42035, netChan.PeerUpdateCh)
	go peers.Transmitter(42035, GetElevIP(), netChan.PeerTxEnable)

	// Goroutine of Assigner
	go assigner.Assigner(hwChan, orderChan, netChan)

	// Goroutine of Distibuter
	go distributer.Reciever(netChan, orderChan)
	go distributer.Transmitter(netChan, orderChan)

	// Goroutine of Executer
	go executer.RunElevator(hwChan, orderChan, netChan)
	executer.ReadFromBackup(hwChan)

	select {}
}
