package Assigner

import (
	"fmt"

	. "../Common"

	hw "../Driver/elevio"

	localip "../Network/network/localip"

	exe "../Executer"

	"os"
)

var AllElevators map[string]Elevator
var OrderBackup map[string][]Order
var SetLights map[string]bool

func GetElevIP() string {
	// Adds elevator-ID (localIP + process ID)
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%s-%d", localIP, os.Getpid())
	return id
}

func AssignOrder(hwChan HardwareChannels, orderChan OrderChannels) {
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
			newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
			//fmt.Println("Sending new order to distributer, via SendOrder")
			//fmt.Printf("%+v\n", newOrder)
			orderChan.SendOrder <- newOrder

		}
	}
}

func UpdateAssigner(orderChan OrderChannels) {
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
}

func PeerUpdate(netChan NetworkChannels) {
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
}

func costFunction(allElev map[string]Elevator, btn hw.ButtonType, floor int) string {
	minTime := -1
	minId := ""
	for id, elev := range allElev {
		if elev.Online {
			time := timeToIdle(elev, btn, floor)
			if time < minTime || minTime == -1 {
				minTime = time
				minId = id
			}
		}
	}
	fmt.Println("Cost function calculated id: ", minId)
	fmt.Println("With minimum time: ", minTime)
	return minId

}

func timeToIdle(elev Elevator, btn hw.ButtonType, floor int) int {
	e := elev
	e.OrderQueue[floor][btn] = true

	arrivedAtRequest := false

	ifEqual := func(inner_b hw.ButtonType, inner_f int) {
		if inner_b == btn && inner_f == floor {
			arrivedAtRequest = true
		}
	}

	duration := 0

	switch elev.State {
	case IDLE:
		elev.Dir = exe.ChooseDirection(elev, elev.Dir)
		if elev.Dir == hw.MD_Stop {
			return duration
		}
		break
	case MOVING:
		duration += TravelTime / 2 //Define travel time later
		elev.Floor += int(elev.Dir)
		break
	case DOOROPEN:
		duration -= DoorOpenTime / 2
	}

	for {
		if exe.ShouldStop(elev) {
			elev = exe.ClearOrdersAtCurrentFloor(elev, ifEqual)
			/* if elev.Dir == hw.MD_Stop {
				return duration
			} */
			if arrivedAtRequest {
				return duration
			}
			duration += DoorOpenTime
			elev.Dir = exe.ChooseDirection(elev, elev.Dir)

		}
		elev.Floor += int(elev.Dir)
		duration += TravelTime
	}
}

func setAllLights() {
	ID := GetElevIP();
	var lightsOff bool
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			for id, elev := range AllElevators {
				SetLights[id] = false
				if btn == hw.BT_Cab && id != ID {
					continue
				}
				if elev.OrderQueue[floor][btn]   { 
					SetLights[id] = true
					hw.SetButtonLamp(hw.ButtonType(btn),floor,true)
				}
			}
			lightsOff = true
			for _, val := range SetLights {		
				if val == true {
					lightsOff = false
				}
			}
			if lightsOff {
				hw.SetButtonLamp(hw.ButtonType(btn),floor,false)
			}
		}
	}
}
