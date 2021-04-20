package Executer

import (
	"fmt"
	"os"
	"time"

	. "../Common"
	hw "../Driver/elevio"
)

func InitElev() {
	hw.Init(fmt.Sprintf("localhost:%s", os.Args[1]), NUMFLOORS)
	clearAllLights()

	hw.SetMotorDirection(hw.MD_Down) // Moving down to the closest floor if elevator is in-between floors
	for hw.GetFloor() == -1 {
	}
	hw.SetMotorDirection(hw.MD_Stop)
	hw.SetFloorIndicator(hw.GetFloor())
}

func RunElevator(hwChan HardwareChannels, orderChan OrderChannels, netChan NetworkChannels) {
	elev := Elevator{
		Id:         GetElevIP(),
		Floor:      hw.GetFloor(),
		Dir:        hw.MD_Stop,
		State:      IDLE,
		Online:     false,
		OrderQueue: [NUMFLOORS][NUMBUTTONS]bool{},
		Mobile:     true,
	}

	doorTimeout := time.NewTimer(DOOROPENTIME * time.Millisecond)
	doorTimeout.Stop()

	engineFailure := time.NewTimer(3 * time.Second)
	engineFailure.Stop()

	var rememberDir hw.MotorDirection
	var obstructionCounter = 0

	for {
		switch elev.State {
		case IDLE:
			rememberDir = elev.Dir
			select {
			case isOnline := <-netChan.IsOnline:
				elev.Online = isOnline
			case newOrder := <-orderChan.LocalOrder:
				elev.Id = newOrder.Id
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(DOOROPENTIME * time.Millisecond)
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
			case isOnline := <-netChan.IsOnline:
				elev.Online = isOnline
			case newOrder := <-orderChan.LocalOrder:
				elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				break
			case newFloor := <-hwChan.HwFloor:
				elev.Floor = newFloor
				elev.Mobile = true
				if ShouldStop(elev) {
					parameters := ClearOrdersParams{Elev: elev}
					elev = ClearOrdersAtCurrentFloor(parameters)
					rememberDir = elev.Dir
					elev.Dir = hw.MD_Stop
					elev.State = DOOROPEN
					doorTimeout.Reset(DOOROPENTIME * time.Millisecond)
					engineFailure.Stop()
				} else {
					engineFailure.Reset((3 * time.Second))
				}
				break
			case <-engineFailure.C:
				fmt.Println("ENGINE FAILURE")
				if elev.Mobile {
					elev.Mobile = false
					netChan.InmobileElev <- elev
				}
				engineFailure.Reset((1 * time.Second))
			}
		case DOOROPEN:
			select {
			case isOnline := <-netChan.IsOnline:
				elev.Online = isOnline
			case newOrder := <-orderChan.LocalOrder:
				if elev.Floor == newOrder.Floor {
					elev.State = DOOROPEN
					doorTimeout.Reset(DOOROPENTIME * time.Millisecond)
				} else {
					elev.OrderQueue[newOrder.Floor][newOrder.Button] = true
				}
				break
			case <-doorTimeout.C:
				obstructed := hw.GetObstruction()
				elev.Dir = ChooseDirection(elev, rememberDir)
				if obstructed {
					fmt.Println("OBSTRUCTED")
					doorTimeout.Reset(DOOROPENTIME * time.Millisecond)
					elev.State = DOOROPEN
					elev.Dir = hw.MD_Stop
					obstructionCounter++
					if obstructionCounter == 3 { // Reassign order
						obstructionCounter = 0
						if elev.Mobile {
							elev.Mobile = false
							netChan.InmobileElev <- elev
						}
					}
				} else if elev.Dir == hw.MD_Stop {
					elev.State = IDLE
					elev.Mobile = true
					engineFailure.Stop()
					obstructionCounter = 0
				} else {
					elev.State = MOVING
					elev.Mobile = true
					engineFailure.Reset((3 * time.Second))
					obstructionCounter = 0
				}
				break
			}
		}
		enrollHardware(elev)
		writeToBackup(elev)

		orderChan.LocalElevUpdate <- elev
		orderChan.RecieveElevUpdate <- elev
	}
}

func clearAllLights() {
	hw.SetDoorOpenLamp(false)
	for floor := 0; floor < NUMFLOORS; floor++ {
		for btn := 0; btn < NUMBUTTONS; btn++ {
			hw.SetButtonLamp(hw.ButtonType(btn), floor, false)
		}
	}
}

func enrollHardware(elev Elevator) {
	hw.SetFloorIndicator(elev.Floor)
	hw.SetMotorDirection(elev.Dir)
	hw.SetDoorOpenLamp(DOOROPEN == elev.State)
}
