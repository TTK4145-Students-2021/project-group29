package Executer

import "../Driver/elevio"
import "fmt"
import "../Assigner"

func orderAbove(Elev e){
	for fl:=e.Floor+1; fl<NumFloors;fl++{
		for btn:=0; btn<NumButtons;btn++{
			if e.OrderQueue[fl][btn]{
				return 1;
			}
		}
	}
	return 0;
}


func orderBelow(Elev e) {
	for f:0; f<e.Floor; f++{
		for btn:=0; btn<NumButtons;btn++{
			if e.OrderQueue[fl][btn]{
				return 1;
			}
		}
	}
}

func orderChooseDir(Elev e){
	dir:= e.Dir;
	switch dir {
		case MD_Up:
			if orderAbove(e){
				return DM_UP;
			}
			else if orderBelow(e){
				return DM_DOWN;
			}
			else{
				return DM_STOP;
			}
		case DM_DOWN:
		case DM_STOP:
			if orderBelow(e){
				return DM_DOWN;
			}
			else if orderAbove(e){
				return DM_UP;
			}
			else{
				return DM_STOP;
			}
	}
}


func orderShouldStop(Elev e) {
	switch dir := e.Dir; dir {
	case MD_Down:
		return 
			e.OrderQueue[e.Floor][BT_HallDown] 	||
			e.OrderQueue[e.Floor][BT_Cab] 		||
			!orderBelow(e);
	case MD_Up:
		return
			e.orderAbove[e.Floor][BT_HallUp]   ||
            e.requests[e.Floor][BT_Cab]  	   ||
            !orderAbove(e);
	case MD_Stop
	default:
		return 1;
	}
}

func orderClearAtCurrentFloor() {
	
}