package Executer

import (
	"fmt"
	"time"

	hw "../Driver/elevio"

	. "../Common"
)

func InitElev() {
	hw.Init("localhost:15657", NumFloors)

	hw.SetMotorDirection(hw.MD_Down)
	for hw.GetFloor() != 0 {

	}
	hw.SetMotorDirection(hw.MD_Stop)
	hw.SetFloorIndicator(0)
	hw.SetDoorOpenLamp(false)
}

//Moove to localOrderHandler??
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

// Where should we use this?
/*func checkForObstruction(obstructionChan chan bool, elev Elevator) {
	for {
		select {
		case obstructionButton := <-obstructionChan:
			fmt.Println("checking checking cheking obstruction")
			elev.Obstructed = obstructionButton
		}
	}
}*/

func enrollHardware(elev Elevator) {

	hw.SetFloorIndicator(elev.Floor) // Does it harm to set this more times than necessary?
	hw.SetMotorDirection(elev.Dir)

	switch elev.State {
	case DOOROPEN:
		hw.SetDoorOpenLamp(true)
	case MOVING:
		hw.SetDoorOpenLamp(false)
	case IDLE:
		hw.SetDoorOpenLamp(false)
	}

	if !elev.Online {
		hw.SetMotorDirection(hw.MD_Stop)
		for i := 0; i < 5; i++ {
			hw.SetStopLamp(true)
			time.Sleep(200 * time.Millisecond)
			hw.SetStopLamp(false)
		}
		hw.SetMotorDirection(elev.Dir)
	}
}

func RunElevator(hwChan HardwareChannels, orderChan OrderChannels) {

	// Initializing elevator
	elev := Elevator{
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

	for {
		switch elev.State {
		case IDLE:
			fmt.Println("Inside IDLE")
			rememberDir = elev.Dir
			select {
			case newOrder := <-orderChan.LocalOrder:
				fmt.Println("Inside IDLE --- New order!")
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
					elev.State = MOVING
					elev.Dir = chooseDirection(elev, rememberDir)
					enrollHardware(elev)
					engineFailure.Reset(3 * time.Second)
				}
				break
			}
		case MOVING:
			select {
			case newOrder := <-orderChan.LocalOrder:
				fmt.Println("Inside MOVING --- New order, adding to queue!")
				elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				break
			case newFloor := <-hwChan.HwFloor: //change to elev.Floor := <-hwChan.HwFloor
				fmt.Println("Inside MOVING --- Detecting floor!")
				elev.Online = true
				elev.Floor = newFloor //remove this?? So that the code is alike
				enrollHardware(elev)

				if shouldStop(elev) {
					fmt.Println("Stopping!")
					elev = clearOrdersAtCurrentFloor(elev)
					rememberDir = elev.Dir
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
					fmt.Printf("%+v\n", elev)
					// Here we need to set Order to Finished and send it to Assigner, so it can update global map
					engineFailure.Stop()
				} else {
					engineFailure.Reset((3 * time.Second)) // If reached floor, reset engineFailure-timer
				}

				break
			case <-engineFailure.C:
				fmt.Printf("Inside MOVING --- Engine failure!")
				elev.Online = false
				enrollHardware(elev)
				engineFailure.Reset(5 * time.Second)

			}
		case DOOROPEN:
			select {
			case newOrder := <-orderChan.LocalOrder:
				fmt.Println("Inside DOOROPEN --- New order!")
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				}
				break
			case <-doorTimeout.C:
				elev.Obstructed = hw.GetObstruction()
				fmt.Println("Inside DOOROPEN --- Door timeout!")
				elev.Dir = chooseDirection(elev, rememberDir)
				fmt.Printf("%+v\n", elev)
				if elev.Obstructed {
					fmt.Println("OBSTRUCTED!!!!!")
					doorTimeout.Reset(3 * time.Second) // Does the door have to be open 3 seconds after not obstructed????
					elev.State = DOOROPEN
					elev.Dir = hw.MD_Stop
					enrollHardware(elev)
				} else if elev.Dir == hw.MD_Stop {
					fmt.Println("Inside    elev.Dir == hw.MD_Stop")
					elev.State = IDLE
					engineFailure.Stop()
					enrollHardware(elev)
				} else {
					elev.State = MOVING
					engineFailure.Reset((3 * time.Second)) // engineFailure resets whenever an elevator starts moving and has reached a floor.
					enrollHardware(elev)
				}
				break
			}
		}
		//Implement again when more than one elevator
		//orderChan.StateUpdate <- elev // Have to implement these more places?
	}
}
