package Executer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	hw "../Driver/elevio"

	. "../Common"
)

func errors(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func InitElev() {

	hw.Init(fmt.Sprintf("localhost:%s", os.Args[1]), NumFloors)

	clearAllLights()

	hw.SetMotorDirection(hw.MD_Down)
	for hw.GetFloor() == -1 {
	}
	hw.SetMotorDirection(hw.MD_Stop)
	hw.SetFloorIndicator(hw.GetFloor())
}

// fuctions to save and read backup of caborders
func writeToBackup(elev Elevator) {
	filename := "cabOrder " + os.Args[1] + ".txt"
	f, err := os.Create(filename)
	errors(err)

	caborders := make([]bool, 0)
	for _, row := range elev.OrderQueue {
		caborders = append(caborders, row[NumButtons-1])
	}
	cabordersString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(caborders)), " "), "[]")
	_, err = f.WriteString(cabordersString)
	defer f.Close()
}

func readFromBackup(orderChan OrderChannels) {
	filename := "cabOrder " + os.Args[1] + ".txt"
	f, err := ioutil.ReadFile(filename)
	errors(err)
	caborders := make([]bool, 0)
	if err == nil {
		s := strings.Split(string(f), " ")
		for _, item := range s {
			result, _ := strconv.ParseBool(item)
			caborders = append(caborders, result)
		}
	}
	id := GetElevIP()
	time.Sleep(15 * time.Millisecond) // A small wait such that my elevator is connected as peer or something with tx message ticker?
	for f, order := range caborders {
		if order {
			newOrder := Order{Floor: f, Button: hw.BT_Cab, Id: id}
			orderChan.SendOrder <- newOrder
		}
	}
}

func clearAllLights() {
	hw.SetDoorOpenLamp(false)
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
		}
	}
}

//Moove to localOrderHandler??

func enrollHardware(elev Elevator) {

	hw.SetFloorIndicator(elev.Floor) // Does it harm to set this more times than necessary?
	hw.SetMotorDirection(elev.Dir)
	hw.SetDoorOpenLamp(DOOROPEN == elev.State)
}

func RunElevator(hwChan HardwareChannels, orderChan OrderChannels) {

	// Initializing elevator
	elev := Elevator{
		Id:         GetElevIP(),
		Floor:      hw.GetFloor(),
		Dir:        hw.MD_Stop,
		State:      IDLE,
		Online:     true,
		OrderQueue: [NumFloors][NumButtons]bool{},
		Obstructed: false,
	}
	// Check if we have backup of cab orders
	readFromBackup(orderChan)

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

	for {
		switch elev.State {
		case IDLE:
			rememberDir = elev.Dir
			select {
			case newOrder := <-orderChan.LocalOrder:
				fmt.Println("Reciving order: ", newOrder)
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
				engineFailure.Reset((3 * time.Second))
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
		writeToBackup(elev)
		//Implement again when more than one elevator
		orderChan.LocalElevUpdate <- elev // Have to implement these more places?
		// fmt.Println("Orderqueue from local exe: ", elev.OrderQueue)

	}
}
