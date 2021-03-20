package Executer

import (
	"time"

	hw "../Driver/elevio"

	. "../Common"
)

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
func checkForObstruction(obstructionChan chan bool, elev Elevator) {
	for {
		select {
		case obstructionButton := <-obstructionChan:
			elev.Obstructed = obstructionButton
		}
	}
}

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
		Floor:      0, // Have to fix this to correct Floor value, but orker ikke nå
		Dir:        hw.MD_Stop,
		State:      IDLE,
		Online:     true,
		OrderQueue: [NumFloors][NumButtons]bool{},
		Obstructed: false,
	}

	// Hardware channels
	go hw.PollButtons(hwChan.HwButtons)
	go hw.PollFloorSensor(hwChan.HwFloor)
	go hw.PollObstructionSwitch(hwChan.HwObstruction)

	// Executing channels
	go checkForObstruction(hwChan.HwObstruction, elev)

	// Timer in Go
	doorTimeout := time.NewTimer(3 * time.Second)
	doorTimeout.Stop()

	// Check for engine failure
	engineFailure := time.NewTimer(3 * time.Second)
	engineFailure.Stop()

	for {
		switch elev.State {
		case IDLE:
			select {
			case newOrder := <-orderChan.LocalOrder:
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					enrollHardware(elev)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
					elev.State = MOVING
					elev.Dir = chooseDirection(elev)
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
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
					clearOrdersAtCurrentFloor(elev)
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
				elev.Dir = chooseDirection(elev)

				if elev.Obstructed {
					doorTimeout.Reset(3 * time.Second)
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
		//orderChan.StateUpdate <- elev // Have to implement these more places?
	}
}
