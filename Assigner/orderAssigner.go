package Assigner

import (
	"fmt"

	. "../Common"

	hw "../Driver/elevio"

	exe "../Executer"
)

var AllElevators map[string]Elevator
var OrderBackup map[string][]Order
var SetLights map[string]bool

func Assigner(hwChan HardwareChannels, orderChan OrderChannels, netChan NetworkChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			id := "No ID"
			myId := GetElevIP()
			myElev := AllElevators[myId]
			if myElev.Online && NumElevators > 1 { // If no pwwer (not ourself) is detected (packet loss =1) or If one we only have ourself
				if buttonPress.Button == hw.BT_Cab {
					id = GetElevIP()
				} else {
					id = costFunction(AllElevators, buttonPress.Button, buttonPress.Floor)
				}
				if !duplicateOrder(buttonPress.Button, buttonPress.Floor) {
					newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
					orderChan.SendOrder <- newOrder
				}
			} else {
				newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: myId}
				fmt.Println(newOrder)
				orderChan.LocalOrder <- newOrder
				fmt.Println("Going into SINGLE MODE!!!!")
			}
		case newOrder := <-orderChan.OrderBackupUpdate:
			OrderBackup[newOrder.Id] = append(OrderBackup[newOrder.Id], newOrder)

		case updatedElev := <-orderChan.RecieveElevUpdate:
			AllElevators[updatedElev.Id] = updatedElev
			setAllLights()

		case myElev := <-netChan.InMobileElev:
			fmt.Println("Reassigning inmobile elev")
			AllElevators[myElev.Id] = myElev
			setAllLights()
			reassignOrders(myElev, orderChan) // assume that it is not offline

		case peer := <-netChan.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peer.Peers)
			fmt.Printf("  New:      %q\n", peer.New)
			fmt.Printf("  Lost:     %q\n", peer.Lost)
			if peer.New != "" {
				if elev, foundPeer := AllElevators[peer.New]; foundPeer { // If elevator is found again, going online
					elev.Online = true
					AllElevators[peer.New] = elev
				} else { // If elevator is new, needs to be created
					elev := Elevator{
						Id:         peer.New,
						Floor:      0,
						Dir:        hw.MD_Stop,
						State:      IDLE,
						Online:     true,
						OrderQueue: [NumFloors][NumButtons]bool{},
						Mobile:     true,
					}
					AllElevators[peer.New] = elev
				}
				if peer.New == GetElevIP() {
					netChan.IsOnline <- true
					fmt.Println("Sending true")
				}
				orderChan.LocalElevUpdate <- AllElevators[GetElevIP()] // send update on myself such that newly connected elevators know where I am
				NumElevators++
			}
			if len(peer.Lost) > 0 {
				for _, lostPeer := range peer.Lost { // If elevator is lost, going offline
					elev := AllElevators[lostPeer]
					elev.Online = false
					if lostPeer == GetElevIP() {
						netChan.IsOnline <- false
					}
					AllElevators[lostPeer] = elev
					NumElevators--
					reassignOrders(elev, orderChan)
				}
			}
		}
	}
}

func costFunction(allElev map[string]Elevator, btn hw.ButtonType, floor int) string {
	minTime := -1
	minId := ""
	for id, elev := range allElev {
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

func reassignOrders(offlineElev Elevator, orderChan OrderChannels) { // change name
	myId := GetElevIP()
	for id, elev := range AllElevators {
		if elev.Online {
			if id == myId {
				for floor := 0; floor < NumFloors; floor++ {
					for btn := 0; btn < NumButtons-1; btn++ {
						if offlineElev.OrderQueue[floor][btn] {
							id := costFunction(AllElevators, hw.ButtonType(btn), floor)
							if !duplicateOrder(hw.ButtonType(btn), floor) {
								newOrder := Order{Floor: floor, Button: hw.ButtonType(btn), Id: id}
								fmt.Println("Sending new order to channel")
								orderChan.SendOrder <- newOrder
							}
						}
					}
				}
			}
			break
		}
	}
}

func timeToServeRequest(elev Elevator, btn hw.ButtonType, floor int) int {
	e := elev
	e.OrderQueue[floor][btn] = true

	arrivedAtRequest := false
	ifEqual := func(inner_b hw.ButtonType, inner_f int) {
		if inner_b == btn && inner_f == floor {
			arrivedAtRequest = true
		}
	}

	duration := 0

	switch e.State {
	case IDLE:
		e.Dir = exe.ChooseDirection(e, e.Dir)
		if e.Dir == hw.MD_Stop {
			return duration
		}
		break
	case MOVING:
		duration += TravelTime / 2
		e.Floor += int(e.Dir)
		break
	case DOOROPEN:
		duration -= DoorOpenTime / 2
	}

	for {
		if exe.ShouldStop(e) {
			parameters := ClearOrdersParams{Elev: e, IfEqual: ifEqual}
			e = exe.ClearOrdersAtCurrentFloor(parameters)
			if arrivedAtRequest {
				return duration
			}
			duration += DoorOpenTime
			e.Dir = exe.ChooseDirection(e, e.Dir)

		}
		e.Floor += int(e.Dir)
		duration += TravelTime
	}
}

func setAllLights() {
	ID := GetElevIP()
	var lightsOff bool
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			for id, elev := range AllElevators {
				SetLights[id] = false
				if btn == hw.BT_Cab && id != ID {
					continue
				}

				if elev.OrderQueue[floor][btn] && elev.Online && elev.Mobile { // make this better
					SetLights[id] = true
					hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
				}
				if elev.OrderQueue[floor][btn] && !elev.Online && id == ID {
					SetLights[id] = true
					hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
				}
				if elev.OrderQueue[floor][btn] && !elev.Mobile && btn == hw.BT_Cab {
					SetLights[id] = true
					hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
				}
			}
			lightsOff = true
			for _, val := range SetLights {
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
	ID := GetElevIP()
	for id, elev := range AllElevators {
		if btn == hw.BT_Cab && id != ID {
			continue
		}
		if elev.OrderQueue[floor][btn] && elev.Online && elev.Mobile {
			fmt.Println("true")
			return true
		}
	}
	return false
}
