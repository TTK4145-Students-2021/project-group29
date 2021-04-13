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

func ChooseDirection(elev Elevator, rememberDir hw.MotorDirection) hw.MotorDirection {
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

func ShouldStop(elev Elevator) bool {
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

func ClearOrdersAtCurrentFloor(params ClearOrdersParams) Elevator {
	params.Elev.OrderQueue[params.Elev.Floor][hw.BT_Cab] = false
	haveFunction := !(params.IfEqual == nil)
	switch params.Elev.Dir {
	case hw.MD_Up:
		if haveFunction { // check ifRequest
			params.IfEqual(hw.BT_HallUp, params.Elev.Floor)
		}
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallUp] = false
		if !ordersAbove(params.Elev) {
			if haveFunction { // check ifRequest
				params.IfEqual(hw.BT_HallDown, params.Elev.Floor)
			}
			params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallDown] = false
		}
		break

	case hw.MD_Down:
		if haveFunction { // check ifRequest
			params.IfEqual(hw.BT_HallDown, params.Elev.Floor)
		}
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallDown] = false
		if !ordersBelow(params.Elev) {
			if haveFunction { // check ifRequest
				params.IfEqual(hw.BT_HallUp, params.Elev.Floor)
			}
			params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallUp] = false
		}
		break

	case hw.MD_Stop:
		if haveFunction { // check ifRequest
			params.IfEqual(hw.BT_HallUp, params.Elev.Floor)
			params.IfEqual(hw.BT_HallDown, params.Elev.Floor)
		}
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallUp] = false
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallDown] = false
		break

	default:
		if haveFunction { // check ifRequest
			params.IfEqual(hw.BT_HallUp, params.Elev.Floor)
			params.IfEqual(hw.BT_HallDown, params.Elev.Floor)
		}
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallUp] = false
		params.Elev.OrderQueue[params.Elev.Floor][hw.BT_HallDown] = false
		break
	}

	return params.Elev
}
