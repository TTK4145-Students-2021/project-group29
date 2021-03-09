package Executer

import  hw "../Driver/elevio"
import "fmt"
import "../Assigner"

func initalizeElevator(Elev elevator) {
	elevator := Elev{
		Floor: hw.getFloor(),
		Dir: MD_STOP,
		State: Idle,
		Online: true,
		OrderQueue: [NumFloors][NumButtons]bool{},
	}
}


func runElevator() {

	elevator Elev
	initializeElevator(elevator)

	for {
		select{
		case newOrder <= newOrder:
			switch elevator.State {
				case Idle:
					if (elevator.Floor == newOrder.Floor) {
						hw.SetDoorOpenLamp(true)
						timerStart(3)
						elevator.State = DoorOpen 
					}
					else {
						elevator.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
						elevator.Dir = localOrderHandler_chooseDirection(elevator) // Denne må implementeres i LocalOrderHandler guro
						hw.SetMotorDirection(elevator.Dir)
						elevator.State = Moving
					}
					break

				case Dooropen:
					if (elevator.Floor == newOrder.Floor) {
						// Timer logikk timerStart
					}
					else {
						elevator.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true
					}
					
					break

				case Moving:
					elevator.OrderQueue[newOrder.Floor][newOrder.ButtonType] = true

					break


			}
		
		
		case newFloor <= arrivedAtFloor:
			elevator.Floor = newFloor
			hw.setFloorIndicator(newFloor)

			switch elevator.State {
				case Moving:
					if localOrderHandler_shouldStop(elevator) {
						elevator.Dir = MD_STOP
						hw.SetMotorDirection(elevator.Dir)
						hw.SetDoorOpenLamp(true)
						// Timer logikk
						// Clear requests at current floor
						elevator.State = DoorOpen
					}
					break
				default:
					break
			}
		case doorTimeout <-  DoorTimeout:

			switch elevator.State {
				case DoorOpen:
					hw.SetDoorOpenLamp(false) // Setting door open lamp to 0 
					elevator.Dir = localOrderHandler_chooseDirection(elevator) // Chooses direction from localOrderHandler
					hw.SetMotorDirection(elevator.Dir) // sets motor direction

					if (elevator.Dir = MD_STOP) {
						elevator.State = Idle
					}
					else {
						elevator.State = Moving
					}
					break

				default:
					break
			}


		}


		
	}
}