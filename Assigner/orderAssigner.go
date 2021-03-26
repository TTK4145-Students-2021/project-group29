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
			// Cost function returning ID of elevator taking the order
			id := costFunction(AllElevators)
			newOrder := Order{Floor: buttonPress.Floor, Button: buttonPress.Button, Id: id}
			fmt.Println("Sending new order to distributer, via SendOrder")
			fmt.Printf("%+v\n", newOrder)
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
			}
			if len(p.Lost) > 0 {
				for _, lostPeer := range p.Lost { // If elevator is lost, going offline
					elev := AllElevators[lostPeer]
					elev.Online = false
					AllElevators[lostPeer] = elev
				}
			}
		}
	}
}

func costFunction(allElev map[string]Elevator) string {
	for id, _ := range allElev {
		if id != GetElevIP() {
			return id
		}
	}
	return "error"
}

func timeToIdle(elev Elevator) int {
	duration := 0
	switch elev.State {
	case IDLE:
		elev.Dir = exe.ChooseDirection(elev, elev.Dir)
		if elev.Dir == hw.MD_Stop {
			return duration
		}
		break
	case MOVING:
		duration += TRAVELTIME / 2 //Define travel time later
		elev.Floor += int(elev.Dir)
		break
	case DOOROPEN:
		duration -= 3 / 2
	}
	for {
		if exe.ShouldStop(elev) {
			elev = exe.ClearOrdersAtCurrentFloor(elev)
			duration += 3
			elev.Dir = exe.ChooseDirection(elev, elev.Dir)
			if elev.Dir == hw.MD_Stop {
				return duration
			}
		}
		elev.Floor += int(elev.Dir)
		duration += TRAVELTIME
	}
}

func setAllLights(elev Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] == true {
				hw.SetButtonLamp(hw.ButtonType(btn), floor, true)
			} else {
				hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
			}
		}
	}
}

/*
func UpdateAssigner() {
	for {
		select {
		case updateLocalElev := <-LocalElevChannels.LocalElevUpdate:
			updateLocalElevator(updateLocalElev)
			AssignerChannels.SendElevUpdate <- updateLocalElev

		case updateExternalElev := <-AssignerChannels.RecieveElevUpdate:
			updateElevators(updateExternalElev)

		case updateOrderList := <-AssignerChannels.OrderBackupUpdate:
			updateOrderBackup(updateOrderList)

		}
	}
}
*/
/*
func costFunc(order Order, elevatorInfo []Elevator) {
	//...
}

*/

/*
func getRecommendedExecuter() {

}



func updateAllElevatorInfo(msg net.Message) {

}

func updateOrderBackup() {
	// Make a map with id and order
}


func RemoveElevFromNetwork() {
	// If PeersUpdate (p.Lost)
	// Remove that elevator from network
	// Only needs ID to the elevator that is lost
	elev = ElevList[ID]
	PeerHandler <- elev
	// Sends all orders to AssignCer through a channel
}
func AddElevToNetwork() {
	// If PeersUpdate (p.New)
	// Add Elevator to network
}
*/
