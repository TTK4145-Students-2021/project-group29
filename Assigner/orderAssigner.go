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
			// Cost function returning ID of elevator taking the order
			if buttonPress.Button == hw.BT_Cab {
				id = GetElevIP()
			} else {
				id = costFunction(AllElevators, buttonPress.Button, buttonPress.Floor)
			}
			if !duplicateOrder(buttonPress.Button, buttonPress.Floor) {
				newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
				orderChan.SendOrder <- newOrder
			}
		case newOrder := <-orderChan.OrderBackupUpdate:
			OrderBackup[newOrder.Id] = append(OrderBackup[newOrder.Id], newOrder)
			
		case updatedElev := <-orderChan.RecieveElevUpdate:
			fmt.Println("Recieving elev update")
			AllElevators[updatedElev.Id] = updatedElev
			setAllLights()
			if !updatedElev.Online && updatedElev.Id == GetElevIP() {
				netChan.PeerTxEnable <- false
			} else if updatedElev.Online && updatedElev.Id == GetElevIP() {
				netChan.PeerTxEnable <- true
				fmt.Println("Enabling peer again")
			}	
		case peer := <-netChan.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peer.Peers)
			fmt.Printf("  New:      %q\n", peer.New)
			fmt.Printf("  Lost:     %q\n", peer.Lost)
			fmt.Println(len(peer.Lost))
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
						Obstructed: false,
					}
					AllElevators[peer.New] = elev
				}
				NumElevators++
			}
			if len(peer.Lost) > 0 {
				for _, lostPeer := range peer.Lost { // If elevator is lost, going offline
					elev := AllElevators[lostPeer]
					elev.Online = false
					AllElevators[lostPeer] = elev
					NumElevators--
					id := "No ID"
					fmt.Println("Recieved lost peer")
					for floor := 0; floor < NumFloors; floor++ {
						for btn := 0; btn < NumButtons-1; btn++ {
							if elev.OrderQueue[floor][btn] {
								id = costFunction(AllElevators, hw.ButtonType(btn),floor)
								if !duplicateOrder(hw.ButtonType(btn), floor) {
									newOrder := Order{Floor: floor, Button: hw.ButtonType(btn), Id: id}
									fmt.Println("Sending reassigned order!")
									orderChan.SendOrder <- newOrder
								}
							}
						}
					}

				}
			}
		
		}
	}
}

/* func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
	for {
		select {
		case buttonPress := <-hwChan.HwButtons:
			id := "No ID"
			// Cost function returning ID of elevator taking the order
			if buttonPress.Button == hw.BT_Cab {
				id = GetElevIP()
			} else {
				id = costFunction(AllElevators, buttonPress.Button, buttonPress.Floor)
			}
			if !duplicateOrder(buttonPress.Button, buttonPress.Floor) {

				newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
				orderChan.SendOrder <- newOrder
			}

		}
	}
}*/

/*func UpdateAssigner(orderChan OrderChannels) {
	for {
		select {
		case newOrder := <-orderChan.OrderBackupUpdate:
			OrderBackup[newOrder.Id] = append(OrderBackup[newOrder.Id], newOrder)
			// Make function that deletes orders from backup when finished
		case updatedElev := <-orderChan.RecieveElevUpdate:
			AllElevators[updatedElev.Id] = updatedElev
			setAllLights()

		}
	}
}*/

/*func PeerUpdate(netChan NetworkChannels) {
	for {
		select {
		case p := <-netChan.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			if p.New != "" {
				if elev, found := AllElevators[p.New]; found { // If elevator is found again, going online
					elev.Online = true
					AllElevators[p.New] = elev

				} else { // If elevator is new, needs to be created
					elev := Elevator{
						Id:         p.New,
						Floor:      0,
						Dir:        hw.MD_Stop,
						State:      IDLE,
						Online:     true,
						OrderQueue: [NumFloors][NumButtons]bool{},
						Obstructed: false,
					}
					AllElevators[p.New] = elev
				}
				NumElevators++
			}
			if len(p.Lost) > 0 {
				for _, lostPeer := range p.Lost { // If elevator is lost, going offline
					elev := AllElevators[lostPeer]
					elev.Online = false
					AllElevators[lostPeer] = elev
					NumElevators--
				}
			}
		}
	}
}*/

func costFunction(allElev map[string]Elevator, btn hw.ButtonType, floor int) string {
	minTime := -1
	minId := ""
	for id, elev := range allElev {
		if elev.Online {
			time := timeToServeRequest(elev, btn, floor)
			if time < minTime || minTime == -1 {
				minTime = time
				minId = id
			}
		}
	}
	return minId

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
		duration += TravelTime / 2 //Define travel time later
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
				
				if elev.OrderQueue[floor][btn] && (elev.Online ||  btn == hw.BT_Cab) { // make this better 
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
				fmt.Println("Turning lights off: ", btn, floor)
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
		if elev.OrderQueue[floor][btn] && elev.Online {
			fmt.Println("Returning TRUE to duplicate")
			return true
		}
	}
	fmt.Println("Returning FALSE to duplicate")
	return false
}