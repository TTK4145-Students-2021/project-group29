package Executer

import (
	"reflect"

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

func ClearOrdersAtCurrentFloor(p Params) Elevator { //, onClearedRequest func(hw.ButtonType,int)
	p.Elev.OrderQueue[p.Elev.Floor][hw.BT_Cab] = false
	switch p.Elev.Dir {
	case hw.MD_Up:
		// check ifRequest
		if !reflect.ValueOf(p.Func).IsZero() {
			p.Func(hw.BT_HallUp, p.Elev.Floor)
		}
		p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallUp] = false
		if !ordersAbove(p.Elev) {
			// check ifRequest
			if !reflect.ValueOf(p.Func).IsZero() {
				p.Func(hw.BT_HallDown, p.Elev.Floor)
			}
			p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallDown] = false
		}
		break

	case hw.MD_Down:
		// check if request
		if !reflect.ValueOf(p.Func).IsZero() {
			p.Func(hw.BT_HallDown, p.Elev.Floor)
		}
		p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallDown] = false
		if !ordersBelow(p.Elev) {
			// check ifRequest
			if !reflect.ValueOf(p.Func).IsZero() {
				p.Func(hw.BT_HallUp, p.Elev.Floor)
			}
			p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallUp] = false
		}
		break

	case hw.MD_Stop:

	default:
		p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallUp] = false
		p.Elev.OrderQueue[p.Elev.Floor][hw.BT_HallDown] = false
		break
	}

	return p.Elev
}
