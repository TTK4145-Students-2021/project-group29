package Executer

import co "../Common"

func ordersAbove(elev co.Elevator) {
	for floor := elev.Floor + 1; floor < co.NumFloors; floor++ {
		for btn := 0; btn < co.NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elev co.Elevator) {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < co.NumButtons; btn++ {
			if elev.OrderQueue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func chooseDirection(elev co.Elevator) {
	switch elev.Dir {
	case MD_Up:
		if ordersAbove(elev) {
			return MD_UP
		} else if ordersBelow(elev) {
			return MD_DOWN
		} else {
			return MD_STOP
		}
	case MD_DOWN:
		if ordersBelow(elev) {
			return MD_DOWN
		} else if ordersAbove(elev) {
			return MD_UP
		} else {
			return MD_STOP
		}
	case MD_STOP:
		if ordersAbove(elev) {
			return MD_UP
		} else if ordersAbove(elev) {
			return MD_DOWN
		} else {
			return MD_STOP
		}
	}
}

func shouldStop(elev co.Elevator) {
	switch elev.Dir {
	case MD_Down:
		return
		elev.OrderQueue[elev.Floor][BT_HallDown] ||
			elev.OrderQueue[elev.Floor][BT_Cab] ||
			!ordersBelow(elev)
	case MD_Up:
		return
		elev.orderAbove[elev.Floor][BT_HallUp] ||
			elev.requests[elev.Floor][BT_Cab] ||
			!ordersAbove(elev)
	case MD_Stop:
	default:
		return 1
	}
}

func clearOrdersAtCurrentFloor(elev co.Elevator) {
	//Assuming that everyone enters at the current floor
	/* for(Button btn = 0; btn < N_BUTTONS; btn++){
	       e.requests[e.floor][btn] = 0;
	   }
	   break;
	*/
	// Assuming that only those that want to travel in the current direction
	// enter the elevator, and keep waiting outside otherwise
	elev.OrderQueue[elev.Floor][BT_Cab] = 0
	switch dir := elev.Dir; dir {
	case MD_UP:
		elev.OrderQueue[elev.Floor][BT_HallUp] = 0
		if !orderAbove(elev) {
			elev.OrderQueue[elev.Floor][BT_HallDown] = 0
		}
		break

	case DM_DOWN:
		elev.OrderQueue[elev.Floor][BT_HallDown] = 0
		if !orderBelow(elev) {
			elev.OrderQueue[elev.Floor][BT_HallUp] = 0
		}
		break

	case DM_STOP:
	default:
		elev.OrderQueue[elev.Floor][BT_HallUp] = 0
		elev.OrderQueue[elev.Floor][BT_HallDown] = 0
		break
	}

}
