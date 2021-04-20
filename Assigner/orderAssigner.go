package Assigner

import (
	"fmt"
	. "../Common"
	executer "../Executer"
	hw "../Driver/elevio"
)

var AllElevs map[string]Elevator 
var SetLights map[string]bool
var myId = GetElevIP()

func Assigner(hwChan HardwareChannels, orderChan OrderChannels, netChan NetworkChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			id := "No ID"
			myElev := AllElevs[myId]
			if myElev.Online && NumElevs > 1 {  // More than one elevator is online, distribute order over network 
				if buttonPress.Button == hw.BT_Cab {
					id = myId
				} else {
					id = costFunction(AllElevs, buttonPress.Button, buttonPress.Floor)
				}
				if !duplicateOrder(buttonPress.Button, buttonPress.Floor) {
					newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
					orderChan.SendOrder <- newOrder
				}
			} else { // If no other elevator is online, go in "single-elevator mode", send order directly to local executer
				newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: myId}
				orderChan.LocalOrder <- newOrder
				fmt.Println("SINGLE ELEVATOR MODE")
			}
		case updatedElev := <-orderChan.RecieveElevUpdate:
			AllElevs[updatedElev.Id] = updatedElev
			setAllLights()
		case myElev := <-netChan.InmobileElev:
			AllElevs[myElev.Id] = myElev
			setAllLights()
			reassignOrders(myElev, orderChan) 
		case peer := <-netChan.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peer.Peers)
			fmt.Printf("  New:      %q\n", peer.New)
			fmt.Printf("  Lost:     %q\n", peer.Lost)
			if peer.New != "" {
				if elev, foundPeer := AllElevs[peer.New]; foundPeer { // If elevator is found again, go online
					elev.Online = true
					AllElevs[peer.New] = elev
				} else { 
					elev := Elevator{ // Adding new elevator to AllElevs
						Id:         peer.New,
						Floor:      0,
						Dir:        hw.MD_Stop,
						State:  IDLE,
						Online:     true,
						OrderQueue: [NUMFLOORS][NUMBUTTONS]bool{},
						Mobile:     true,
					}
					AllElevs[peer.New] = elev
				}
				if peer.New == myId {
					netChan.IsOnline <- true
				}
				orderChan.LocalElevUpdate <- AllElevs[myId] // Send update on your own elevator such that newly connected elevators know where it is
				NumElevs++
			}
			if len(peer.Lost) > 0 { // If elevator is lost, going offline
				for _, lostPeer := range peer.Lost { 
					elev := AllElevs[lostPeer]
					elev.Online = false
					if lostPeer == myId {
						netChan.IsOnline <- false
					}
					AllElevs[lostPeer] = elev
					NumElevs--
					reassignOrders(elev, orderChan)
				}
			}
		}
	}
}

func costFunction(allElevs map[string]Elevator, btn hw.ButtonType, floor int) string {
	minTime := -1
	minId := "Undefined"
	for id, elev := range allElevs {
		if elev.Online && elev.Mobile {
			time := timeToServeRequest(elev, btn, floor)
			if time < minTime || minTime == -1 {
				minTime = time
				minId = id
			}
		}
	}
	return minId

}

func sendReassignedOrder(absentElev Elevator, orderChan OrderChannels) {
	for floor := 0; floor < NUMFLOORS; floor++ {
		for btn := 0; btn < NUMBUTTONS-1; btn++ { // Only reassigning hall up and hall down orders (not cab orders)
			if absentElev.OrderQueue[floor][btn] && !duplicateOrder(hw.ButtonType(btn), floor){
				id := costFunction(AllElevs, hw.ButtonType(btn), floor)
				newOrder := Order{Floor: floor, Button: hw.ButtonType(btn), Id: id}
				orderChan.SendOrder <- newOrder
			}
		}
	}
}


func reassignOrders(absentElev Elevator, orderChan OrderChannels) { 
	switch absentElev.Mobile{
	case true:
		for id, elev := range AllElevs {
			if elev.Online {
				if id == myId {
					sendReassignedOrder(absentElev,orderChan)
				}
				break
			}
		}
	case false:
		sendReassignedOrder(absentElev,orderChan)
	}
	
}

func timeToServeRequest(elev Elevator, btn hw.ButtonType, floor int) int {
	elev.OrderQueue[floor][btn] = true
	arrivedAtRequest := false
	ifEqual := func(innerBtn hw.ButtonType, innerFloor int) { // Function that checks if the simulated elevator has reached requested floor and button 
		if innerBtn == btn && innerFloor == floor {
			arrivedAtRequest = true
		}
	}
	duration := 0
	switch elev.State { // Simulation of elevator to find time to serve request
	case IDLE:
		elev.Dir = executer.ChooseDirection(elev, elev.Dir)
		if elev.Dir == hw.MD_Stop {
			return duration
		}
		break
	case MOVING:
		duration += TRAVELTIME / 2
		elev.Floor += int(elev.Dir)
		break
	case DOOROPEN:
		duration -= DOOROPENTIME / 2
	}
	for {
		if executer.ShouldStop(elev) {
			parameters := ClearOrdersParams{Elev: elev, IfEqual: ifEqual}
			elev = executer.ClearOrdersAtCurrentFloor(parameters)
			if arrivedAtRequest {
				return duration
			}
			duration += DOOROPENTIME
			elev.Dir = executer.ChooseDirection(elev, elev.Dir)
		}
		elev.Floor += int(elev.Dir)
		duration += TRAVELTIME
	}
}

func setAllLights() {
	var lightsOff bool
	for floor := 0; floor < NUMFLOORS; floor++ {
		for btn := 0; btn < NUMBUTTONS; btn++ {
			for id, elev := range AllElevs {
				SetLights[id] = false
				if btn == hw.BT_Cab && id != myId { 
					continue
				}
				if elev.OrderQueue[floor][btn] && ((elev.Online && elev.Mobile) || (!elev.Online && id == myId) || (!elev.Mobile && btn == hw.BT_Cab)) {  
					SetLights[id] = true
					hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
				}
			}
			lightsOff = true
			for _, val := range SetLights {  // If one of the elevators has an existing order in SetLights-map, do not turn of lights
				if val == true {
					lightsOff = false
				}
			}
			if lightsOff {
				hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
			}
		}
	}
}

func duplicateOrder(btn hw.ButtonType, floor int) bool {
	for id, elev := range AllElevs {
		if btn == hw.BT_Cab && id != myId {
			continue
		}
		if elev.OrderQueue[floor][btn] && elev.Online && elev.Mobile {
			return true
		}
	}
	return false
}
