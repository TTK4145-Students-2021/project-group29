package Executer

import (
	"time"

	hw "../Driver/elevio"

	. "../Common"
)

func InitElev() {
	hw.Init("localhost:15657", NumFloors)

	clearAllLights()

	hw.SetMotorDirection(hw.MD_Down)
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

//Moove to localOrderHandler??

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
		Id: 		"UNDEFINED",
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
			rememberDir = elev.Dir
			select {
			case newOrder := <-orderChan.LocalOrder:
				elev.Id = newOrder.Id // Gets local ID from Peers
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
				elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				break
			case newFloor := <-hwChan.HwFloor: //change to elev.Floor := <-hwChan.HwFloor
				elev.Online = true
				elev.Floor = newFloor //remove this?? So that the code is alike
				enrollHardware(elev)

				if shouldStop(elev) {
					elev = clearOrdersAtCurrentFloor(elev)
					rememberDir = elev.Dir
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
					//fmt.Printf("%+v\n", elev)
					// Here we need to set Order to Finished and send it to Assigner, so it can update global map
					engineFailure.Stop()
				} else {
					engineFailure.Reset((3 * time.Second)) // If reached floor, reset engineFailure-timer
				}

				break
			case <-engineFailure.C:
				elev.Online = false
				enrollHardware(elev)
				engineFailure.Reset(5 * time.Second)

			}
		case DOOROPEN:
			select {
			case newOrder := <-orderChan.LocalOrder:
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
				elev.Dir = chooseDirection(elev, rememberDir)
				//fmt.Printf("%+v\n", elev)
				if elev.Obstructed {
					doorTimeout.Reset(3 * time.Second) // Does the door have to be open 3 seconds after not obstructed????
					elev.State = DOOROPEN
					elev.Dir = hw.MD_Stop
					enrollHardware(elev)
				} else if elev.Dir == hw.MD_Stop {
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
		orderChan.LocalElevUpdate <- elev // Have to implement these more places?
	}
}
