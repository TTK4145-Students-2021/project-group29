package Executer

import (
	"fmt"
	"os"
	"strings"
	"time"

	hw "../Driver/elevio"

	. "../Common"

	"os"

	"fmt"
)

func InitElev() {
	hw.Init(fmt.Sprintf("localhost:%s", os.Args[1]), NumFloors)
	clearAllLights()
	hw.SetMotorDirection(hw.MD_Down)
	//enten gjør caborders som er på fil eller kjør den under
	// caborders må kontinuerlig bli oppdatert! Gjøres excecuter når vi mottar knappetrykk og når vi gjennomfører en order
	for hw.GetFloor() != 0 {

	}
	hw.SetMotorDirection(hw.MD_Stop)
	hw.SetFloorIndicator(0)

}

func clearAllLights() {
	hw.SetDoorOpenLamp(false)
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
		}
	}
}
func errors(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func cabOrderBackup(elev Elevator) {
	filename := "cabOrder " + GetElevIP() + ".txt"
	f, err := os.Create(filename)
	errors(err)

	caborders := make([]bool, 0)
	for _, row := range elev.OrderQueue {
		caborders = append(caborders, row[NumButtons-1])
	}
	cabordersString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(caborders)), " "), "[]")
	fmt.Println(cabordersString)
	_, err = f.WriteString(cabordersString)

	// other way around https://stackoverflow.com/questions/10783405/how-to-convert-string-into-boolean-array
}

//Moove to localOrderHandler??

func enrollHardware(elev Elevator) {

	hw.SetFloorIndicator(elev.Floor) // Does it harm to set this more times than necessary?
	hw.SetMotorDirection(elev.Dir)
	hw.SetDoorOpenLamp(DOOROPEN == elev.State)

	/*if !elev.Online {
		hw.SetMotorDirection(hw.MD_Stop)
		for i := 0; i < 5; i++ {
			hw.SetStopLamp(true)
			time.Sleep(200 * time.Millisecond)
			hw.SetStopLamp(false)
		}
		hw.SetMotorDirection(elev.Dir)
	}*/
}

func RunElevator(hwChan HardwareChannels, orderChan OrderChannels) {

	// Initializing elevator
	elev := Elevator{
		Id:         "UNDEFINED",
		Floor:      0,
		Dir:        hw.MD_Stop,
		State:      IDLE,
		Online:     true,
		OrderQueue: [NumFloors][NumButtons]bool{},
		Obstructed: false,
	}

	// Executing channels
	// go checkForObstruction(hwChan.HwObstruction, elev)

	// Timer in Go
	doorTimeout := time.NewTimer(3 * time.Second)
	doorTimeout.Stop()

	// Check for engine failure
	engineFailure := time.NewTimer(3 * time.Second)
	engineFailure.Stop()

	var rememberDir hw.MotorDirection
	var numberOfTimeouts = 0
	//var recentEngineFailure = false

	/*ifEqualEmpty := func(a hw.ButtonType, b int) {
		fmt.Println(b)
	} // can this be an empty function of some type?
	*/
	for {
		switch elev.State {
		case IDLE:
			rememberDir = elev.Dir
			select {
			case newOrder := <-orderChan.LocalOrder:
				//fmt.Println("Order recieved of executer")
				elev.Id = newOrder.Id // Gets local ID from Peers
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
					elev.State = MOVING
					elev.Dir = ChooseDirection(elev, rememberDir)
					engineFailure.Reset(3 * time.Second)
				}
				break
			}
		case MOVING:
			select {
			case newOrder := <-orderChan.LocalOrder:
				elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				break
			case newFloor := <-hwChan.HwFloor: //change to elev.Floor := <-hwChan.HwFloor
				elev.Floor = newFloor //remove this?? So that the code is alike
				elev.Online = true

				if ShouldStop(elev) {
					parameters := ClearOrdersParams{Elev: elev}
					elev = ClearOrdersAtCurrentFloor(parameters)
					rememberDir = elev.Dir
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					engineFailure.Stop()
				} else {
					engineFailure.Reset((3 * time.Second)) // If reached floor, reset engineFailure-timer
				}

				break
			case <-engineFailure.C:
				fmt.Println("ENGINE FAILURE")
				elev.Online = false
			}
		case DOOROPEN:
			select {
			case newOrder := <-orderChan.LocalOrder:
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				}
				break
			case <-doorTimeout.C:
				elev.Obstructed = hw.GetObstruction()
				elev.Dir = ChooseDirection(elev, rememberDir)
				//fmt.Printf("%+v\n", elev)
				if elev.Obstructed {
					doorTimeout.Reset(3 * time.Second) // Does the door have to be open 3 seconds after not obstructed????
					elev.State = DOOROPEN
					elev.Dir = hw.MD_Stop
					numberOfTimeouts++
					if numberOfTimeouts == 3 {
						elev.Online = false
						numberOfTimeouts = 0
					}
				} else if elev.Dir == hw.MD_Stop {
					elev.State = IDLE
					elev.Online = true
					//NumElevators--
					engineFailure.Stop()
					numberOfTimeouts = 0
				} else {
					elev.State = MOVING
					elev.Online = true

					//NumElevators--
					engineFailure.Reset((3 * time.Second)) // engineFailure resets whenever an elevator starts moving and has reached a floor.
					numberOfTimeouts = 0
				}
				break
			}
		}

		enrollHardware(elev)
		cabOrderBackup(elev)
		//Implement again when more than one elevator
		orderChan.LocalElevUpdate <- elev // Have to implement these more places?
		// fmt.Println("Orderqueue from local exe: ", elev.OrderQueue)

	}
}
