package Executer

import (
	"time"

	hw "../Driver/elevio"

	co "../Common"
)

func initalizeElevator() {
	floorChan := make(chan int)
	elev := co.Elevator{
		hw.PollFloorSensor(floorChan),
		0,
		0,
		true,
		[co.NumFloors][co.NumButtons]bool{},
	}

	return elev

}

func RunElevator(ch co.ExecuterChannels) {

	elev := initializeElevator()

	// Timer in Go
	doorTimeout := time.NewTimer(3 * time.Second)
	doorTimeout.Stop()

	for {
		select {
		case newOrder := <-ch.newOrder:
			switch elev.State {
			case co.Idle:
				if elev.Floor == newOrder.Floor {
					hw.SetDoorOpenLamp(true)
					elev.State = DoorOpen
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
					elev.Dir = chooseDirection(elev)
					hw.SetMotorDirection(elev.Dir)
					elev.State = Moving
				}
				break

			case co.Dooropen:
				if elev.Floor == newOrder.Floor {
					doorTimeout.Reset(3 * time.Second)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
				}

				break

			case co.Moving:
				elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
				break

			}

		case newFloor := <-ch.arrivedAtFloor: // This channel do not need to be connected to orderAssigner
			elev.Floor = newFloor
			hw.SetFloorIndicator(newFloor)

			switch elev.State {
			case co.Moving:
				if shouldStop(elev) { //
					elev.Dir = hw.MD_STOP
					hw.SetMotorDirection(elev.Dir)
					hw.SetDoorOpenLamp(true)
					doorTimeout.Reset(3 * time.Second)

					elev.State = DoorOpen
					clearOrdersAtCurrentFloor(elev)
					// Implement logic concering Order-states (Finished?)

				}
				break
			default:
				break
			}

		case <-doorTimeout.C:

			switch elev.State {
			case DoorOpen:
				hw.SetDoorOpenLamp(false)        // Setting door open lamp to 0
				elev.Dir = chooseDirection(elev) // Chooses direction from localOrderHandler
				hw.SetMotorDirection(elev.Dir)   // sets motor direction

				if elev.Dir == MD_STOP {
					elev.State = Idle
				} else {
					elev.State = Moving
				}
				break

			default:
				break
			}
			ch.stateUpdate <- elev
		}

	}
}
