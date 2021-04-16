package main

import (
	assigner "./Assigner"
	. "./Common"
	distributer "./Distribution"
	hw "./Driver/elevio"
	executer "./Executer"
	bcast "./Network/network/bcast"
	peers "./Network/network/peers"
)

func main() {

	executer.InitElev()

	assigner.AllElevators = make(map[string]Elevator)
	assigner.OrderBackup = make(map[string][]Order)
	assigner.SetLights = make(map[string]bool)
	distributer.PrevRxMsgIDs = make(map[string]int)
	// Making all channels (evt. make a function "InitializeChannels")

	orderChan := OrderChannels{
		//From assigner to distributer
		SendOrder: make(chan Order),
		//From distributer to assigner
		OrderBackupUpdate: make(chan Order),
		RecieveElevUpdate: make(chan Elevator),
		//From distributor to executer
		LocalOrder: make(chan Order),
		//From executer to distributor
		LocalElevUpdate: make(chan Elevator),
		//ReassignOrders:  make(chan string),

	}

	hwChan := HardwareChannels{
		//From elevio to orderassigner
		HwButtons: make(chan hw.ButtonEvent),
		//From elevio to executer
		HwFloor:       make(chan int),
		HwObstruction: make(chan bool),
		// HwStop:        make(chan bool), //Implement this later
	}
	netChan := NetworkChannels{
		//Between OrderAssigner and Network
		PeerUpdateCh: make(chan peers.PeerUpdate),
		PeerTxEnable: make(chan bool),
		//Between OrderDistributor and Network
		BcastMessage:   make(chan Message),
		RecieveMessage: make(chan Message),
	}
	// Goroutines of Hardware
	go hw.PollButtons(hwChan.HwButtons)
	go hw.PollFloorSensor(hwChan.HwFloor)
	go hw.PollObstructionSwitch(hwChan.HwObstruction)

	// Goroutines of Network
	go bcast.Receiver(42034, netChan.RecieveMessage)
	go bcast.Transmitter(42034, netChan.BcastMessage)
	go peers.Receiver(42035, netChan.PeerUpdateCh)
	go peers.Transmitter(42035, GetElevIP(), netChan.PeerTxEnable)

	// Goroutine of Assigner
	go assigner.Assigner(hwChan, orderChan, netChan)

	// Goroutine of Distibuter
	go distributer.Reciever(netChan, orderChan)
	go distributer.Transmitter(netChan, orderChan)

	// Goroutine of runElevator, in executer
	go executer.RunElevator(hwChan, orderChan)

	select {}
}
