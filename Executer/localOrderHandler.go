package Executer

import (
	. "../Common"
	hw "../Driver/elevio"
)

func ordersAbove(elev Elevator) bool {
	for floor := elev.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elev Elevator) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func chooseDirection(elev Elevator, rememberDir hw.MotorDirection) hw.MotorDirection {
	switch rememberDir {
	case hw.MD_Up:
		if ordersAbove(elev) {
			return hw.MD_Up
		} else if ordersBelow(elev) {
			return hw.MD_Down
		} else {
			return hw.MD_Stop
		}
	case hw.MD_Down:
		if ordersBelow(elev) {
			return hw.MD_Down
		} else if ordersAbove(elev) {
			return hw.MD_Up
		} else {
			return hw.MD_Stop
		}
	case hw.MD_Stop:
		if ordersAbove(elev) {
			return hw.MD_Up
		} else if ordersBelow(elev) {
			return hw.MD_Down
		} else {
			return hw.MD_Stop
		}
	}
	return hw.MD_Stop
}

func shouldStop(elev Elevator) bool {
	switch elev.Dir {
	case hw.MD_Down:
		return elev.OrderQueue[elev.Floor][hw.BT_HallDown] ||
			elev.OrderQueue[elev.Floor][hw.BT_Cab] ||
			!ordersBelow(elev)
	case hw.MD_Up:
		return elev.OrderQueue[elev.Floor][hw.BT_HallUp] ||
			elev.OrderQueue[elev.Floor][hw.BT_Cab] ||
			!ordersAbove(elev)
	case hw.MD_Stop:
	}
	return true
}

func clearOrdersAtCurrentFloor(elev Elevator) Elevator {

	elev.OrderQueue[elev.Floor][hw.BT_Cab] = false
	switch elev.Dir {
	case hw.MD_Up:
		elev.OrderQueue[elev.Floor][hw.BT_HallUp] = false
		if !ordersAbove(elev) {
			elev.OrderQueue[elev.Floor][hw.BT_HallDown] = false
		}
		break

	case hw.MD_Down:
		elev.OrderQueue[elev.Floor][hw.BT_HallDown] = false
		if !ordersBelow(elev) {
			elev.OrderQueue[elev.Floor][hw.BT_HallUp] = false
		}
		break

	case hw.MD_Stop:
	
	default:
		elev.OrderQueue[elev.Floor][hw.BT_HallUp] = false
		elev.OrderQueue[elev.Floor][hw.BT_HallDown] = false
		break
	}
	return elev
}
