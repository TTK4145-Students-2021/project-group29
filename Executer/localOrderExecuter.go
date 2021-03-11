package Executer

import (
	"time"

	hw "../Driver/elevio"

	. "../Common"
)

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
}

func RunElevator(hwChan HardwareChannels, orderChan OrderChannels) {

	// Initializing elevator
	elev := Elevator{
		Floor:      0, // Have to fix this to correct Floor value
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
	// go checkForNewOrders(orderChan.newOrder, elev)

	// Timer in Go
	doorTimeout := time.NewTimer(3 * time.Second)
	doorTimeout.Stop()

	for {
		switch elev.State {
		case IDLE:
			select {
			case newOrder := <-orderChan.NewOrder:
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					enrollHardware(elev)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
					elev.State = MOVING
					elev.Dir = chooseDirection(elev)
					enrollHardware(elev)
				}
				break
			}
		case MOVING:
			select {
			case newOrder := <-orderChan.NewOrder:
				elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				break
			case newFloor := <-hwChan.HwFloor:
				elev.Floor = newFloor
				enrollHardware(elev)

				if shouldStop(elev) {
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(3 * time.Second)
					enrollHardware(elev)
					clearOrdersAtCurrentFloor(elev)
				}
				break
			}
		case DOOROPEN:
			select {
			case newOrder := <-orderChan.NewOrder:
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
					enrollHardware(elev)
				} else {
					elev.State = MOVING
					enrollHardware(elev)
				}
				break
			}

		}
		orderChan.StateUpdate <- elev
	}
}
