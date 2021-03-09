package Executer

import  hw "../Driver/elevio"
import "fmt"
import "../Assigner"
import "time"

type ExecuterChannels struct {
	newOrder chan Order
	arrivedAtFloor chan int
	stateUpdate chan Elevator

}

func initalizeElevator() {
	elev := Elevator{
		Floor: hw.getFloor(),
		Dir: MD_STOP,
		State: Idle,
		Online: true,
		OrderQueue: [NumFloors][NumButtons]bool{},
	}

	return elev
}


func runElevator(ch ExecuterChannels) {

	Elevator elev := initializeElevator()
	ch.stateUpdate <- elev

	// Timer in Go 
	doorTimeout := time.NewTimer(3 * time.Second)
	doorTimeout.Stop()



	for {
		select{
		case newOrder <= ch.newOrder:
			switch elev.State {
				case Idle:
					if (elev.Floor == newOrder.Floor) {
						hw.SetDoorOpenLamp(true)	
						elev.State = DoorOpen 
					}
					else {
						elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
						elev.Dir = localOrderHandler_chooseDirection(elev) // Denne må implementeres i LocalOrderHandler guro
						hw.SetMotorDirection(elev.Dir)
						elev.State = Moving
					}
					break

				case Dooropen:
					if (elev.Floor == newOrder.Floor) {
						doorTimeout.Reset(3 * time.Second)
					}
					else {
						elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
					}
					
					break

				case Moving:
					elev.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
					break


			}
			ch.stateUpdate <- elev
		
		
		case newFloor <= ch.arrivedAtFloor:
			elev.Floor = newFloor
			hw.setFloorIndicator(newFloor)

			switch elev.State {
				case Moving:
					if localOrderHandler_shouldStop(elev) { // Feil syntaks?
						elev.Dir = MD_STOP
						hw.SetMotorDirection(elev.Dir)
						hw.SetDoorOpenLamp(true)
						doorTimeout.Reset(3 * time.Second)
						// Clear requests at current floor
						elev.State = DoorOpen
					}
					break
				default:
					break
			}
			ch.stateUpdate <- elev

		case <-doorTimeout.C:

			switch elev.State {
				case DoorOpen:
					hw.SetDoorOpenLamp(false) // Setting door open lamp to 0 
					elev.Dir = localOrderHandler_chooseDirection(elev) // Chooses direction from localOrderHandler
					hw.SetMotorDirection(elev.Dir) // sets motor direction

					if (elev.Dir = MD_STOP) {
						elev.State = Idle
					}
					else {
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