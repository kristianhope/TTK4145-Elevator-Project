package requests

import (
	"Project/config"
	"Project/localElevator/elevator"
	"Project/localElevator/elevio"
)

type Action struct {
	dirn      elevio.MotorDirection
	behaviour elevator.ElevatorBehaviour
}

func requests_above(e elevator.Elevator) int {
	for f := e.Floor; f < config.NumFloors; f++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if e.Requests[f][btn] != 0 {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e elevator.Elevator) int {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if e.Requests[f][btn] != 0 {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e elevator.Elevator) int {
	for btn := 0; btn < config.NumFloors; btn++ {
		if e.Requests[e.Floor][btn] != 0 {
			return 1
		}
	}
	return 0
}

func requests_nextAction(e elevator.Elevator) Action {
	switch e.Dirn {
	case elevio.MD_Up:
		if requests_above(e) != 0 {
			return Action{elevio.MD_Up, elevator.EB_Moving}
		} else if requests_here(e) != 0 {
			return Action{elevio.MD_Down, elevator.EB_DoorOpen}
		} else if requests_below(e) != 0 {
			return Action{elevio.MD_Down, elevator.EB_Moving}
		} else {
			return Action{elevio.MD_Stop, elevator.EB_Idle}
		}
	case elevio.MD_Down:
		if requests_below(e) != 0 {
			return Action{elevio.MD_Down, elevator.EB_Moving}
		} else if requests_here(e) != 0 {
			return Action{elevio.MD_Up, elevator.EB_DoorOpen}
		} else if requests_above(e) != 0 {
			return Action{elevio.MD_Up, elevator.EB_Moving}
		} else {
			return Action{elevio.MD_Stop, elevator.EB_Idle}
		}
	case elevio.MD_Stop:
		if requests_here(e) != 0 {
			return Action{elevio.MD_Stop, elevator.EB_DoorOpen}
		} else if requests_above(e) != 0 {
			return Action{elevio.MD_Up, elevator.EB_Moving}
		} else if requests_below(e) != 0 {
			return Action{elevio.MD_Down, elevator.EB_Moving}
		} else {
			return Action{elevio.MD_Stop, elevator.EB_Idle}
		}
	default:
		return Action{elevio.MD_Stop, elevator.EB_Idle}
	}
}

func Requests_shouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return (e.Requests[e.Floor][elevio.BT_HallDown] != 0) || (e.Requests[e.Floor][elevio.BT_Cab] != 0) || (requests_below(e) == 0)
	case elevio.MD_Up:
		return (e.Requests[e.Floor][elevio.BT_HallUp] != 0) || (e.Requests[e.Floor][elevio.BT_Cab] != 0) || (requests_above(e) == 0)
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	switch e.config.clearRequestVariant {
	case elevator.CV_All:
		return e.Floor == btn_floor
	case elevator.CV_InDirn:
		return e.Floor == btn_floor &&
			((e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
				(e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
				e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab)

	default:
		return false
	}
}

func requests_clearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.config.clearRequestVariant {
	case elevator.CV_All:
		for btn := 0; btn < config.NumButtons; btn++ {
			e.Requests[e.Floor][btn] = 0
		}
	case elevator.CV_InDirn:
		e.Requests[e.Floor][elevio.BT_Cab] = 0
		switch e.Dirn {
		case elevio.MD_Up:
			if (requests_above(e) == 0) && (e.Requests[e.Floor][elevio.BT_HallUp] == 0) {
				e.Requests[e.Floor][elevio.BT_HallDown] = 0
			}
			e.Requests[e.Floor][elevio.BT_HallUp] = 0
			break
		case elevio.MD_Down:
			if (requests_below(e) == 0) && (e.Requests[e.Floor][elevio.BT_HallDown] == 0) {
				e.Requests[e.Floor][elevio.BT_HallUp] = 0
			}
			e.Requests[e.Floor][elevio.BT_HallDown] = 0
			break
		case elevio.MD_Stop:
		default:
			e.Requests[e.Floor][elevio.BT_HallUp] = 0
			e.Requests[e.Floor][elevio.BT_HallDown] = 0
			break
		}
		break
	default:
		break
	}
	return e
}